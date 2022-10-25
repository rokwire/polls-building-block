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

// Survey wraps the entire record
type Survey struct {
	ID          string                `json:"id" bson:"_id"`
	CreatorID   string                `json:"creator_id" bson:"creator_id"`
	OrgID       string                `json:"org_id" bson:"org_id"`
	AppID       string                `json:"app_id" bson:"app_id"`
	Questions   map[string]SurveyData `json:"questions" bson:"questions"`
	Scored      bool                  `json:"scored" bson:"scored"`
	ResultRule  string                `json:"result_rule" bson:"result_rule"`
	Type        string                `json:"type" bson:"type"`
	SurveyStats *SurveyStats          `json:"stats" bson:"stats"`
	Sensitive   bool                  `json:"sensitive" bson:"sensitive"`
	DateCreated time.Time             `json:"date_created" bson:"date_created"`
	DateUpdated *time.Time            `json:"date_updated" bson:"date_updated"`
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
	Response            interface{} `json:"response" bson:"response"`

	Type string `json:"type" bson:"type"`

	// Shared
	OkAnswer  interface{}              `json:"ok_answer,omitempty" bson:"ok_answer,omitempty"`
	OkAnswers []interface{}            `json:"ok_answers,omitempty" bson:"ok_answers,omitempty"`
	Options   []map[string]interface{} `json:"options,omitempty" bson:"options,omitempty"`
	Action    *ActionData              `json:"action,omitempty" bson:"action,omitempty"`

	// True/False
	YesNo *bool `json:"yes_no,omitempty" bson:"yes_no,omitempty"`

	// Multiple Choice
	AllowMultiple *bool `json:"allow_multiple,omitempty" bson:"allow_multiple,omitempty"`
	CheckAll      *bool `json:"check_all,omitempty" bson:"check_all,omitempty"`

	// DateTime
	StartTime *time.Time `json:"start_time,omitempty" bson:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty" bson:"end_time,omitempty"`
	AskTime   *bool      `json:"ask_time,omitempty" bson:"ask_time,omitempty"`

	// Numeric
	Minimum   *float64 `json:"minimum,omitempty" bson:"minimum,omitempty"`
	Maximum   *float64 `json:"maximum,omitempty" bson:"maximum,omitempty"`
	WholeNum  *bool    `json:"whole_num,omitempty" bson:"whole_num,omitempty"`
	Slider    *bool    `json:"slider,omitempty" bson:"slider,omitempty"`
	SelfScore *bool    `json:"self_score,omitempty" bson:"self_score,omitempty"`

	// Text
	MinLength *int `json:"min_length,omitempty" bson:"min_length,omitempty"`
	MaxLength *int `json:"max_length,omitempty" bson:"max_length,omitempty"`

	// DataEntry
	DataFormat map[string]string `json:"data_format,omitempty" bson:"data_format,omitempty"`

	// Response
	Body *string `json:"body,omitempty" bson:"body,omitempty"`

	// Survey
	Survey *Survey `json:"survey,omitempty" bson:"survey,omitempty"`
}

// SurveyResponseFollowUp is the response follow up to Survey
type SurveyResponseFollowUp struct {
	Key   interface{} `json:"key" bson:"key"`
	Value SurveyData  `json:"value" bson:"value"`
}

// ActionData is the wrapped within SurveyData
type ActionData struct {
	Type  string  `json:"type" bson:"type"`
	Label *string `json:"label" bson:"label"`
	Data  string  `json:"data" bson:"data"`
}
