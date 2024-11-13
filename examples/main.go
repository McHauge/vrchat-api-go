package main

import (
	"github.com/lyzcoote/vrchat-api-go"
)

func main() {
	client := vrchat.NewClient("https://api.vrchat.cloud/api/1", "MyApp-Test/1.0")

	err := client.AuthenticateTOTP("username", "password", "totp")
	if err != nil {
		panic(err)
	}

	user, err := client.GetCurrentUser()
	if err != nil {
		panic(err)
	}

	println("logged in as ", user.DisplayName)
}
