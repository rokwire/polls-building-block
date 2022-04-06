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
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CollectionListener listens for collection updates
type CollectionListener interface {
	OnCollectionUpdated(name string)
}

type database struct {
	listener CollectionListener

	mongoDBAuth  string
	mongoDBName  string
	mongoTimeout time.Duration

	db       *mongo.Database
	dbClient *mongo.Client

	polls *collectionWrapper
}

func (m *database) start() error {

	log.Println("database -> start")

	//connect to the database
	clientOptions := options.Client().ApplyURI(m.mongoDBAuth)
	connectContext, cancel := context.WithTimeout(context.Background(), m.mongoTimeout)
	client, err := mongo.Connect(connectContext, clientOptions)
	cancel()
	if err != nil {
		return err
	}

	//ping the database
	pingContext, cancel := context.WithTimeout(context.Background(), m.mongoTimeout)
	err = client.Ping(pingContext, nil)
	cancel()
	if err != nil {
		return err
	}

	//apply checks
	db := client.Database(m.mongoDBName)

	//asign the db, db client and the collections
	m.db = db
	m.dbClient = client

	polls := &collectionWrapper{database: m, coll: db.Collection("polls")}
	err = m.applyRewardTypesChecks(polls)
	if err != nil {
		return err
	}
	go polls.Watch(nil)

	m.polls = polls

	return nil
}

func (m *database) applyRewardTypesChecks(posts *collectionWrapper) error {
	log.Println("apply polls checks.....")

	indexes, _ := posts.ListIndexes()
	indexMapping := map[string]interface{}{}
	if indexes != nil {

		for _, index := range indexes {
			name := index["name"].(string)
			indexMapping[name] = index
		}
	}

	if indexMapping["org_id_1"] == nil {
		err := posts.AddIndex(
			bson.D{
				primitive.E{Key: "org_id", Value: 1},
			}, false)
		if err != nil {
			return err
		}
	}

	if indexMapping["poll.group_id_1"] == nil {
		err := posts.AddIndex(
			bson.D{
				primitive.E{Key: "poll.group_id", Value: 1},
			}, false)
		if err != nil {
			return err
		}
	}

	if indexMapping["poll.pin_1_poll.status_1__id_1"] == nil {
		err := posts.AddIndex(
			bson.D{
				primitive.E{Key: "poll.pin", Value: 1},
				primitive.E{Key: "poll.status", Value: 1},
				primitive.E{Key: "_id", Value: 1},
			}, false)
		if err != nil {
			return err
		}
	}

	if indexMapping["poll.userid_1_poll.status_1__id_1"] == nil {
		err := posts.AddIndex(
			bson.D{
				primitive.E{Key: "poll.userid", Value: 1},
				primitive.E{Key: "poll.status", Value: 1},
				primitive.E{Key: "_id", Value: 1},
			}, false)
		if err != nil {
			return err
		}
	}

	if indexMapping["responses.userid_1_poll.status_1__id_1"] == nil {
		err := posts.AddIndex(
			bson.D{
				primitive.E{Key: "responses.userid", Value: 1},
				primitive.E{Key: "poll.status", Value: 1},
				primitive.E{Key: "_id", Value: 1},
			}, false)
		if err != nil {
			return err
		}
	}

	log.Println("polls checks passed")
	return nil
}
