package main

import (
	"github.com/lyzcoote/vrchat-api-go"
	"log"
)

func main() {
	client := vrchat.NewClient("https://api.vrchat.cloud/api/1", "MyApp-Test/1.0")

	resp, err := client.Authenticate("username", "password")
	if err != nil {
		panic(err)
	}

	if resp == "{\"requiresTwoFactorAuth\":[\"totp\",\"otp\"]}" {
		resp, err = client.VerifyTOTP("username", "password", "123456")
		if err != nil {
			panic(err)
		}
		log.Println(resp)
	}

	user, err := client.GetCurrentUser()
	if err != nil {
		panic(err)
	}

	println("logged in as ", user.DisplayName)
}
