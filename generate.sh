#!/bin/bash

# If you have Wget
# wget https://vrchat.community/openapi.yaml -O openapi.yaml

# If you have Curl
curl.exe --output openapi.yaml --url https://vrchat.community/openapi.yaml

# Fix openapi.yaml schema issues
cd utils && go run fix_schema.go && cd ..

# Old Style generator
go install github.com/mchauge/openapi-codegen@latest
openapi-codegen -i ./openapi.yaml -o . -p vrchat -f=false

# fix Schema fail with Group_All
sed -i -E 's/GroupPermissions GroupPermissions/GroupPermissionsAll GroupPermissions/' ./schema.gen.go

# Format code and schema
go fmt client.gen.go
go fmt schema.gen.go

# Added UserAgent to NewClient function, otherwise it will return 403
sed -i 's/func NewClient(baseURL string)/func NewClient(baseURL string, UserAgent string)/' ./client.gen.go
sed -i 's/client: resty.New().SetBaseURL(baseURL),/client: resty.New().SetBaseURL(baseURL).SetHeader("User-Agent", UserAgent),/' ./client.gen.go

# Added SetClient and GetClient function to easily get the underlying http.client, for better cokkies handeling particulaly
sed -i '$a \
\
func (c *Client) SetClient(client *resty.Client) {\
    c.client = client\
}\n\
func (c *Client) GetClient() *resty.Client {\
    return c.client\
}' ./client.gen.go

# Import net/http for the cookie related functions below
# sed -i '/^import (/a \    "net/http"' ./client.gen.go

# Added SetCookie and GetCookies function to easily manage cookies (does not seem to work correctly)
# sed -i '$a \
# \
# func (c *Client) SetCookie(cookie *http.Cookie) {\
#     c.client.SetCookie(cookie)\
# }\n\
# func (c *Client) GetCookies() []*http.Cookie {\
#     return c.client.Cookies\
# }' ./client.gen.go

# (Not needed anymore so commented out)
# sed -i '/^import (/a \    "reflect"' ./client.gen.go
# sed -i '/^import (/a \    "strconv"' ./client.gen.go

# sed -i '/func (c \*Client) UpdateUser(params UpdateUserParams)/,/^}/d' ./client.gen.go

# sed -i '$a \
# \
# func (c *Client) UpdateUser(userId UpdateUserParams, params UpdateUserRequest) (*CurrentUserResponse, error) {\
#     path := "/users/{userId}"\
#     // Replace path parameters and prepare query parameters\
#     queryParams := make(map[string]string)\
# \
#     // Iterate over the fields of the struct and add them to the queryParams\
#     val := reflect.ValueOf(params)\
#     typ := reflect.TypeOf(params)\
# \
#     for i := 0; i < val.NumField(); i++ {\
#         field := val.Field(i)\
#         jsonTag := typ.Field(i).Tag.Get("json")\
#         // Split json tag by comma and get the first part\
#         jsonTag = strings.Split(jsonTag, ",")[0]\
# \
#         // Skip empty fields\
#         if !field.IsValid() || (field.Kind() == reflect.String && field.String() == "") ||\
#             (field.Kind() == reflect.Slice && field.Len() == 0) ||\
#             (field.Kind() == reflect.Bool && !field.Bool()) {\
#             continue\
#         }\
# \
#         // Adding field to queryParams\
#         switch field.Kind() {\
#         case reflect.String:\
#             queryParams[jsonTag] = field.String()\
#         case reflect.Slice:\
#             if field.Len() > 0 {\
#                 // Handle string arrays like bioLinks or tags\
#                 var strSlice []string\
#                 for j := 0; j < field.Len(); j++ {\
#                     strSlice = append(strSlice, field.Index(j).String())\
#                 }\
#                 queryParams[jsonTag] = strings.Join(strSlice, ",")\
#             }\
#         case reflect.Bool:\
#             queryParams[jsonTag] = strconv.FormatBool(field.Bool())\
#         case reflect.Int64, reflect.Int:\
#             queryParams[jsonTag] = fmt.Sprintf("%d", field.Int())\
#         case reflect.Uint:\
#             queryParams[jsonTag] = fmt.Sprintf("%d", field.Uint())\
#         default:\
#             panic("unhandled default case")\
#         }\
#     }\
# \
#     path = strings.ReplaceAll(path, "{userId}", fmt.Sprintf("%v", userId.UserId))\
# \
#     // Create request\
#     req := c.client.R()\
#     // Set body\
#     req.SetBody(queryParams)\
#     // Set response object\
#     var result CurrentUserResponse\
#     req.SetResult(&result)\
# \
#     // Send request\
#     resp, err := req.Put(path)\
#     if err != nil {\
#         return nil, fmt.Errorf("error sending request: %w", err)\
#     }\
# \
#     // Check for successful status code\
#     if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {\
#         return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode(), resp.String())\
#     }\
#     return &result, nil\
# }' ./client.gen.go
