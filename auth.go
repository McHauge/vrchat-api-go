package vrchat

import (
	"fmt"
)

// Authenticate authenticates the client with the VRChat API using the username and password.
func (c *Client) Authenticate(username, password string) (string, error) {
	c.client.SetHeaders(map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	})
	resp, err := c.client.R().
		SetBasicAuth(username, password).
		Get("/auth/user")
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("failed to authenticate: %s", resp.String())
	}

	cookies := resp.Cookies()
	c.client.SetCookies(cookies)

	return resp.String(), nil
}

// VerifyRecoveryOTP authenticates the client with the VRChat API using the email recovery OTP code.
func (c *Client) VerifyRecoveryOTP(username, password, totp string) (string, error) {
	c.client.SetHeaders(map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	})
	resp, err := c.client.R().
		SetBasicAuth(username, password).
		SetBody(map[string]string{
			"code": totp,
		}).
		Post("/auth/twofactorauth/otp/verify")
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("failed to authenticate: %s", resp.String())
	}

	cookies := resp.Cookies()
	c.client.SetCookies(cookies)

	return resp.String(), nil
}

// VerifyEmailOTP authenticates the client with the VRChat API using the email OTP code.
func (c *Client) VerifyEmailOTP(username, password, totp string) (string, error) {
	c.client.SetHeaders(map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	})
	resp, err := c.client.R().
		SetBasicAuth(username, password).
		SetBody(map[string]string{
			"code": totp,
		}).
		Post("/auth/twofactorauth/emailotp/verify")
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("failed to authenticate: %s", resp.String())
	}

	cookies := resp.Cookies()
	c.client.SetCookies(cookies)

	return resp.String(), nil
}

// VerifyTOTP authenticates the client with the VRChat API using the TOTP code.
func (c *Client) VerifyTOTP(username, password, totp string) (string, error) {
	c.client.SetHeaders(map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	})
	resp, err := c.client.R().
		SetBasicAuth(username, password).
		SetBody(map[string]string{
			"code": totp,
		}).
		Post("/auth/twofactorauth/totp/verify")
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("failed to authenticate: %s", resp.String())
	}

	cookies := resp.Cookies()
	c.client.SetCookies(cookies)

	return resp.String(), nil
}
