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

package storage

import (
	"fmt"
	"log"
	"polls/core/model"
	"polls/driven/groups"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// PollStatusCreated status created
	PollStatusCreated = "created"

	// PollStatusStarted status started
	PollStatusStarted = "started"

	// PollStatusTerminated status terminated
	PollStatusTerminated = "terminated"

	settingsKey   = "stadium"
	eventInterval = 100 * time.Millisecond
)

// Adapter implements the Storage interface
type Adapter struct {
	db     *database
	config *model.Config
}

// Start starts the storage
func (sa *Adapter) Start() error {
	err := sa.db.start()
	if err != nil {
		return err
	}

	err = sa.applyMultiTenancy()
	return err
}

// NewStorageAdapter creates a new storage adapter instance
func NewStorageAdapter(config *model.Config) *Adapter {
	timeout, err := strconv.Atoi(config.MongoTimeout)
	if err != nil {
		log.Println("Set default timeout - 500")
		timeout = 500
	}
	timeoutMS := time.Millisecond * time.Duration(timeout)

	db := &database{mongoDBAuth: config.MongoDBAuth, mongoDBName: config.MongoDBName, mongoTimeout: timeoutMS}
	return &Adapter{db: db, config: config}
}

// GetPolls retrieves all polls with an ability to filter
func (sa *Adapter) GetPolls(user *model.User, filter model.PollsFilter, filterByToMembers bool, membership *groups.GroupMembership) ([]model.Poll, error) {
	mongoFilter := bson.D{}
	if user != nil {
		mongoFilter = append(mongoFilter, primitive.E{Key: "org_id", Value: user.Claims.OrgID})
	}

	if len(filter.PollIDs) > 0 {
		reconstructedIDs := []primitive.ObjectID{}
		for _, id := range filter.PollIDs {
			if objID, err := primitive.ObjectIDFromHex(id); err == nil {
				reconstructedIDs = append(reconstructedIDs, objID)
			}
		}
		mongoFilter = append(mongoFilter, primitive.E{Key: "_id", Value: bson.M{"$in": reconstructedIDs}})
	}

	if filter.MyPolls != nil && *filter.MyPolls == true && filter.RespondedPolls != nil && *filter.RespondedPolls == true {
		mongoFilter = append(mongoFilter, primitive.E{Key: "$or", Value: []primitive.M{
			{"poll.userid": user.Claims.Subject},
			{"responses.userid": user.Claims.Subject},
		}})
	} else {
		if filter.MyPolls != nil && *filter.MyPolls == true {
			mongoFilter = append(mongoFilter, primitive.E{Key: "poll.userid", Value: user.Claims.Subject})
		}

		if filter.RespondedPolls != nil && *filter.RespondedPolls == true {
			mongoFilter = append(mongoFilter, primitive.E{Key: "responses.userid", Value: user.Claims.Subject})
		}
	}

	if filter.Pin != nil {
		mongoFilter = append(mongoFilter, primitive.E{Key: "poll.pin", Value: *filter.Pin})
	}

	if len(filter.GroupIDs) > 0 {
		mongoFilter = append(mongoFilter, primitive.E{Key: "poll.group_id", Value: bson.M{"$in": filter.GroupIDs}})
	}

	if len(filter.Statuses) > 0 {
		mongoFilter = append(mongoFilter, primitive.E{Key: "poll.status", Value: bson.M{"$in": filter.Statuses}})
	}

	if filterByToMembers {
		var innerFilter primitive.M
		if membership != nil && len(membership.GroupIDsAsAdmin) > 0 {
			innerFilter = primitive.M{"$or": []primitive.M{
				primitive.M{"poll.group_id": bson.M{"$in": membership.GroupIDsAsAdmin}},
				primitive.M{"poll.to_members.user_id": user.Claims.Subject},
			}}
		} else {
			innerFilter = primitive.M{"poll.to_members.user_id": user.Claims.Subject}
		}

		mongoFilter = append(mongoFilter, primitive.E{Key: "$or", Value: []primitive.M{
			primitive.M{"poll.to_members": primitive.Null{}},
			primitive.M{"poll.to_members": primitive.M{"$exists": true, "$size": 0}},
			primitive.M{"poll.user_id": user.Claims.Subject},
			innerFilter,
		}})
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"poll.status", 1}, {"_id", -1}})

	if filter.Limit != nil {
		findOptions.SetLimit(*filter.Limit)
	}
	if filter.Offset != nil {
		findOptions.SetSkip(*filter.Offset)
	}

	var list []model.Poll
	err := sa.db.polls.Find(mongoFilter, &list, findOptions)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetPoll retrieves a single poll
func (sa *Adapter) GetPoll(user *model.User, id string, filterByToMembers bool, membership *groups.GroupMembership) (*model.Poll, error) {

	if objID, err := primitive.ObjectIDFromHex(id); err == nil {
		filter := bson.D{
			primitive.E{Key: "org_id", Value: user.Claims.OrgID},
			primitive.E{Key: "_id", Value: objID},
		}

		if filterByToMembers {
			var innerFilter primitive.M
			if membership != nil && len(membership.GroupIDsAsAdmin) > 0 {
				innerFilter = primitive.M{"$or": []primitive.M{
					primitive.M{"poll.group_id": bson.M{"$in": membership.GroupIDsAsAdmin}},
					primitive.M{"poll.to_members.user_id": user.Claims.Subject},
				}}
			} else {
				innerFilter = primitive.M{"poll.to_members.user_id": user.Claims.Subject}
			}

			filter = append(filter, primitive.E{Key: "$or", Value: []primitive.M{
				primitive.M{"poll.to_members": primitive.Null{}},
				primitive.M{"poll.to_members": primitive.M{"$exists": true, "$size": 0}},
				primitive.M{"poll.user_id": user.Claims.Subject},
				innerFilter,
			}})
		}

		var list []model.Poll
		err := sa.db.polls.Find(filter, &list, &options.FindOptions{})
		if err != nil {
			fmt.Printf("error storage.Adapter.GetPoll(%s) - %s", id, err)
			return nil, fmt.Errorf("error storage.Adapter.GetPoll(%s) - %s", id, err)
		}

		if len(list) > 0 {
			entry := list[0]
			return &entry, nil
		}
	} else {
		fmt.Printf("error storage.Adapter.GetPoll(%s) - unable to construct obj id", id)
		return nil, fmt.Errorf("error storage.Adapter.GetPoll(%s) - unable to construct obj id", id)
	}
	return nil, nil
}

// CreatePoll creates a poll
func (sa *Adapter) CreatePoll(user *model.User, poll model.Poll) (*model.Poll, error) {
	poll.OrgID = user.Claims.OrgID
	poll.ID = primitive.NewObjectID()
	poll.UserID = user.Claims.Subject
	poll.UserName = user.Claims.Name
	poll.DateCreated = time.Now()

	_, err := sa.db.polls.InsertOne(poll)
	if err != nil {
		fmt.Printf("error storage.Adapter.CreatePoll(%s) - %s", poll.ID, err)
		return nil, fmt.Errorf("error storage.Adapter.CreatePoll(%s) - %s", poll.ID, err)
	}

	return &poll, nil
}

// UpdatePoll updates a poll
func (sa *Adapter) UpdatePoll(user *model.User, poll model.Poll) (*model.Poll, error) {

	if len(poll.ID) > 0 {

		poll.DateUpdated = time.Now().UTC()
		filter := bson.D{
			primitive.E{Key: "org_id", Value: user.Claims.OrgID},
			primitive.E{Key: "_id", Value: poll.ID},
		}

		update := bson.D{
			primitive.E{Key: "$set", Value: bson.D{
				primitive.E{Key: "poll.date_updated", Value: poll.DateUpdated},
				primitive.E{Key: "poll.to_members", Value: poll.ToMembersList},
				primitive.E{Key: "poll.pin", Value: poll.Pin},
				primitive.E{Key: "poll.question", Value: poll.Question},
				primitive.E{Key: "poll.options", Value: poll.Options},
				primitive.E{Key: "poll.group_id", Value: poll.GroupID},
				primitive.E{Key: "poll.multi_choice", Value: poll.MultiChoice},
				primitive.E{Key: "poll.repeat", Value: poll.Repeat},
				primitive.E{Key: "poll.show_results", Value: poll.ShowResults},
				primitive.E{Key: "poll.stadium", Value: poll.Stadium},
				primitive.E{Key: "poll.geo_fence", Value: poll.Geo},
				primitive.E{Key: "poll.status", Value: poll.Status},
			}},
		}

		_, err := sa.db.polls.UpdateOne(filter, update, nil)
		if err != nil {
			fmt.Printf("error storage.Adapter.UpdatePoll(%s) - %s", poll.ID, err)
			return nil, fmt.Errorf("error storage.Adapter.UpdatePoll(%s) - %s", poll.ID, err)
		}
	}

	return &poll, nil
}

// StartPoll starts an existing poll
func (sa *Adapter) StartPoll(user *model.User, pollID string) error {

	poll, err := sa.GetPoll(user, pollID, true, nil)
	if err != nil {
		fmt.Printf("error storage.Adapter.StartPoll(%s) - %s", pollID, err)
		return fmt.Errorf("error storage.Adapter.StartPoll(%s) - %s", pollID, err)
	}

	if poll != nil {
		if poll.Status != PollStatusStarted {
			poll.Status = PollStatusStarted
			_, err = sa.UpdatePoll(user, *poll)
			if err != nil {
				fmt.Printf("error storage.Adapter.StartPoll(%s) - %s", pollID, err)
				return fmt.Errorf("error storage.Adapter.StartPoll(%s) - %s", pollID, err)
			}
		}
	} else {
		return fmt.Errorf("error storage.Adapter.EndPoll(%s) - poll not found: %s", pollID, err)
	}

	return nil
}

// EndPoll ends an existing poll
func (sa *Adapter) EndPoll(user *model.User, pollID string) error {
	poll, err := sa.GetPoll(user, pollID, true, nil)
	if err != nil {
		fmt.Printf("error storage.Adapter.EndPoll(%s) - %s", pollID, err)
		return fmt.Errorf("error storage.Adapter.EndPoll(%s) - %s", pollID, err)
	}

	if poll != nil {
		if poll.Status != PollStatusTerminated {
			poll.Status = PollStatusTerminated
			_, err = sa.UpdatePoll(user, *poll)
			if err != nil {
				fmt.Printf("error storage.Adapter.EndPoll(%s) - %s", pollID, err)
				return fmt.Errorf("error storage.Adapter.EndPoll(%s) - %s", pollID, err)
			}
		}
	} else {
		return fmt.Errorf("error storage.Adapter.EndPoll(%s) - poll not found: %s", pollID, err)
	}

	return nil
}

// DeletePoll deletes a poll
func (sa *Adapter) DeletePoll(user *model.User, id string) error {
	if objID, err := primitive.ObjectIDFromHex(id); err == nil {
		filter := bson.D{
			primitive.E{Key: "org_id", Value: user.Claims.OrgID},
			primitive.E{Key: "_id", Value: objID},
		}
		_, err := sa.db.polls.DeleteOne(filter, nil)
		if err != nil {
			fmt.Printf("error storage.Adapter.DeletePoll(): error while delete poll (%s) - %s", id, err)
			return fmt.Errorf("error storage.Adapter.DeletePoll(): error while delete poll (%s) - %s", id, err)
		}

	}
	return nil

}

// VotePoll votes a poll
func (sa *Adapter) VotePoll(user *model.User, pollID string, vote model.PollVote) error {

	if objID, err := primitive.ObjectIDFromHex(pollID); err == nil {
		now := time.Now().UTC()
		vote.Created = now

		filter := bson.D{
			primitive.E{Key: "_id", Value: objID},
		}

		update := bson.D{
			primitive.E{Key: "$set", Value: bson.D{
				primitive.E{Key: "poll.date_updated", Value: now},
			}},
			primitive.E{Key: "$push", Value: bson.D{
				primitive.E{Key: "responses", Value: vote},
			}},
		}

		_, err := sa.db.polls.UpdateOne(filter, update, nil)
		if err != nil {
			fmt.Printf("error storage.Adapter.VotePoll(%s) - %s", pollID, err)
			return fmt.Errorf("error storage.Adapter.VotePoll(%s) - %s", pollID, err)
		}
	}
	return nil
}

// SetListener sets the upper layer listener for sending collection changed callbacks
func (sa *Adapter) SetListener(listener CollectionListener) {
	sa.db.listener = listener
}

// Event

func (m *database) onDataChanged(changeDoc map[string]interface{}) {
	if changeDoc == nil {
		return
	}
	log.Printf("onDataChanged: %+v\n", changeDoc)
	ns := changeDoc["ns"]
	if ns == nil {
		return
	}
	nsMap := ns.(map[string]interface{})
	coll := nsMap["coll"]

	record := changeDoc["fullDocument"]
	var recordMap map[string]interface{}
	if record != nil {
		recordMap = record.(map[string]interface{})
	}

	if m.listener != nil {
		m.listener.OnCollectionUpdated(coll.(string), recordMap)
	}
}

// Get Survey
func (sa *Adapter) GetSurvey(id string) (*model.Survey, error) {

	filter := bson.M{"_id": id}
	var list []model.Survey
	err := sa.db.surveys.Find(filter, &list, &options.FindOptions{})
	if err != nil {
		fmt.Printf("error storage.Adapter.GetSurvey(%s) - %s", id, err)
		return nil, fmt.Errorf("error storage.Adapter.GetSurvey(%s) - %s", id, err)
	}

	if len(list) > 0 {
		entry := list[0]
		return &entry, nil
	} else {
		return nil, nil
	}
}

// Create Survey
func (sa *Adapter) CreateSurvey(survey model.Survey) (*model.Survey, error) {

	_, err := sa.db.surveys.InsertOne(survey)
	if err != nil {
		fmt.Printf("error storage.Adapter.CreateSurvey(%s) - %s", survey.ID, err)
		return nil, fmt.Errorf("error storage.Adapter.CreateSurvey(%s) - %s", survey.ID, err)
	}

	return &survey, nil
}

// Update Survey
func (sa *Adapter) UpdateSurvey(user *model.User, survey model.Survey) error {

	if len(survey.ID) > 0 {

		now := time.Now().UTC()
		filter := bson.M{"_id": survey.ID, "creator_id": user.Claims.Subject}
		update := bson.M{"$set": bson.M{
			"questions":    survey.Questions,
			"scored":       survey.Scored,
			"result_rule":  survey.ResultRule,
			"type":         survey.Type,
			"stats":        survey.SurveyStats,
			"sensitive":    survey.Sensitive,
			"date_updated": now,
		}}

		result, err := sa.db.surveys.UpdateOne(filter, update, nil)
		if err != nil {
			fmt.Printf("error storage.Adapter.UpdateSurvey(%s) - %s", survey.ID, err)
			return fmt.Errorf("error storage.Adapter.UpdateSurvey(%s) - %s", survey.ID, err)
		}

		modifiedCount := result.ModifiedCount
		if modifiedCount == 0 {
			return fmt.Errorf("error storage.Adapter.DeleteSurvey(): 403 - Not the creator ", err)
		}
	}

	return nil
}

// Delete Survey
func (sa *Adapter) DeleteSurvey(user *model.User, id string) error {

	filter := bson.M{"_id": id, "creator_id": user.Claims.Subject}

	result, err := sa.db.surveys.DeleteOne(filter, nil)
	if err != nil {
		return fmt.Errorf("error storage.Adapter.DeleteSurvey(): error while delete survey (%s) - %s", id, err)
	}

	deletedCount := result.DeletedCount
	if deletedCount == 0 {
		return fmt.Errorf("error storage.Adapter.DeleteSurvey(): 403 - Not the creator (%s) - %s", id, err)
	}

	return nil
}

// Get SurveyResponse
func (sa *Adapter) GetSurveyResponse(user *model.User, id string) (*model.SurveyResponse, error) {

	filter := bson.M{"_id": id}
	var entry model.SurveyResponse
	err := sa.db.surveyResponses.FindOne(filter, &entry, nil)
	if err != nil {
		fmt.Printf("error storage.Adapter.GetSurveyResponse(%s) - %s", id, err)
		return nil, fmt.Errorf("error storage.Adapter.GetSurveyResponse(%s) - %s", id, err)
	}

	if entry.AppID != user.Claims.AppID || entry.OrgID != user.Claims.OrgID {
		fmt.Printf("error storage.Adapter.GetSurveyResponse(%s) - 403 incorrect AppID or OrgID", id)
		return nil, fmt.Errorf("error storage.Adapter.GetSurveyResponse(%s) - 403 incorrect AppID or OrgID", id)
	} else {
		return &entry, nil
	}
}

// Create SurveyResponse
func (sa *Adapter) CreateSurveyResponse(surveyResponse model.SurveyResponse) (*model.SurveyResponse, error) {

	_, err := sa.db.surveyResponses.InsertOne(surveyResponse)
	if err != nil {
		fmt.Printf("error storage.Adapter.CreateSurveyResponse(%s) - %s", surveyResponse.ID, err)
		return nil, fmt.Errorf("error storage.Adapter.CreateSurveyResponse(%s) - %s", surveyResponse.ID, err)
	}

	return &surveyResponse, nil
}

// Update SurveyResponse
func (sa *Adapter) UpdateSurveyResponse(user *model.User, surveyResponse model.SurveyResponse) error {

	if len(surveyResponse.ID) > 0 {

		now := time.Now().UTC()
		filter := bson.M{"_id": surveyResponse.ID, "user_id": user.Claims.Subject}
		update := bson.M{"$set": bson.M{
			"survey":       surveyResponse.Survey,
			"date_updated": now,
		}}

		result, err := sa.db.surveyResponses.UpdateOne(filter, update, nil)
		if err != nil {
			fmt.Printf("error storage.Adapter.UpdateSurveyResponse(%s) - %s", surveyResponse.ID, err)
			return fmt.Errorf("error storage.Adapter.UpdateSurveyResponse(%s) - %s", surveyResponse.ID, err)
		}

		modifiedCount := result.ModifiedCount
		if modifiedCount == 0 {
			return fmt.Errorf("error storage.Adapter.DeleteSurveyResponse(): 403 - Not the user ")
		}
	}

	return nil
}

// Delete SurveyResponse
func (sa *Adapter) DeleteSurveyResponse(user *model.User, id string) error {

	filter := bson.M{"_id": id, "user_id": user.Claims.Subject}

	result, err := sa.db.surveyResponses.DeleteOne(filter, nil)
	if err != nil {
		return fmt.Errorf("error storage.Adapter.DeleteSurveyResponse(): error while delete survey response (%s) - %s", id, err)
	}

	deletedCount := result.DeletedCount
	if deletedCount == 0 {
		return fmt.Errorf("error storage.Adapter.DeleteSurveyResponse(): 403 - Not the creator (%s) - %s", id, err)
	}

	return nil
}
