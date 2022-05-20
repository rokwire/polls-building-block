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
	"polls/driven/storage"
)

func (app *Application) getVersion() string {
	return app.version
}

func (app *Application) getPolls(user *model.User, filter model.PollsFilter, filterByToMembers bool) ([]model.Poll, error) {
	storageFilter := filter.ToPollsInnerFilter()

	if filter.GroupPolls != nil && *filter.GroupPolls {
		groupIDs, err := app.groups.GetGroupsMembership(user.Token)
		if err != nil {
			log.Printf("error app.getPolls() - unable to retrieve user groups - %s", err)
			return nil, fmt.Errorf("error app.getPolls() - unable to retrieve user groups - %s", err)
		}

		if len(groupIDs) > 0 {
			storageFilter.GroupIDs = groupIDs
		}
	}

	return app.storage.GetPolls(user, storageFilter, filterByToMembers)
}

func (app *Application) getPoll(user *model.User, id string) (*model.Poll, error) {
	return app.storage.GetPoll(user, id)
}

func (app *Application) createPoll(user *model.User, poll model.Poll) (*model.Poll, error) {
	return app.storage.CreatePoll(user, poll)
}

func (app *Application) updatePoll(user *model.User, poll model.Poll) (*model.Poll, error) {
	updatedPoll, err := app.storage.UpdatePoll(user, poll)
	if err != nil {
		return nil, err
	}

	return updatedPoll, nil
}

func (app *Application) deletePoll(user *model.User, id string) error {
	err := app.storage.DeletePoll(user, id)
	if err != nil {
		return err
	}

	app.sseServer.NotifyPollForEvent(id, "poll_deleted")
	app.sseServer.ClosePoll(id)

	return nil
}

func (app *Application) startPoll(user *model.User, pollID string) error {
	poll, err := app.storage.GetPoll(user, pollID)
	if err != nil {
		return err
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

	topic := "polls"
	err = app.notifications.SendNotification(
		nil,
		&topic,
		"Illinois",
		fmt.Sprintf("Poll '%s' has been started", poll.Question),
		map[string]string{
			"type":        "poll",
			"operation":   "poll_created",
			"entity_type": "group",
			"entity_id":   poll.ID.Hex(),
			"entity_name": poll.Question,
		},
	)
	if err != nil {
		log.Printf("error while sending notification for new event: %s", err) // dont fail
	}

	app.sseServer.NotifyPollForEvent(pollID, "poll_started")

	return nil
}

func (app *Application) endPoll(user *model.User, pollID string) error {
	poll, err := app.storage.GetPoll(user, pollID)
	if err != nil {
		return err
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

	topic := "polls"
	err = app.notifications.SendNotification(
		nil,
		&topic,
		"Illinois",
		fmt.Sprintf("Poll '%s' has been edned", poll.Question),
		map[string]string{
			"type":        "poll",
			"operation":   "poll_ended",
			"entity_type": "group",
			"entity_id":   poll.ID.Hex(),
			"entity_name": poll.Question,
		},
	)
	if err != nil {
		log.Printf("error while sending notification for poll end: %s", err) // dont fail
	}

	app.sseServer.NotifyPollForEvent(pollID, "poll_end")
	app.sseServer.ClosePoll(pollID)

	return nil
}

func (app *Application) votePoll(user *model.User, pollID string, vote model.PollVote) error {
	return app.storage.VotePoll(user, pollID, vote)
}

func (app *Application) subscribeToPoll(user *model.User, pollID string, resultChan chan map[string]interface{}) error {
	app.sseServer.RegisterUserForPoll(user.Claims.Subject, pollID, resultChan)
	return nil
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
