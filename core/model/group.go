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
	Description     string `json:"description"`
	ImageURL        string `json:"image_url"`
	WebURL          string `json:"web_url"`
	Tags            string `json:"tags"`
	CurrentMember   *struct {
		ID                       string `json:"id"`
		ClientID                 string `json:"client_id"`
		GroupID                  string `json:"group_id"`
		UserID                   string `json:"user_id"`
		ExternalID               string `json:"external_id"`
		Name                     string `json:"name"`
		NetID                    string `json:"net_id"`
		Email                    string `json:"email"`
		PhotoURL                 string `json:"photo_url"`
		Status                   string `json:"status"`
		Admin                    bool   `json:"admin"`
		RejectReason             string `json:"reject_reason"`
		NotificationsPreferences struct {
			OverridePreferences bool `json:"override_preferences"`
			AllMute             bool `json:"all_mute"`
			InvitationsMute     bool `json:"invitations_mute"`
			PostsMute           bool `json:"posts_mute"`
			EventsMute          bool `json:"events_mute"`
			PollsMute           bool `json:"polls_mute"`
		} `json:"notifications_preferences"`
		DateCreated  time.Time  `json:"date_created"`
		DateUpdated  *time.Time `json:"date_updated"`
		DateAttended *time.Time `json:"date_attended"`
	} `json:"current_member"`
	Stats struct {
		TotalCount      int `json:"total_count"`
		AdminsCount     int `json:"admins_count"`
		MemberCount     int `json:"member_count"`
		PendingCount    int `json:"pending_count"`
		RejectedCount   int `json:"rejected_count"`
		AttendanceCount int `json:"attendance_count"`
	} `json:"stats"`
	DateCreated                time.Time  `json:"date_created"`
	DateUpdated                *time.Time `json:"date_updated"`
	AuthmanEnabled             bool       `json:"authman_enabled"`
	AuthmanGroup               string     `json:"authman_group"`
	OnlyAdminsCanCreatePolls   bool       `json:"only_admins_can_create_polls"`
	CanJoinAutomatically       bool       `json:"can_join_automatically"`
	BlockNewMembershipRequests bool       `json:"block_new_membership_requests"`
	AttendanceGroup            bool       `json:"attendance_group"`
}

// IsCurrentUserAdmin checks if the user is a group admin
func (g *Group) IsCurrentUserAdmin(currentUserID string) bool {
	if g.CurrentMember != nil {
		if g.CurrentMember.UserID == currentUserID && g.CurrentMember.Status == "admin" {
			return true
		}
	}
	return false
}
