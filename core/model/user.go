package model

import "github.com/rokwire/core-auth-library-go/tokenauth"

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

// DeletedUserData represents a user-deleted
type DeletedUserData struct {
	AppID       string              `json:"app_id"`
	Memberships []DeletedMembership `json:"memberships"`
	OrgID       string              `json:"org_id"`
}

// DeletedMembership defines model for DeletedMembership.
type DeletedMembership struct {
	AccountID string                  `json:"account_id"`
	Context   *map[string]interface{} `json:"context,omitempty"`
}

// CoreAccount represents an account in the Core BB
type CoreAccount struct {
	ID      string      `json:"id" bson:"id"`
	Profile CoreProfile `json:"profile" bson:"profile"`
} //@name CoreAccount

// CoreProfile represents a profile in the Core BB
type CoreProfile struct {
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
} //@name CoreProfile
