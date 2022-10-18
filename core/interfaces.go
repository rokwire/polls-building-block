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
	"polls/driven/groups"
	"polls/driven/storage"
)

// Services exposes APIs for the driver adapters
type Services interface {
	GetVersion() string

	// CRUD Polls
	GetPolls(user *model.User, filter model.PollsFilter, filterByToMembers bool) ([]model.Poll, error)
	GetPoll(user *model.User, id string) (*model.Poll, error)
	CreatePoll(user *model.User, poll model.Poll) (*model.Poll, error)
	UpdatePoll(user *model.User, poll model.Poll) (*model.Poll, error)
	DeletePoll(user *model.User, id string) error

	VotePoll(user *model.User, pollID string, vote model.PollVote) error
	StartPoll(user *model.User, pollID string) error
	EndPoll(user *model.User, pollID string) error

	SubscribeToPoll(user *model.User, pollID string, resultChan chan map[string]interface{}) error

	//CRUD Surveys
	GetSurvey(id string) (*model.Survey, error)
	CreateSurvey(user *model.User, survey model.Survey) (*model.Survey, error)
	UpdateSurvey(user *model.User, survey model.Survey, id string) error
	DeleteSurvey(user *model.User, id string) error
}

type servicesImpl struct {
	app *Application
}

func (s *servicesImpl) GetVersion() string {
	return s.app.getVersion()
}

func (s *servicesImpl) GetPolls(user *model.User, filter model.PollsFilter, filterByToMembers bool) ([]model.Poll, error) {
	return s.app.getPolls(user, filter, filterByToMembers)
}

func (s *servicesImpl) GetPoll(user *model.User, id string) (*model.Poll, error) {
	return s.app.getPoll(user, id)
}

func (s *servicesImpl) CreatePoll(user *model.User, poll model.Poll) (*model.Poll, error) {
	return s.app.createPoll(user, poll)
}

func (s *servicesImpl) UpdatePoll(user *model.User, poll model.Poll) (*model.Poll, error) {
	return s.app.updatePoll(user, poll)
}

func (s *servicesImpl) DeletePoll(user *model.User, id string) error {
	return s.app.deletePoll(user, id)
}

func (s *servicesImpl) StartPoll(user *model.User, pollID string) error {
	return s.app.startPoll(user, pollID)
}

func (s *servicesImpl) EndPoll(user *model.User, pollID string) error {
	return s.app.endPoll(user, pollID)
}

func (s *servicesImpl) VotePoll(user *model.User, pollID string, vote model.PollVote) error {
	return s.app.votePoll(user, pollID, vote)
}

func (s *servicesImpl) SubscribeToPoll(user *model.User, pollID string, resultChan chan map[string]interface{}) error {
	return s.app.subscribeToPoll(user, pollID, resultChan)
}

func (s *servicesImpl) GetSurvey(id string) (*model.Survey, error) {
	return s.app.getSurvey(id)
}

func (s *servicesImpl) CreateSurvey(user *model.User, survey model.Survey) (*model.Survey, error) {
	return s.app.createSurvey(user, survey)
}

func (s *servicesImpl) UpdateSurvey(user *model.User, survey model.Survey, id string) error {
	return s.app.updateSurvey(user, survey, id)
}

func (s *servicesImpl) DeleteSurvey(user *model.User, id string) error {
	return s.app.deleteSurvey(user, id)
}

// Storage is used by core to storage data - DB storage adapter, file storage adapter etc
type Storage interface {
	GetPolls(user *model.User, filter model.PollsFilter, filterByToMembers bool, membership *groups.GroupMembership) ([]model.Poll, error)
	GetPoll(user *model.User, id string, filterByToMembers bool, membership *groups.GroupMembership) (*model.Poll, error)
	CreatePoll(user *model.User, poll model.Poll) (*model.Poll, error)
	UpdatePoll(user *model.User, poll model.Poll) (*model.Poll, error)

	DeletePoll(user *model.User, id string) error

	VotePoll(user *model.User, pollID string, vote model.PollVote) error

	SetListener(listener storage.CollectionListener)

	GetSurvey(id string) (*model.Survey, error)
	CreateSurvey(survey model.Survey) (*model.Survey, error)
	UpdateSurvey(user *model.User, survey model.Survey) error
	DeleteSurvey(user *model.User, id string) error
}
