package model

import (
	"time"
)

// Group struct wrapper
type Group struct {
	ID              string `json:"id"`
	ClientID        string `json:"client_id"`
	Category        string `json:"category"`
	Title           string `json:"title"`
	Privacy         string `json:"privacy"`
	HiddenForSearch bool   `json:"hidden_for_search"`
	Members         []struct {
		ID           string      `json:"id"`
		UserID       string      `json:"user_id"`
		ExternalID   string      `json:"external_id"`
		Name         string      `json:"name"`
		Email        string      `json:"email"`
		Status       string      `json:"status"`
		RejectReason string      `json:"reject_reason"`
		DateCreated  time.Time   `json:"date_created"`
		DateUpdated  time.Time   `json:"date_updated"`
		DateAttended interface{} `json:"date_attended"`
	} `json:"members"`
	DateCreated                time.Time `json:"date_created"`
	DateUpdated                time.Time `json:"date_updated"`
	OnlyAdminsCanCreatePolls   bool      `json:"only_admins_can_create_polls"`
	BlockNewMembershipRequests bool      `json:"block_new_membership_requests"`
	AttendanceGroup            bool      `json:"attendance_group"`
}

// IsGroupAdmin Checks if the userID is an admin of the group
func (g *Group) IsGroupAdmin(userID string) bool {
	if len(g.Members) > 0 {
		for _, member := range g.Members {
			if member.UserID == userID && member.Status == "admin" {
				return true
			}
		}
	}
	return false
}

// GetMembersAsNotificationRecipients Gets members as notification recipients
func (g *Group) GetMembersAsNotificationRecipients(currentUserID string, subMembers []ToMember) []NotificationRecipient {
	var recipients []NotificationRecipient
	if len(g.Members) > 0 {
		if len(subMembers) > 0 {

			//
			// Send notification to the group admins & the sub member list
			//
			userIDmapping := map[string]bool{}
			for _, member := range g.Members {
				if member.UserID != "" && member.UserID != currentUserID && (member.Status == "admin") {
					recipients = append(recipients, NotificationRecipient{
						UserID: member.UserID,
						Name:   member.Name,
					})
					userIDmapping[member.UserID] = true
				}
			}

			for _, toMember := range subMembers {
				if toMember.UserID != "" && toMember.UserID != currentUserID && !userIDmapping[toMember.UserID] {
					recipients = append(recipients, NotificationRecipient{
						UserID: toMember.UserID,
						Name:   toMember.Name,
					})
				}
			}
		} else {

			//
			// Send notification to the group members & group admins
			//
			for _, member := range g.Members {
				if member.UserID != "" && member.UserID != currentUserID && (member.Status == "member" || member.Status == "admin") {
					recipients = append(recipients, NotificationRecipient{
						UserID: member.UserID,
						Name:   member.Name,
					})
				}
			}
		}
	}
	return recipients
}
