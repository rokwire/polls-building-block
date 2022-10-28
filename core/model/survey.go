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

package model

import (
	"time"
)

// SurveyResponse wraps the entire survey response
type SurveyResponse struct {
	ID          string     `json:"id" bson:"_id"`
	UserID      string     `json:"user_id" bson:"user_id"`
	OrgID       string     `json:"org_id" bson:"org_id"`
	AppID       string     `json:"app_id" bson:"app_id"`
	Survey      Survey     `json:"survey" bson:"survey"`
	DateCreated time.Time  `json:"date_created" bson:"date_created"`
	DateUpdated *time.Time `json:"date_updated" bson:"date_updated"`
}

// Survey wraps the entire record
type Survey struct {
	ID                 string                 `json:"id" bson:"_id"`
	CreatorID          string                 `json:"creator_id" bson:"creator_id"`
	OrgID              string                 `json:"org_id" bson:"org_id"`
	AppID              string                 `json:"app_id" bson:"app_id"`
	Title              string                 `json:"title" bson:"title"`
	MoreInfo           *string                `json:"more_info" bson:"more_info"`
	Data               map[string]SurveyData  `json:"data" bson:"data"`
	Scored             bool                   `json:"scored" bson:"scored"`
	ResultRules        string                 `json:"result_rules" bson:"result_rules"`
	Type               string                 `json:"type" bson:"type"`
	SurveyStats        *SurveyStats           `json:"stats" bson:"stats"`
	Sensitive          bool                   `json:"sensitive" bson:"sensitive"`
	DefaultDataKey     *string                `json:"default_data_key" bson:"default_data_key"`
	DefaultDataKeyRule *string                `json:"default_data_key_rule" bson:"default_data_key_rule"`
	Constants          map[string]interface{} `json:"constants" bson:"constants"`
	Strings            map[string]interface{} `json:"strings" bson:"strings"`
	SubRules           map[string]interface{} `json:"sub_rules" bson:"sub_rules"`
	DateCreated        time.Time              `json:"date_created" bson:"date_created"`
	DateUpdated        *time.Time             `json:"date_updated" bson:"date_updated"`
}

// SurveyStats are stats of a Survey
type SurveyStats struct {
	Total    int                `json:"total" bson:"total"`
	Complete int                `json:"complete" bson:"complete"`
	Scored   int                `json:"scored" bson:"scored"`
	Scores   map[string]float32 `json:"scores" bson:"scores"`
}

// SurveyData is data stored for a Survey
type SurveyData struct {
	Section             *string     `json:"section" bson:"section"`
	AllowSkip           bool        `json:"allow_skip" bson:"allow_skip"`
	Text                string      `json:"text" bson:"text"`
	MoreInfo            string      `json:"more_info" bson:"more_info"`
	DefaultFollowUpKey  *string     `json:"default_follow_up_key" bson:"default_follow_up_key"`
	DefaultResponseRule *string     `json:"default_response_rule" bson:"default_response_rule"`
	FollowUpRule        *string     `json:"follow_up_rule" bson:"follow_up_rule"`
	ScoreRule           *string     `json:"score_rule" bson:"score_rule"`
	Replace             bool        `json:"replace" bson:"replace"`
	Response            interface{} `json:"response" bson:"response"`

	Type string `json:"type" bson:"type"`

	// Shared
	CorrectAnswer  interface{}              `json:"correct_answer,omitempty" bson:"correct_answer,omitempty"`
	CorrectAnswers []interface{}            `json:"correct_answers,omitempty" bson:"correct_answers,omitempty"`
	Options        []map[string]interface{} `json:"options,omitempty" bson:"options,omitempty"`
	Actions        []ActionData             `json:"actions,omitempty" bson:"actions,omitempty"`
	SelfScore      *bool                    `json:"self_score,omitempty" bson:"self_score,omitempty"`

	// True/False
	YesNo *bool `json:"yes_no,omitempty" bson:"yes_no,omitempty"`

	// Multiple Choice
	AllowMultiple *bool `json:"allow_multiple,omitempty" bson:"allow_multiple,omitempty"`

	// DateTime
	StartTime *time.Time `json:"start_time,omitempty" bson:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty" bson:"end_time,omitempty"`
	AskTime   *bool      `json:"ask_time,omitempty" bson:"ask_time,omitempty"`

	// Numeric
	Minimum  *float64 `json:"minimum,omitempty" bson:"minimum,omitempty"`
	Maximum  *float64 `json:"maximum,omitempty" bson:"maximum,omitempty"`
	WholeNum *bool    `json:"whole_num,omitempty" bson:"whole_num,omitempty"`
	Slider   *bool    `json:"slider,omitempty" bson:"slider,omitempty"`

	// Text
	MinLength *int `json:"min_length,omitempty" bson:"min_length,omitempty"`
	MaxLength *int `json:"max_length,omitempty" bson:"max_length,omitempty"`

	// DataEntry
	DataFormat map[string]string `json:"data_format,omitempty" bson:"data_format,omitempty"`
}

// ActionData is the wrapped within SurveyData
type ActionData struct {
	Type  string  `json:"type" bson:"type"`
	Label *string `json:"label" bson:"label"`
	Data  string  `json:"data" bson:"data"`
}
