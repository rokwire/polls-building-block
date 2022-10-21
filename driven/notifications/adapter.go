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
	internalAPIKey string
	baseURL        string
}

// NewNotificationsAdapter creates a new Notifications BB adapter instance
func NewNotificationsAdapter(config *model.Config) *Adapter {
	return &Adapter{internalAPIKey: config.InternalAPIKey, baseURL: config.NotificationsHost}
}

// SendNotification Sends a direct notification trough Notifications BB
func (a *Adapter) SendNotification(notification model.NotificationMessage) {
	go a.sendNotification(notification)
}

// SendNotification sends notification to a user
func (a *Adapter) sendNotification(notification model.NotificationMessage) {
	if len(notification.Recipients) > 0 && notification.Subject != "" && notification.Body != "" {
		url := fmt.Sprintf("%s/api/int/message", a.baseURL)

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
