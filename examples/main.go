package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/lyzcoote/vrchat-api-go"
)

func main() {
	client := vrchat.NewClient("https://api.vrchat.cloud/api/1", "My-App-Name/1.0")

	resp, err := client.Authenticate("Username", "password")
	if err != nil {
		panic(err)
	}

	var authCode string
	if strings.Contains(resp, "requiresTwoFactorAuth") {
		log.Println("2FA required, please enter code:")
		if _, err := fmt.Scan(&authCode); err != nil {
			panic(fmt.Sprintf("Unable to read 2FA code: %v", err))
		}
	}

	if resp == "{\"requiresTwoFactorAuth\":[\"totp\",\"otp\"]}" || resp == "{\"requiresTwoFactorAuth\":[\"emailOtp\"]}" {
		resp, err = client.VerifyTOTP("username", "password", authCode)
		if err != nil {
			panic(err)
		}
		log.Println(resp)
	}

	user, err := client.GetCurrentUser()
	if err != nil {
		panic(err)
	}

	println("logged in as ", user.DisplayName, user.Id)

	update := vrchat.UpdateUserParams{
		UserId:            user.Id,
		Status:            vrchat.UserStatusActive,
		StatusDescription: "I'm a bot",
	}

	resp1, err := client.UpdateUser(update)
	if err != nil {
		panic(err)
	}
	println(resp1.Status, resp1.StatusDescription)

}
