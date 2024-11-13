package vrchat

import (
	"fmt"
)

// AuthenticateRecoveryOTP authenticates the client with the VRChat API using the email recovery OTP code.
func (c *Client) AuthenticateRecoveryOTP(username, password, totp string) error {
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
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to authenticate: %s", resp.String())
	}

	cookies := resp.Cookies()
	c.client.SetCookies(cookies)

	return nil
}

// AuthenticateEmailOTP authenticates the client with the VRChat API using the email OTP code.
func (c *Client) AuthenticateEmailOTP(username, password, totp string) error {
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
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to authenticate: %s", resp.String())
	}

	cookies := resp.Cookies()
	c.client.SetCookies(cookies)

	return nil
}

// AuthenticateTOTP authenticates the client with the VRChat API using the TOTP code.
func (c *Client) AuthenticateTOTP(username, password, totp string) error {
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
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to authenticate: %s", resp.String())
	}

	cookies := resp.Cookies()
	c.client.SetCookies(cookies)

	return nil
}
