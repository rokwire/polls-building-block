/*
 *   Copyright (c) 2020 Board of Trustees of the University of Illinois.
 *   All rights reserved.

 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at

 *   http://www.apache.org/licenses/LICENSE-2.0

 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package core

import (
	"encoding/json"
	"fmt"
	"log"
	"polls/core/model"
	"polls/driven/groups"
	"polls/driven/storage"
)

func (app *Application) getVersion() string {
	return app.version
}

func (app *Application) getPolls(user *model.User, filter model.PollsFilter, filterByToMembers bool) ([]model.Poll, error) {

	var membership *groups.GroupMembership
	if len(filter.GroupIDs) > 0 {
		groupMembership, err := app.groups.GetGroupsMembership(user.Token)
		if err != nil {
			log.Printf("error app.getPolls() - unable to retrieve user groups - %s", err)
			return nil, fmt.Errorf("error app.getPolls() - unable to retrieve user groups - %s", err)
		}

		membership = groupMembership
	}

	return app.storage.GetPolls(user, filter, filterByToMembers, membership)
}

func (app *Application) getPoll(user *model.User, id string) (*model.Poll, error) {
	groupMembership, err := app.groups.GetGroupsMembership(user.Token)
	if err != nil {
		log.Printf("error app.getPoll() - unable to retrieve user groups - %s", err)
		return nil, fmt.Errorf("error app.getPoll() - unable to retrieve user groups - %s", err)
	}

	return app.storage.GetPoll(user, id, true, groupMembership)
}

func (app *Application) createPoll(user *model.User, poll model.Poll) (*model.Poll, error) {
	createdPoll, err := app.storage.CreatePoll(user, poll)
	if err != nil {
		return nil, err
	}

	err = app.notifyNotificationsBBForPoll(user, createdPoll, "polls", "poll_created", fmt.Sprintf("Poll '%s' has been created", createdPoll.Question))
	if err != nil {
		log.Printf("error while sending notification for created poll: %s", err) // dont fail
	}

	return createdPoll, nil
}

func (app *Application) updatePoll(user *model.User, poll model.Poll) (*model.Poll, error) {
	persistedPoll, err := app.storage.GetPoll(user, poll.ID.Hex(), true, nil)
	if err != nil {
		return nil, err
	}

	if persistedPoll == nil || user.Claims.Subject != poll.UserID {
		return nil, fmt.Errorf("only the creator of a poll can update it")
	}

	updatedPoll, err := app.storage.UpdatePoll(user, poll)
	if err != nil {
		return nil, err
	}

	return updatedPoll, nil
}

func (app *Application) deletePoll(user *model.User, id string) error {
	poll, err := app.storage.GetPoll(user, id, true, nil)
	if err != nil {
		return err
	}

	if poll == nil || user.Claims.Subject != poll.UserID {
		if poll.GroupID != nil && len(*poll.GroupID) > 0 {
			group, err := app.groups.GetGroupDetails(*poll.GroupID)
			if err != nil {
				return err
			}
			if group != nil {
				if !group.IsCurrentUserAdmin(user.Claims.Subject) {
					return fmt.Errorf("only the creator of a poll or a group admin can delete it")
				}
			} else {
				return fmt.Errorf("only the creator of a poll or a group admin can delete it")
			}
		} else {
			return fmt.Errorf("only the creator of a poll can delete it")
		}
	}

	err = app.storage.DeletePoll(user, id)
	if err != nil {
		return err
	}

	app.sseServer.NotifyPollForEvent(id, "poll_deleted")
	app.sseServer.ClosePoll(id)

	return nil
}

func (app *Application) startPoll(user *model.User, pollID string) error {
	poll, err := app.storage.GetPoll(user, pollID, true, nil)
	if err != nil {
		return err
	}

	if poll == nil || user.Claims.Subject != poll.UserID {
		return fmt.Errorf("only the creator of a poll can start it")
	}

	if poll != nil {
		if poll.Status != storage.PollStatusStarted {
			poll.Status = storage.PollStatusStarted
			_, err = app.storage.UpdatePoll(user, *poll)
			if err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("error app.startPoll() - poll not found: %s", pollID)
	}

	err = app.notifyNotificationsBBForPoll(user, poll, "polls", "poll_created", fmt.Sprintf("Poll '%s' has been started", poll.Question))
	if err != nil {
		log.Printf("error while sending notification for started poll: %s", err) // dont fail
	}

	app.sseServer.NotifyPollForEvent(pollID, "poll_started")

	return nil
}

func (app *Application) endPoll(user *model.User, pollID string) error {
	poll, err := app.storage.GetPoll(user, pollID, true, nil)
	if err != nil {
		return err
	}

	if poll == nil || user.Claims.Subject != poll.UserID {
		if poll.GroupID != nil && len(*poll.GroupID) > 0 {
			group, err := app.groups.GetGroupDetails(*poll.GroupID)
			if err != nil {
				return err
			}
			if group != nil {
				if !group.IsCurrentUserAdmin(user.Claims.Subject) {
					return fmt.Errorf("only the creator of a poll or a group admin can end it")
				}
			} else {
				return fmt.Errorf("only the creator of a poll or a group admin can end it")
			}
		} else {
			return fmt.Errorf("only the creator of a poll can end it")
		}
	}

	if poll != nil {
		if poll.Status != storage.PollStatusTerminated {
			poll.Status = storage.PollStatusTerminated
			_, err = app.storage.UpdatePoll(user, *poll)
			if err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("error app.startPoll() - poll not found: %s", pollID)
	}

	err = app.notifyNotificationsBBForPoll(user, poll, "polls", "poll_ended", fmt.Sprintf("Poll '%s' has been edned", poll.Question))
	if err != nil {
		log.Printf("error while sending notification for ended poll: %s", err) // dont fail
	}

	app.sseServer.NotifyPollForEvent(pollID, "poll_end")
	app.sseServer.ClosePoll(pollID)

	return nil
}

func (app *Application) notifyNotificationsBBForPoll(user *model.User, poll *model.Poll, topic string, operation string, message string) error {
	recipients, err := app.buildPollNotificationRecipients(user, poll)
	if err != nil {
		log.Printf("error while building recipients for notification for poll: %s", err) // dont fail
		return fmt.Errorf("error while building recipients for notification for poll: %s", err)
	}

	err = app.notifications.SendNotification(
		recipients,
		&topic,
		"Illinois",
		message,
		map[string]string{
			"type":        "poll",
			"operation":   operation,
			"entity_type": "poll",
			"entity_id":   poll.ID.Hex(),
			"entity_name": poll.Question,
		},
	)
	if err != nil {
		log.Printf("error while sending notification for poll end: %s", err) // dont fail
	}

	return nil
}

func (app *Application) votePoll(user *model.User, pollID string, vote model.PollVote) error {
	return app.storage.VotePoll(user, pollID, vote)
}

func (app *Application) subscribeToPoll(user *model.User, pollID string, resultChan chan map[string]interface{}) error {
	app.sseServer.RegisterUserForPoll(user.Claims.Subject, pollID, resultChan)
	return nil
}

func (app *Application) buildPollNotificationRecipients(user *model.User, poll *model.Poll) ([]model.NotificationRecipient, error) {
	if poll.GroupID != nil {
		group, err := app.groups.GetGroupDetails(*poll.GroupID)
		if err != nil {
			return nil, fmt.Errorf("error while retriving group details: %s", err)
		}
		return group.GetMembersAsNotificationRecipients(user.Claims.Subject, poll.ToMembersList), nil
	} else if len(poll.ToMembersList) > 0 {
		return poll.GetPollNotificationRecipients(user.Claims.Subject), nil
	}

	return nil, nil
}

// OnCollectionUpdated callback that indicates the reward types collection is changed
func (app *Application) OnCollectionUpdated(collection string, record map[string]interface{}) {
	if "polls" == collection && record != nil {
		data, err := json.Marshal(record)
		if err != nil {
			log.Printf("Error on Application.OnCollectionUpdated: %s", err)
			return
		}

		if data != nil {
			var poll model.PollNotification
			err = json.Unmarshal(data, &poll)
			if err != nil {
				log.Printf("Error on Application.OnCollectionUpdated: %s", err)
				return
			}

			app.sseServer.NotifyPollUpdate(poll.ID.Hex(), poll)
		}
	}
}
