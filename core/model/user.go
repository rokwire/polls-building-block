package model

import (
	"github.com/rokwire/core-auth-library-go/v2/tokenauth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

// UserDataResponse wraps polls user data
type UserDataResponse struct {
	PollsUserData          []PollsUserData          `json:"polls"`
	PollsResponseUserData  []PollsResponseUserData  `json:"polls_responses"`
	SurveysUserData        []SurveysUserData        `json:"surveys"`
	SurveyResponseUserData []SurveyResponseUserData `json:"surveys_responses"`
} //@name UserDataResponse

// PollsUserData  wraps polls user data
type PollsUserData struct {
	ID     primitive.ObjectID `json:"id"`
	UserID string             `json:"user_id"`
} // @name PollsUserData

// PollsResponseUserData  wraps polls user data
type PollsResponseUserData struct {
	ID     primitive.ObjectID `json:"id"`
	UserID string             `json:"user_id"`
} // @name PollsResponseUserData

// SurveyResponseUserData wraps the user data survey response
type SurveyResponseUserData struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

// SurveysUserData wraps the user data record
type SurveysUserData struct {
	ID        string `json:"id"`
	CreatorID string `json:"creator_id"`
}
