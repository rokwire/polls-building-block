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
	"github.com/rokwire/core-auth-library-go/tokenauth"
	"log"
	"polls/core/model"
)

func (app *Application) getVersion() string {
	return app.version
}

func (app *Application) getPolls(user *tokenauth.Claims, filter model.PollsFilter, filterByToMembers bool) ([]model.Poll, error) {
	return app.storage.GetPolls(user, filter, filterByToMembers)
}

func (app *Application) getPoll(user *tokenauth.Claims, id string) (*model.Poll, error) {
	return app.storage.GetPoll(user, id)
}

func (app *Application) createPoll(user *tokenauth.Claims, poll model.Poll) (*model.Poll, error) {
	return app.storage.CreatePoll(user, poll)
}

func (app *Application) updatePoll(user *tokenauth.Claims, poll model.Poll) (*model.Poll, error) {
	updatedPoll, err := app.storage.UpdatePoll(user, poll)
	if err != nil {
		return nil, err
	}

	return updatedPoll, nil
}

func (app *Application) deletePoll(user *tokenauth.Claims, id string) error {
	err := app.storage.DeletePoll(user, id)
	if err != nil {
		return err
	}

	app.sseServer.NotifyPollForEvent(id, "poll_deleted")
	app.sseServer.ClosePoll(id)

	return nil
}

func (app *Application) startPoll(user *tokenauth.Claims, pollID string) error {
	err := app.storage.StartPoll(user, pollID)
	if err != nil {
		return err
	}

	app.sseServer.NotifyPollForEvent(pollID, "poll_started")

	return nil
}

func (app *Application) endPoll(user *tokenauth.Claims, pollID string) error {
	err := app.storage.EndPoll(user, pollID)
	if err != nil {
		return err
	}

	app.sseServer.NotifyPollForEvent(pollID, "poll_end")
	app.sseServer.ClosePoll(pollID)

	return nil
}

func (app *Application) votePoll(user *tokenauth.Claims, pollID string, vote model.PollVote) error {
	return app.storage.VotePoll(user, pollID, vote)
}

func (app *Application) subscribeToPoll(user *tokenauth.Claims, pollID string, resultChan chan map[string]interface{}, closeChan chan interface{}) error {
	app.sseServer.RegisterUserForPoll(user.Subject, pollID, resultChan, closeChan)
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
