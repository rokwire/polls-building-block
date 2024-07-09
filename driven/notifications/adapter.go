package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"polls/core/model"
)

// Adapter implements the Notifications interface
type Adapter struct {
	baseURL        string
	internalAPIKey string
	appID          string
	orgID          string
}

// NewNotificationsAdapter creates a new Notifications BB adapter instance
func NewNotificationsAdapter(notificationHost string, internalAPIKey string, appID string, orgID string) *Adapter {
	return &Adapter{baseURL: notificationHost, internalAPIKey: internalAPIKey, appID: appID, orgID: orgID}
}

// SendNotification Sends a direct notification trough Notifications BB
func (a *Adapter) SendNotification(notification model.NotificationMessage) {
	go a.sendNotification(notification)
}

// SendNotification sends notification to a user
func (a *Adapter) sendNotification(notification model.NotificationMessage) {
	if notification.Message.Subject != "" && notification.Message.Body != "" {
		url := fmt.Sprintf("%s/api/int/v2/message", a.baseURL)

		bodyBytes, err := json.Marshal(notification)
		if err != nil {
			log.Printf("error creating notification request - %s", err)
			return
		}

		client := &http.Client{}
		req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
		if err != nil {
			log.Printf("error creating load user data request - %s", err)
			return
		}
		req.Header.Set("INTERNAL-API-KEY", a.internalAPIKey)

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("error loading user data - %s", err)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Printf("error with response code - %d", resp.StatusCode)
		}
	}
}

// SendMail sends email to a user
func (a *Adapter) SendMail(user *model.User, toEmail string, subject string, body string) {
	go a.sendMail(user, toEmail, subject, body)
}

func (a *Adapter) sendMail(user *model.User, toEmail string, subject string, body string) error {
	if len(toEmail) > 0 && len(subject) > 0 && len(body) > 0 {
		url := fmt.Sprintf("%s/api/int/mail", a.baseURL)

		bodyData := map[string]interface{}{
			"to_mail": toEmail,
			"subject": subject,
			"body":    body,
		}
		bodyBytes, err := json.Marshal(bodyData)
		if err != nil {
			log.Printf("error creating notification request - %s", err)
			return err
		}

		client := &http.Client{}
		req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
		if err != nil {
			log.Printf("error creating load user data request - %s", err)
			return err
		}
		req.Header.Set("INTERNAL-API-KEY", a.internalAPIKey)

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("error sending an email - %s", err)
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Printf("error with response code - %d", resp.StatusCode)
			return fmt.Errorf("error with response code != 200")
		}
	}
	return nil
}
