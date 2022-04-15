package core

import "polls/core/model"

// SSEClient struct
type SSEClient struct {
	pollID     string
	userID     string
	resultChan chan map[string]interface{}
	closeChan  chan interface{}
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
func (s *SSEServer) RegisterUserForPoll(userID, pollID string, resultChan chan map[string]interface{}, closeChan chan interface{}) {
	var list []SSEClient
	if val, ok := s.PollClientsMapping[pollID]; ok {
		list = val
	} else {
		list = []SSEClient{}
	}

	list = append(list, SSEClient{pollID: pollID, userID: userID, resultChan: resultChan, closeChan: closeChan})
	s.PollClientsMapping[pollID] = list
}

// UnregisterUser unregisters user for poll updates
func (s *SSEServer) UnregisterUser(userID string, pollID string) {
	if clients, ok := s.PollClientsMapping[pollID]; ok {
		if len(clients) > 0 {
			var newList []SSEClient
			for _, client := range clients {
				if client.userID == userID {
					client.closeChan <- 1
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
				client.closeChan <- 1
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
				go func() {
					client.resultChan <- map[string]interface{}{
						"poll_id":    pollID,
						"event_type": eventType,
					}
				}()
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
