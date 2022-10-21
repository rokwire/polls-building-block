package model

// GroupNotification wrapper for sending a notifications to members of a group
type GroupNotification struct {
	MemberStatuses []string          `json:"member_statuses"` // default: ["admin", "member"]
	Members        []UserRef         `json:"members"`
	Sender         *Sender           `json:"sender"`
	Subject        string            `json:"subject"`
	Topic          *string           `json:"topic"`
	Body           string            `json:"body"`
	Data           map[string]string `json:"data"`
} // @name GroupNotification
