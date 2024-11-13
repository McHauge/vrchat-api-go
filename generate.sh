#!/bin/bash

wget https://vrchatapi.github.io/specification/openapi.yaml -O openapi.yaml

go install github.com/mayocream/openapi-codegen@latest

openapi-codegen -i ./openapi.yaml -o . -p vrchat

# Added UserAgent to NewClient function, otherwise it will return 403
sed -i 's/func NewClient(baseURL string)/func NewClient(baseURL string, UserAgent string)/' ./client.gen.go
sed -i 's/client: resty.New().SetBaseURL(baseURL),/client: resty.New().SetBaseURL(baseURL).SetHeader("User-Agent", UserAgent),/' ./client.gen.go

# Import net/http for the cookie related functions below
sed -i '/^import (/a \    "net/http"' ./client.gen.go

# Added SetCookie and GetCookies function to easily manage cookies
sed -i '$a \
\
func (c *Client) SetCookie(cookie *http.Cookie) {\
    c.client.SetCookie(cookie)\
}\n\
func (c *Client) GetCookies() []*http.Cookie {\
    return c.client.Cookies\
}' ./client.gen.go