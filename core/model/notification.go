package model

// NotificationMessage wrapper for internal message
type NotificationMessage struct {
	Priority   int                     `json:"priority" bson:"priority"`
	Recipients []NotificationRecipient `json:"recipients" bson:"recipients"`
	Topic      *string                 `json:"topic" bson:"topic"`
	Subject    string                  `json:"subject" bson:"subject"`
	Sender     *NotificationSender     `json:"sender,omitempty" bson:"sender,omitempty"`
	Body       string                  `json:"body" bson:"body"`
	Data       map[string]string       `json:"data" bson:"data"`
}

// NotificationRecipient recipients wrapper struct
type NotificationRecipient struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

// NotificationSender notification sender
type NotificationSender struct {
	Type string       `json:"type" bson:"type"` // user or system
	User *CoreUserRef `json:"user,omitempty" bson:"user,omitempty"`
}

// CoreUserRef user reference that contains ExternalID & Name
type CoreUserRef struct {
	UserID *string `json:"user_id" bson:"user_id"`
	Name   *string `json:"name" bson:"name"`
}
