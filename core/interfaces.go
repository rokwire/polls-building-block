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
	"time"
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
	GetSurvey(user *model.User, id string) (*model.Survey, error)
	CreateSurvey(user *model.User, survey model.Survey, admin bool) (*model.Survey, error)
	UpdateSurvey(user *model.User, survey model.Survey, id string, admin bool) error
	DeleteSurvey(user *model.User, id string, admin bool) error

	//CRUD Survey Response
	GetSurveyResponse(user *model.User, id string) (*model.SurveyResponse, error)
	GetSurveyResponses(user *model.User, surveyIDs []string, surveyTypes []string, startDate *time.Time, endDate *time.Time, limit *int, offset *int) ([]model.SurveyResponse, error)
	CreateSurveyResponse(user *model.User, survey model.Survey) (*model.SurveyResponse, error)
	UpdateSurveyResponse(user *model.User, id string, survey model.Survey) error
	DeleteSurveyResponse(user *model.User, id string) error
	DeleteSurveyResponses(user *model.User, surveyIDs []string, surveyTypes []string, startDate *time.Time, endDate *time.Time) error

	//CRUD Survey Alerts
	GetAlertContacts(user *model.User) ([]model.AlertContact, error)
	GetAlertContact(user *model.User, id string) (*model.AlertContact, error)
	CreateAlertContact(user *model.User, alertContact model.AlertContact) (*model.AlertContact, error)
	UpdateAlertContact(user *model.User, id string, alertContact model.AlertContact) error
	DeleteAlertContact(user *model.User, id string) error
	CreateSurveyAlert(user *model.User, surveyAlert model.SurveyAlert) error

	GetUserData(userID string) (*model.UserDataResponse, error)
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

func (s *servicesImpl) GetSurvey(user *model.User, id string) (*model.Survey, error) {
	return s.app.getSurvey(user, id)
}

func (s *servicesImpl) CreateSurvey(user *model.User, survey model.Survey, admin bool) (*model.Survey, error) {
	return s.app.createSurvey(user, survey, admin)
}

func (s *servicesImpl) UpdateSurvey(user *model.User, survey model.Survey, id string, admin bool) error {
	return s.app.updateSurvey(user, survey, id, admin)
}

func (s *servicesImpl) DeleteSurvey(user *model.User, id string, admin bool) error {
	return s.app.deleteSurvey(user, id, admin)
}

func (s *servicesImpl) DeleteSurveyResponses(user *model.User, surveyIDs []string, surveyTypes []string, startDate *time.Time, endDate *time.Time) error {
	return s.app.deleteSurveyResponses(user, surveyIDs, surveyTypes, startDate, endDate)
}

func (s *servicesImpl) GetSurveyResponse(user *model.User, id string) (*model.SurveyResponse, error) {
	return s.app.getSurveyResponse(user, id)
}

func (s *servicesImpl) GetSurveyResponses(user *model.User, surveyIDs []string, surveyTypes []string, startDate *time.Time, endDate *time.Time, limit *int, offset *int) ([]model.SurveyResponse, error) {
	return s.app.getSurveyResponses(user, surveyIDs, surveyTypes, startDate, endDate, limit, offset)
}

func (s *servicesImpl) CreateSurveyResponse(user *model.User, survey model.Survey) (*model.SurveyResponse, error) {
	return s.app.createSurveyResponse(user, survey)
}

func (s *servicesImpl) UpdateSurveyResponse(user *model.User, id string, survey model.Survey) error {
	return s.app.updateSurveyResponse(user, id, survey)
}

func (s *servicesImpl) DeleteSurveyResponse(user *model.User, id string) error {
	return s.app.deleteSurveyResponse(user, id)
}

func (s *servicesImpl) GetAlertContacts(user *model.User) ([]model.AlertContact, error) {
	return s.app.getAlertContacts(user)
}

func (s *servicesImpl) GetAlertContact(user *model.User, id string) (*model.AlertContact, error) {
	return s.app.getAlertContact(user, id)
}

func (s *servicesImpl) CreateAlertContact(user *model.User, alertContact model.AlertContact) (*model.AlertContact, error) {
	return s.app.createAlertContact(user, alertContact)
}

func (s *servicesImpl) UpdateAlertContact(user *model.User, id string, alertContact model.AlertContact) error {
	return s.app.updateAlertContact(user, id, alertContact)
}

func (s *servicesImpl) DeleteAlertContact(user *model.User, id string) error {
	return s.app.deleteAlertContact(user, id)
}

func (s *servicesImpl) CreateSurveyAlert(user *model.User, surveyAlert model.SurveyAlert) error {
	return s.app.createSurveyAlert(user, surveyAlert)
}

func (s *servicesImpl) GetUserData(userID string) (*model.UserDataResponse, error) {
	return s.app.getUserData(userID)
}

// Storage is used by core to storage data - DB storage adapter, file storage adapter etc
type Storage interface {
	GetPolls(user *model.User, filter model.PollsFilter, filterByToMembers bool, membership *groups.GroupMembership) ([]model.Poll, error)
	GetPoll(user *model.User, id string, filterByToMembers bool, membership *groups.GroupMembership) (*model.Poll, error)
	CreatePoll(user *model.User, poll model.Poll) (*model.Poll, error)
	UpdatePoll(user *model.User, poll model.Poll) (*model.Poll, error)

	DeletePoll(user *model.User, id string) error

	VotePoll(user *model.User, pollID string, vote model.PollVote) error
	DeletePollsWithIDs(orgID string, accountsIDs []string) error

	SetListener(listener storage.CollectionListener)

	GetSurvey(user *model.User, id string) (*model.Survey, error)
	CreateSurvey(survey model.Survey) (*model.Survey, error)
	UpdateSurvey(user *model.User, survey model.Survey, admin bool) error
	DeleteSurvey(user *model.User, id string, admin bool) error
	DeleteSurveysWithIDs(appID string, orgID string, accountsIDs []string) error

	GetSurveyResponse(user *model.User, id string) (*model.SurveyResponse, error)
	GetSurveyResponses(user *model.User, surveyIDs []string, surveyTypes []string, startDate *time.Time, endDate *time.Time, limit *int, offset *int) ([]model.SurveyResponse, error)
	CreateSurveyResponse(surveyResponse model.SurveyResponse) (*model.SurveyResponse, error)
	UpdateSurveyResponse(user *model.User, id string, surveyResponse model.Survey) error
	DeleteSurveyResponse(user *model.User, id string) error
	DeleteSurveyResponses(user *model.User, surveyIDs []string, surveyTypes []string, startDate *time.Time, endDate *time.Time) error
	DeleteSurveyResponsesWithIDs(appID string, orgID string, accountsIDs []string) error

	GetAlertContacts(user *model.User) ([]model.AlertContact, error)
	GetAlertContact(user *model.User, id string) (*model.AlertContact, error)
	CreateAlertContact(alertContact model.AlertContact) (*model.AlertContact, error)
	UpdateAlertContact(user *model.User, id string, alertContact model.AlertContact) error
	DeleteAlertContact(user *model.User, id string) error
	GetAlertContactsByKey(key string, user *model.User) ([]model.AlertContact, error)
}

// Core exposes Core APIs for the driver adapters
type Core interface {
	LoadDeletedMemberships() ([]model.DeletedUserData, error)
}
