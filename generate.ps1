Invoke-WebRequest -Uri "https://vrchatapi.github.io/specification/openapi.yaml" -OutFile "openapi.yaml"

go install github.com/mayocream/openapi-codegen@latest

openapi-codegen -i ./openapi.yaml -o . -p vrchat

# Added UserAgent to NewClient function, otherwise it will return 403
(Get-Content ./client.gen.go) -replace 'func NewClient\(baseURL string\)', 'func NewClient(baseURL string, UserAgent string)' | Set-Content ./client.gen.go
(Get-Content ./client.gen.go) -replace 'client: resty.New\(\).SetBaseURL\(baseURL\),', 'client: resty.New().SetBaseURL(baseURL).SetHeader("User-Agent", UserAgent),' | Set-Content ./client.gen.go

# Import net/http for the cookie related functions below
(Get-Content ./client.gen.go) -replace '^import \(', 'import (    
    "net/http"' | Set-Content ./client.gen.go

# Added SetCookie and GetCookies function to easily manage cookies
Add-Content -Path ./client.gen.go -Value @"

func (c *Client) SetCookie(cookie *http.Cookie) {
    c.client.SetCookie(cookie)
}

func (c *Client) GetCookies() []*http.Cookie {
    return c.client.Cookies
}
"@