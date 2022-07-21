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

import (
	"polls/core/model"
	"polls/driven/storage"

	"github.com/rokwire/core-auth-library-go/tokenauth"
)

// Services exposes APIs for the driver adapters
type Services interface {
	GetVersion() string

	// CRUD
	GetPolls(user *tokenauth.Claims, filter model.PollsFilter, filterByToMembers bool) ([]model.Poll, error)
	GetPoll(user *tokenauth.Claims, id string) (*model.Poll, error)
	CreatePoll(user *tokenauth.Claims, poll model.Poll) (*model.Poll, error)
	UpdatePoll(user *tokenauth.Claims, poll model.Poll) (*model.Poll, error)
	DeletePoll(user *tokenauth.Claims, id string) error

	VotePoll(user *tokenauth.Claims, pollID string, vote model.PollVote) error
	StartPoll(user *tokenauth.Claims, pollID string) error
	EndPoll(user *tokenauth.Claims, pollID string) error

	SubscribeToPoll(user *tokenauth.Claims, pollID string, resultChan chan map[string]interface{}) error
}

type servicesImpl struct {
	app *Application
}

func (s *servicesImpl) GetVersion() string {
	return s.app.getVersion()
}

func (s *servicesImpl) GetPolls(user *tokenauth.Claims, filter model.PollsFilter, filterByToMembers bool) ([]model.Poll, error) {
	return s.app.getPolls(user, filter, filterByToMembers)
}

func (s *servicesImpl) GetPoll(user *tokenauth.Claims, id string) (*model.Poll, error) {
	return s.app.getPoll(user, id)
}

func (s *servicesImpl) CreatePoll(user *tokenauth.Claims, poll model.Poll) (*model.Poll, error) {
	return s.app.createPoll(user, poll)
}

func (s *servicesImpl) UpdatePoll(user *tokenauth.Claims, poll model.Poll) (*model.Poll, error) {
	return s.app.updatePoll(user, poll)
}

func (s *servicesImpl) DeletePoll(user *tokenauth.Claims, id string) error {
	return s.app.deletePoll(user, id)
}

func (s *servicesImpl) StartPoll(user *tokenauth.Claims, pollID string) error {
	return s.app.startPoll(user, pollID)
}

func (s *servicesImpl) EndPoll(user *tokenauth.Claims, pollID string) error {
	return s.app.endPoll(user, pollID)
}

func (s *servicesImpl) VotePoll(user *tokenauth.Claims, pollID string, vote model.PollVote) error {
	return s.app.votePoll(user, pollID, vote)
}

func (s *servicesImpl) SubscribeToPoll(user *tokenauth.Claims, pollID string, resultChan chan map[string]interface{}) error {
	return s.app.subscribeToPoll(user, pollID, resultChan)
}

// Storage is used by core to storage data - DB storage adapter, file storage adapter etc
type Storage interface {
	GetPolls(user *tokenauth.Claims, filter model.PollsFilter, filterByToMembers bool) ([]model.Poll, error)
	GetPoll(user *tokenauth.Claims, id string) (*model.Poll, error)
	CreatePoll(user *tokenauth.Claims, poll model.Poll) (*model.Poll, error)
	UpdatePoll(user *tokenauth.Claims, poll model.Poll) (*model.Poll, error)
	StartPoll(user *tokenauth.Claims, pollID string) error
	EndPoll(user *tokenauth.Claims, pollID string) error
	DeletePoll(user *tokenauth.Claims, id string) error

	VotePoll(user *tokenauth.Claims, pollID string, vote model.PollVote) error

	SetListener(listener storage.CollectionListener)
}
