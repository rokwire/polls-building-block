// Copyright 2022 Board of Trustees of the University of Illinois.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import "polls/core/model"

// SSEClient struct
type SSEClient struct {
	pollID     string
	userID     string
	resultChan chan map[string]interface{}
}

// SSEServer struct
type SSEServer struct {
	PollClientsMapping map[string][]SSEClient
}

// NewSSEServer new instance
func NewSSEServer() *SSEServer {
	return &SSEServer{PollClientsMapping: map[string][]SSEClient{}}
}

// RegisterUserForPoll registers a user for a poll updates
func (s *SSEServer) RegisterUserForPoll(userID, pollID string, resultChan chan map[string]interface{}) {
	var list []SSEClient
	if val, ok := s.PollClientsMapping[pollID]; ok {
		list = val
	} else {
		list = []SSEClient{}
	}

	list = append(list, SSEClient{pollID: pollID, userID: userID, resultChan: resultChan})
	s.PollClientsMapping[pollID] = list
}

// UnregisterUser unregisters user for poll updates
func (s *SSEServer) UnregisterUser(userID string, pollID string) {
	if clients, ok := s.PollClientsMapping[pollID]; ok {
		if len(clients) > 0 {
			var newList []SSEClient
			for _, client := range clients {
				if client.userID == userID {
					close(client.resultChan)
				} else {
					newList = append(newList, client)
				}
			}
			s.PollClientsMapping[userID] = newList
		}
	}
}

// ClosePoll notifies all subscribers the poll is closed and remove the client
func (s *SSEServer) ClosePoll(pollID string) {
	if clients, ok := s.PollClientsMapping[pollID]; ok {
		if len(clients) > 0 {
			for _, client := range clients {
				close(client.resultChan)
			}
			delete(s.PollClientsMapping, pollID)
		}
	}
}

// NotifyPollForEvent notifies all subscribers for changed poll
func (s *SSEServer) NotifyPollForEvent(pollID string, eventType string) {
	if list, ok := s.PollClientsMapping[pollID]; ok {
		if len(list) > 0 {
			for _, client := range list {
				client.resultChan <- map[string]interface{}{
					"poll_id":    pollID,
					"event_type": eventType,
				}
			}
		}
	}
}

// NotifyPollUpdate notifies all subscribers for changed poll
func (s *SSEServer) NotifyPollUpdate(pollID string, poll model.PollNotification) {
	if list, ok := s.PollClientsMapping[pollID]; ok {
		if len(list) > 0 {
			for _, client := range list {
				go func() {
					client.resultChan <- map[string]interface{}{
						"poll_id":    pollID,
						"event_type": "poll_updated",
						"result":     poll.ToPollResult(client.userID).Results,
					}
				}()
			}
		}
	}
}
