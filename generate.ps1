Invoke-WebRequest -Uri "https://vrchatapi.github.io/specification/openapi.yaml" -OutFile "openapi.yaml"

go install github.com/mayocream/openapi-codegen@latest

openapi-codegen -i ./openapi.yaml -o . -p vrchat

# Added UserAgent to NewClient function, otherwise it will return 403
(Get-Content ./client.gen.go) -replace 'func NewClient\(baseURL string\)', 'func NewClient(baseURL string, UserAgent string)' | Set-Content ./client.gen.go
(Get-Content ./client.gen.go) -replace 'client: resty.New\(\).SetBaseURL\(baseURL\),', 'client: resty.New().SetBaseURL(baseURL).SetHeader("User-Agent", UserAgent),' | Set-Content ./client.gen.go

# Import net/http for the cookie related functions below
(Get-Content ./client.gen.go) -replace '^import \(', 'import (    
    "net/http"' | Set-Content ./client.gen.go

# Import reflect for the cookie related functions below
(Get-Content ./client.gen.go) -replace '^import \(', 'import (
    "reflect"' | Set-Content ./client.gen.go

# Import strconv for the cookie related functions below
(Get-Content ./client.gen.go) -replace '^import \(', 'import (
    "strconv"' | Set-Content ./client.gen.go

# Added SetCookie and GetCookies function to easily manage cookies
Add-Content -Path ./client.gen.go -Value @"

func (c *Client) SetCookie(cookie *http.Cookie) {
    c.client.SetCookie(cookie)
}

func (c *Client) GetCookies() []*http.Cookie {
    return c.client.Cookies
}
"@


(Get-Content ./client.gen.go) |
        ForEach-Object {
            if ($_ -match "func \(c \*Client\) UpdateUser\(params UpdateUserParams\)") {
                $skip = $true
            }
            if ($skip) {
                if ($_ -match "type GetUserGroupsParams") {
                    $skip = $false
                }
            }
            if (-not $skip) {
                $_
            }
        } | Set-Content ./client.gen.go




Add-Content -Path ./client.gen.go -Value @"

func (c *Client) UpdateUser(userId UpdateUserParams, params UpdateUserRequest) (*CurrentUserResponse, error) {
	path := "/users/{userId}"
	// Replace path parameters and prepare query parameters
	queryParams := make(map[string]string)

	// Iterate over the fields of the struct and add them to the queryParams
	val := reflect.ValueOf(params)
	typ := reflect.TypeOf(params)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		jsonTag := typ.Field(i).Tag.Get("json")
		// Split json tag by comma and get the first part
		jsonTag = strings.Split(jsonTag, ",")[0]

		// Skip empty fields
		if !field.IsValid() || (field.Kind() == reflect.String && field.String() == "") ||
			(field.Kind() == reflect.Slice && field.Len() == 0) ||
			(field.Kind() == reflect.Bool && !field.Bool()) {
			continue
		}

		// Adding field to queryParams
		switch field.Kind() {
		case reflect.String:
			queryParams[jsonTag] = field.String()
		case reflect.Slice:
			if field.Len() > 0 {
				// Handle string arrays like bioLinks or tags
				var strSlice []string
				for j := 0; j < field.Len(); j++ {
					strSlice = append(strSlice, field.Index(j).String())
				}
				queryParams[jsonTag] = strings.Join(strSlice, ",")
			}
		case reflect.Bool:
			queryParams[jsonTag] = strconv.FormatBool(field.Bool())
		case reflect.Int64, reflect.Int:
			queryParams[jsonTag] = fmt.Sprintf("%d", field.Int())
		case reflect.Uint:
			queryParams[jsonTag] = fmt.Sprintf("%d", field.Uint())
		default:
			panic("unhandled default case")
		}
	}

	path = strings.ReplaceAll(path, "{userId}", fmt.Sprintf("%v", userId.UserId))

	// Create request
	req := c.client.R()
	// Set body
	req.SetBody(queryParams)
	// Set response object
	var result CurrentUserResponse
	req.SetResult(&result)

	// Send request
	resp, err := req.Put(path)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	// Check for successful status code
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode(), resp.String())
	}
	return &result, nil
}

"@