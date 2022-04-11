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
	"github.com/rokwire/core-auth-library-go/tokenauth"
	"polls/core/model"
)

func (app *Application) getVersion() string {
	return app.version
}

func (app *Application) getPolls(user *tokenauth.Claims, IDs []string, userID *string, offset *int64, limit *int64, order *string, filterByToMembers bool) ([]model.Poll, error) {
	return app.storage.GetPolls(user, IDs, userID, offset, limit, order, filterByToMembers)
}

func (app *Application) getPoll(user *tokenauth.Claims, id string) (*model.Poll, error) {
	return app.storage.GetPoll(user, id)
}

func (app *Application) createPoll(user *tokenauth.Claims, poll model.Poll) (*model.Poll, error) {
	return app.storage.CreatePoll(user, poll)
}

func (app *Application) updatePoll(user *tokenauth.Claims, poll model.Poll) (*model.Poll, error) {
	return app.storage.UpdatePoll(user, poll)
}

func (app *Application) deletePoll(user *tokenauth.Claims, id string) error {
	return app.storage.DeletePoll(user, id)
}

func (app *Application) startPoll(user *tokenauth.Claims, pollID string) error {
	return app.storage.StartPoll(user, pollID)
}

func (app *Application) endPoll(user *tokenauth.Claims, pollID string) error {
	return app.storage.EndPoll(user, pollID)
}

func (app *Application) votePoll(user *tokenauth.Claims, pollID string, vote model.PollVote) error {
	return app.storage.VotePoll(user, pollID, vote)
}

// OnCollectionUpdated callback that indicates the reward types collection is changed
func (app *Application) OnCollectionUpdated(name string) {

}
