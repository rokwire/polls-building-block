package model

// NotificationMessage wrapper for internal message
type NotificationMessage struct {
	Recipients []UserRef         `json:"recipients" bson:"recipients"`
	Topic      *string           `json:"topic" bson:"topic"`
	Subject    string            `json:"subject" bson:"subject"`
	Sender     *Sender           `json:"sender,omitempty" bson:"sender,omitempty"`
	Body       string            `json:"body" bson:"body"`
	Data       map[string]string `json:"data" bson:"data"`
	AppID      string            `json:"app_id" bson:"app_id"`
	OrgID      string            `json:"org_id" bson:"org_id"`
}
