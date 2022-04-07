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
	"polls/core/model"
	"polls/driven/storage"
)

// Services exposes APIs for the driver adapters
type Services interface {
	GetVersion() string

	// CRUD
	GetPolls(IDs []string, userID *string, offset *int64, limit *int64, order *string) ([]model.Poll, error)
	GetPoll(id string) (*model.Poll, error)
	CreatePoll(poll model.Poll) (*model.Poll, error)
	UpdatePoll(poll model.Poll) (*model.Poll, error)

	VotePoll(pollID string, vote model.PollVote) error
	StartPoll(pollID string) error
	EndPoll(pollID string) error

	DeletePoll(id string) error
}

type servicesImpl struct {
	app *Application
}

func (s *servicesImpl) GetVersion() string {
	return s.app.getVersion()
}

func (s *servicesImpl) GetPolls(IDs []string, userID *string, offset *int64, limit *int64, order *string) ([]model.Poll, error) {
	return s.app.getPolls(IDs, userID, offset, limit, order)
}

func (s *servicesImpl) GetPoll(id string) (*model.Poll, error) {
	return s.app.getPoll(id)
}

func (s *servicesImpl) CreatePoll(poll model.Poll) (*model.Poll, error) {
	return s.app.createPoll(poll)
}

func (s *servicesImpl) UpdatePoll(poll model.Poll) (*model.Poll, error) {
	return s.app.updatePoll(poll)
}

func (s *servicesImpl) DeletePoll(id string) error {
	return s.app.deletePoll(id)
}

func (s *servicesImpl) StartPoll(pollID string) error {
	return s.app.startPoll(pollID)
}

func (s *servicesImpl) EndPoll(pollID string) error {
	return s.app.endPoll(pollID)
}

func (s *servicesImpl) VotePoll(pollID string, vote model.PollVote) error {
	return s.app.votePoll(pollID, vote)
}

// Storage is used by core to storage data - DB storage adapter, file storage adapter etc
type Storage interface {
	GetPolls(IDs []string, userID *string, offset *int64, limit *int64, order *string) ([]model.Poll, error)
	GetPoll(id string) (*model.Poll, error)
	CreatePoll(poll model.Poll) (*model.Poll, error)
	UpdatePoll(poll model.Poll) (*model.Poll, error)
	StartPoll(pollID string) error
	EndPoll(pollID string) error
	DeletePoll(id string) error

	VotePoll(pollID string, vote model.PollVote) error

	SetListener(listener storage.CollectionListener)
}
