package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// PollData data stored for a poll
type PollData struct {
	UserID      string    `json:"userid" bson:"userid" validate:"required"`
	UserName    string    `json:"username" bson:"username" validate:"required"`
	Question    string    `json:"question" bson:"question" validate:"required"`
	Options     []string  `json:"options" bson:"options" validate:"required,min=2,dive,required"`
	GroupID     string    `json:"group_id,omitempty" bson:"group_id"`
	Pin         int       `json:"pin,omitempty" bson:"pin" validate:"min=0,max=9999"`
	MultiChoice bool      `json:"multi_choice" bson:"multi_choice"`
	Repeat      bool      `json:"repeat" bson:"repeat"`
	ShowResults bool      `json:"show_results" bson:"show_results"`
	Stadium     string    `json:"stadium" bson:"stadium"`
	Geo         bool      `json:"geo_fence" bson:"geo_fence"`
	Status      string    `json:"status" bson:"status" validate:"required,oneof=created started"`
	DateCreated time.Time `json:"date_created" bson:"date_created"`
	DateUpdated time.Time `json:"date_updated" bson:"date_updated"`
} // @name PollData

// Poll wraps the entire record
type Poll struct {
	PollData  `json:"poll" bson:"poll"`
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Responses []PollVote         `json:"responses" bson:"responses,omitempty" validate:"max=0"`
	Results   []int              `json:"results" bson:"results,omitempty" validate:"max=0"`
} // @name Poll

// ToPollResult converts to PollResult
func (poll *Poll) ToPollResult() PollResult {

	result := PollResult{
		PollData: poll.PollData,
	}

	votersMap := make(map[string]bool)

	result.ID = poll.ID
	result.PollData = poll.PollData

	count := len(result.PollData.Options)
	result.Results = make([]int, count)

	if len(poll.Responses) > 0 {
		for _, e := range poll.Responses {
			votersMap[e.UserID] = true

			for _, a := range e.Answer {
				if a >= 0 && a < count {
					result.Results[a]++
				}
			}
		}
	} else {
		copy(result.Results, poll.Results)
	}

	result.UniqueVotersCount = len(votersMap)

	for _, n := range result.Results {
		result.Total += n
	}

	return result
}

// PollVote data stored for each response
type PollVote struct {
	UserID  string    `json:"userid" validate:"required"`
	Answer  []int     `json:"answer" validate:"required,min=1"`
	Created time.Time `json:"created"`
} // @name PollVote

// PollResult wraps poll result
type PollResult struct {
	PollData          `json:"poll" bson:""`
	ID                primitive.ObjectID `json:"id"`
	Voted             []int              `json:"voted,omitempty"`
	Results           []int              `json:"results"`
	UniqueVotersCount int                `json:"unique_voters_count"`
	Total             int                `json:"total"`
} // @name PollResult
