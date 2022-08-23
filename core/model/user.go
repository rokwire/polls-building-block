package model

import "github.com/rokwire/core-auth-library-go/tokenauth"

// User auth wrapper
type User struct {
	Token  string
	Claims tokenauth.Claims
}
