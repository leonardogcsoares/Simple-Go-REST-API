// Package models provides a structs that will be Marshaled and Unmarshaled between
package models

type (
	// User is an abstraction for IMEI number to be received from request
	User struct {
		Imei string `json:"imei"`
	}

	// Token is an abstraction for the token to be received from request
	Token struct {
		TokenJwt string `json:"token"`
	}
)
