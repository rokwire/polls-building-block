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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RecordCount wraps count aggregation
type RecordCount struct {
	Count int `json:"count" bson:"count"`
}

func (sa *Adapter) applyMultiTenancy() error {
	log.Printf("applyMultiTenancy started ")
	settingsCount, err := sa.getNonOrgRecordCount(sa.db.settings)
	if err != nil {
		log.Printf("error storage.Adapter.applyMultiTenancy() - %s", err)
		return fmt.Errorf("error storage.Adapter.applyMultiTenancy() - %s", err)
	}

	if settingsCount != nil && settingsCount.Count > 0 {
		err = sa.migrateNonOrgRecords(sa.db.settings)
		if err != nil {
			log.Printf("error storage.Adapter.applyMultiTenancy() - %s", err)
			return fmt.Errorf("error storage.Adapter.applyMultiTenancy() - %s", err)
		}
		log.Printf("migrate %d settings records successfully", settingsCount.Count)
	}

	pollsCount, err := sa.getNonOrgRecordCount(sa.db.polls)
	if err != nil {
		log.Printf("error storage.Adapter.applyMultiTenancy() - %s", err)
		return fmt.Errorf("error storage.Adapter.applyMultiTenancy() - %s", err)
	}

	if pollsCount != nil && pollsCount.Count > 0 {
		err = sa.migrateNonOrgRecords(sa.db.polls)
		if err != nil {
			log.Printf("error storage.Adapter.applyMultiTenancy() - %s", err)
			return fmt.Errorf("error storage.Adapter.applyMultiTenancy() - %s", err)
		}
		log.Printf("migrate %d settings records successfully", pollsCount.Count)
	}

	log.Printf("applyMultiTenancy ended")

	return nil
}

// getNonOrgRecordCount Gets the count of non org polls
func (sa *Adapter) getNonOrgRecordCount(wrapper *collectionWrapper) (*RecordCount, error) {

	pipeline := []bson.M{
		{"$match": bson.M{"org_id": bson.M{"$exists": false}}},
		{"$group": bson.M{"_id": "non_org_polls", "count": bson.M{"$sum": 1}}},
	}

	var result []RecordCount
	err := wrapper.Aggregate(pipeline, &result, nil)
	if err != nil {
		log.Printf("storage.GetNonOrgPollsCount error: %s", err)
		return nil, fmt.Errorf("storage.GetNonOrgPollsCount error: %s", err)
	}

	if len(result) > 0 {
		first := result[0]
		return &first, nil
	}
	return &RecordCount{Count: 0}, nil
}

// migrateNonOrgRecords Migrate all non org records
func (sa *Adapter) migrateNonOrgRecords(wrapper *collectionWrapper) error {

	filter := bson.D{
		primitive.E{Key: "org_id", Value: bson.M{"$exists": false}},
	}

	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "org_id", Value: sa.config.UiucOrgID},
		}},
	}

	_, err := wrapper.UpdateMany(filter, update, nil)
	if err != nil {
		fmt.Printf("error storage.Adapter.MigrateNonOrgPolls() - %s", err)
		return fmt.Errorf("error storage.Adapter.MigrateNonOrgPolls() - %s", err)
	}
	return nil
}
