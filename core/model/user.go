package model

import "github.com/rokwire/core-auth-library-go/v2/tokenauth"

// User auth wrapper
type User struct {
	Token  string
	Claims tokenauth.Claims
}

// UserRef reference for a concrete user which is member of a group
type UserRef struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
} // @name MemberRecipient

// Sender Wraps sender type and user ref
type Sender struct {
	Type string   `json:"type" bson:"type"` // user or system
	User *UserRef `json:"user,omitempty" bson:"user,omitempty"`
} // @name Sender
