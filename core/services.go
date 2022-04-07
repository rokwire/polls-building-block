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

import "polls/core/model"

func (app *Application) getVersion() string {
	return app.version
}

func (app *Application) getPolls(IDs []string, userID *string, offset *int64, limit *int64, order *string) ([]model.Poll, error) {
	return app.storage.GetPolls(IDs, userID, offset, limit, order)
}

func (app *Application) getPoll(id string) (*model.Poll, error) {
	return app.storage.GetPoll(id)
}

func (app *Application) createPoll(poll model.Poll) (*model.Poll, error) {
	return app.storage.CreatePoll(poll)
}

func (app *Application) updatePoll(poll model.Poll) (*model.Poll, error) {
	return app.storage.UpdatePoll(poll)
}

func (app *Application) deletePoll(id string) error {
	return app.storage.DeletePoll(id)
}

func (app *Application) startPoll(pollID string) error {
	return app.storage.StartPoll(pollID)
}

func (app *Application) endPoll(pollID string) error {
	return app.storage.EndPoll(pollID)
}

func (app *Application) votePoll(pollID string, vote model.PollVote) error {
	return app.storage.VotePoll(pollID, vote)
}

// OnCollectionUpdated callback that indicates the reward types collection is changed
func (app *Application) OnCollectionUpdated(name string) {

}
