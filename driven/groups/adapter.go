package groups

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"polls/core/model"
)

// Adapter groups adapter
type Adapter struct {
	internalAPIKey string
	baseURL        string
}

// NewGroupsAdapter creates a new Groups BB adapter instance
func NewGroupsAdapter(config *model.Config) *Adapter {
	return &Adapter{internalAPIKey: config.InternalAPIKey, baseURL: config.GroupsHost}
}

type userGroup struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Privacy          string `json:"privacy"`
	MembershipStatus string `json:"membership_status"`
}

// GetGroupsMembership retrieves all groups that a user is a member
func (a *Adapter) GetGroupsMembership(userToken string) ([]string, error) {
	if userToken != "" {

		url := fmt.Sprintf("%s/api/user/group-memberships", a.baseURL)
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userToken))
		if err != nil {
			log.Printf("error GetGroupsMembership: request - %s", err)
			return nil, fmt.Errorf("error GetGroupsMembership: request - %s", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("error GetGroupsMembership: request - %s", err)
			return nil, fmt.Errorf("error GetGroupsMembership: request - %s", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			errorBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("error GetGroupsMembership: request - %s", err)
				return nil, fmt.Errorf("error GetGroupsMembership: request - %s", err)
			}

			log.Printf("error GetGroupsMembership: request - %d. Error: %s, Body: %s", resp.StatusCode, err, string(errorBody))
			return nil, fmt.Errorf("error GetGroupsMembership: request - %d. Error: %s, Body: %s", resp.StatusCode, err, string(errorBody))
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("error GetGroupsMembership: request - %s", err)
			return nil, fmt.Errorf("error GetGroupsMembership: request - %s", err)
		}

		var groups []userGroup
		err = json.Unmarshal(data, &groups)
		if err != nil {
			log.Printf("error GetGroupsMembership: request - %s", err)
			return nil, fmt.Errorf("error GetGroupsMembership: request - %s", err)
		}

		ids := []string{}
		if len(groups) > 0 {
			for _, group := range groups {
				if group.MembershipStatus == "member" || group.MembershipStatus == "admin" {
					ids = append(ids, group.ID)
				}
			}
		}

		return ids, nil
	}
	return nil, nil
}
