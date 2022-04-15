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

package storage

import (
	"fmt"
	"github.com/rokwire/core-auth-library-go/tokenauth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"polls/core/model"
	"strconv"
	"time"
)

const (
	statusCreated    = "created"
	statusStarted    = "started"
	statusTerminated = "terminated"
	settingsKey      = "stadium"
	eventInterval    = 100 * time.Millisecond
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
func (sa *Adapter) GetPolls(user *tokenauth.Claims, IDs []string, userID *string, offset *int64, limit *int64, order *string, filterByToMembers bool) ([]model.Poll, error) {
	filter := bson.D{
		primitive.E{Key: "org_id", Value: user.OrgID},
	}
	innerFilter := []interface{}{}
	if userID != nil {
		innerFilter = append(innerFilter, bson.D{primitive.E{Key: "poll.user_id", Value: userID}})
	}
	if len(IDs) > 0 {
		reconstructedIDs := []primitive.ObjectID{}
		for _, id := range IDs {
			if objID, err := primitive.ObjectIDFromHex(id); err == nil {
				reconstructedIDs = append(reconstructedIDs, objID)
			}
		}
		filter = append(filter, primitive.E{Key: "_id", Value: bson.M{"$in": reconstructedIDs}})
	}
	if filterByToMembers {
		filter = append(filter, primitive.E{Key: "$or", Value: []primitive.M{
			primitive.M{"poll.to_members": primitive.Null{}},
			primitive.M{"poll.to_members": primitive.M{"$exists": true, "$size": 0}},
			primitive.M{"poll.to_members.user_id": user.Subject},
			primitive.M{"poll.user_id": user.Subject},
		}})
	}

	findOptions := options.Find()
	if order != nil && *order == "asc" {
		findOptions.SetSort(bson.D{{"_id", 1}})
	} else {
		findOptions.SetSort(bson.D{{"_id", -1}})
	}
	if limit != nil {
		findOptions.SetLimit(*limit)
	}
	if offset != nil {
		findOptions.SetSkip(*offset)
	}

	var list []model.Poll
	err := sa.db.polls.Find(filter, &list, findOptions)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetPoll retrieves a single poll
func (sa *Adapter) GetPoll(user *tokenauth.Claims, id string) (*model.Poll, error) {

	if objID, err := primitive.ObjectIDFromHex(id); err == nil {
		filter := bson.D{
			primitive.E{Key: "org_id", Value: user.OrgID},
			primitive.E{Key: "_id", Value: objID},
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
func (sa *Adapter) CreatePoll(user *tokenauth.Claims, poll model.Poll) (*model.Poll, error) {
	poll.DateCreated = time.Now()
	poll.OrgID = user.OrgID

	_, err := sa.db.polls.InsertOne(poll)
	if err != nil {
		fmt.Printf("error storage.Adapter.CreatePoll(%s) - %s", poll.ID, err)
		return nil, fmt.Errorf("error storage.Adapter.CreatePoll(%s) - %s", poll.ID, err)
	}

	return &poll, nil
}

// UpdatePoll updates a poll
func (sa *Adapter) UpdatePoll(user *tokenauth.Claims, poll model.Poll) (*model.Poll, error) {
	if len(poll.ID) > 0 {

		poll.DateUpdated = time.Now().UTC()
		filter := bson.D{
			primitive.E{Key: "org_id", Value: user.OrgID},
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
func (sa *Adapter) StartPoll(user *tokenauth.Claims, pollID string) error {

	poll, err := sa.GetPoll(user, pollID)
	if err != nil {
		fmt.Printf("error storage.Adapter.StartPoll(%s) - %s", pollID, err)
		return fmt.Errorf("error storage.Adapter.StartPoll(%s) - %s", pollID, err)
	}

	if poll != nil {
		if poll.Status != statusStarted {
			poll.Status = statusStarted
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
func (sa *Adapter) EndPoll(user *tokenauth.Claims, pollID string) error {
	poll, err := sa.GetPoll(user, pollID)
	if err != nil {
		fmt.Printf("error storage.Adapter.EndPoll(%s) - %s", pollID, err)
		return fmt.Errorf("error storage.Adapter.EndPoll(%s) - %s", pollID, err)
	}

	if poll != nil {
		if poll.Status != statusTerminated {
			poll.Status = statusTerminated
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
func (sa *Adapter) DeletePoll(user *tokenauth.Claims, id string) error {
	filter := bson.D{
		primitive.E{Key: "org_id", Value: user.OrgID},
		primitive.E{Key: "_id", Value: id},
	}
	_, err := sa.db.polls.DeleteOne(filter, nil)
	if err != nil {
		fmt.Printf("error storage.Adapter.DeletePoll(): error while delete poll (%s) - %s", id, err)
		return fmt.Errorf("error storage.Adapter.DeletePoll(): error while delete poll (%s) - %s", id, err)
	}

	return nil

}

// VotePoll votes a poll
func (sa *Adapter) VotePoll(user *tokenauth.Claims, pollID string, vote model.PollVote) error {

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
	recordMap := record.(map[string]interface{})

	if m.listener != nil {
		m.listener.OnCollectionUpdated(coll.(string), recordMap)
	}
}
