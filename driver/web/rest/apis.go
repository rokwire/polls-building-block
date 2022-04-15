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

package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rokwire/core-auth-library-go/tokenauth"
	"io/ioutil"
	"log"
	"net/http"
	"polls/core"
	"polls/core/model"
)

const maxUploadSize = 15 * 1024 * 1024 // 15 mb

//ApisHandler handles the rest APIs implementation
type ApisHandler struct {
	app    *core.Application
	config *model.Config
}

//Version gives the service version
// @Description Gives the service version.
// @Tags Client
// @ID Version
// @Produce plain
// @Success 200
// @Router /version [get]
func (h ApisHandler) Version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.app.Services.GetVersion()))
}

// NewApisHandler creates new rest Handler instance
func NewApisHandler(app *core.Application, config *model.Config) ApisHandler {
	return ApisHandler{app: app, config: config}
}

// NewAdminApisHandler creates new rest Handler instance
func NewAdminApisHandler(app *core.Application, config *model.Config) AdminApisHandler {
	return AdminApisHandler{app: app, config: config}
}

// NewInternalApisHandler creates new rest Handler instance
func NewInternalApisHandler(app *core.Application, config *model.Config) InternalApisHandler {
	return InternalApisHandler{app: app, config: config}
}

type pollIDsRequestBody struct {
	IDs []string `json:"ids"`
} // @name pollIDsRequestBody

// GetPolls Retrieves  all polls by a filter params
// @Description Retrieves  all polls by a filter params
// @Tags Client
// @ID GetPolls
// @Param offset query string false "offset"
// @Param limit query string false "limit - limit the result"
// @Param order query string false "order - Possible values: asc, desc. Default: desc"
// @Param data body pollIDsRequestBody false "body json for defined poll ids as request body"
// @Success 200 {array} model.PollResult
// @Security UserAuth
// @Router /polls [get]
func (h ApisHandler) GetPolls(user *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	offsetFilter := getInt64QueryParam(r, "offset")
	limitFilter := getInt64QueryParam(r, "limit")
	orderFilter := getStringQueryParam(r, "order")

	var pollIDs []string
	bodyData, _ := ioutil.ReadAll(r.Body)
	if bodyData != nil {
		var body pollIDsRequestBody
		bodyErr := json.Unmarshal(bodyData, &body)
		if bodyErr == nil {
			pollIDs = body.IDs
		}
	}

	resData, err := h.app.Services.GetPolls(user, pollIDs, nil, offsetFilter, limitFilter, orderFilter, true)
	if err != nil {
		log.Printf("Error on apis.GetPolls(): %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	result := []model.PollResult{}
	if len(resData) > 0 {
		for _, entry := range resData {
			if entry.UserHasAccess(user.Subject) {
				result = append(result, entry.ToPollResult(user.Subject))
			}
		}
	}

	data, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error on apis.GetPolls(): %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetUserPolls Retrieves  all user poll that may include additional filter params
// @Description Retrieves  all user poll that may include additional filter params
// @Tags Client
// @ID GetUserPolls
// @Param offset query string false "offset"
// @Param limit query string false "limit - limit the result"
// @Param order query string false "order - Possible values: asc, desc. Default: desc"
// @Param data body pollIDsRequestBody false "body json for defined poll ids as request body"
// @Success 200 {array} model.PollResult
// @Security UserAuth
// @Router /user/polls [get]
func (h ApisHandler) GetUserPolls(user *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	offsetFilter := getInt64QueryParam(r, "offset")
	limitFilter := getInt64QueryParam(r, "limit")
	orderFilter := getStringQueryParam(r, "order")

	var pollIDs []string
	bodyData, _ := ioutil.ReadAll(r.Body)
	if bodyData != nil {
		var body pollIDsRequestBody
		bodyErr := json.Unmarshal(bodyData, &body)
		if bodyErr == nil {
			pollIDs = body.IDs
		}
	}

	resData, err := h.app.Services.GetPolls(user, pollIDs, &user.Subject, offsetFilter, limitFilter, orderFilter, true)
	if err != nil {
		log.Printf("Error on apis.GetPolls(): %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	result := []model.PollResult{}
	if len(resData) > 0 {
		for _, entry := range resData {
			if entry.UserHasAccess(user.Subject) {
				result = append(result, entry.ToPollResult(user.Subject))
			}
		}
	}

	data, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error on apis.GetPolls(): %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetPoll Retrieves a poll by id
// @Description Retrieves a poll by id
// @Tags Client
// @ID GetPoll
// @Accept json
// @Produce json
// @Success 200 {object} model.Poll
// @Failure 401
// @Security UserAuth
// @Router /polls/{id} [get]
func (h ApisHandler) GetPoll(user *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resData, err := h.app.Services.GetPoll(user, id)
	if err != nil {
		log.Printf("Error on apis.GetPoll(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if resData == nil || !resData.UserHasAccess(user.Subject) {
		log.Printf("Error on apis.GetPoll(%s): access denied", id)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	data, err := json.Marshal(resData.ToPollResult(user.Subject))
	if err != nil {
		log.Printf("Error on apis.GetPoll(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// UpdatePoll Updates a reward type with the specified id
// @Description Updates a reward type with the specified id
// @Tags Client
// @ID UpdatePoll
// @Param data body model.Poll true "body json"
// @Accept json
// @Produce json
// @Success 200 {object} model.Poll
// @Failure 401
// @Security UserAuth
// @Router /polls/{id} [put]
func (h ApisHandler) UpdatePoll(user *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resData, err := h.app.Services.GetPoll(user, id)
	if err != nil {
		log.Printf("Error on apis.UpdatePoll(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if resData == nil || !resData.UserHasAccess(user.Subject) {
		log.Printf("Error on apis.UpdatePoll(%s): access denied", id)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on apis.UpdatePoll(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.Poll
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on apis.UpdatePoll(%s): %s", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resData, err = h.app.Services.UpdatePoll(user, item)
	if err != nil {
		log.Printf("Error on apis.UpdatePoll(%s): %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(resData.ToPollResult(user.Subject))
	if err != nil {
		log.Printf("Error on apis.UpdatePoll(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// CreatePoll Create a new poll
// @Description Create a new poll
// @Tags Client
// @ID CreatePoll
// @Param data body model.Poll true "body json"
// @Accept json
// @Success 200 {object} model.Poll
// @Security UserAuth
// @Router /polls [post]
func (h ApisHandler) CreatePoll(user *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on apis.CreatePoll: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.Poll
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on apis.CreatePoll: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdItem, err := h.app.Services.CreatePoll(user, item)
	if err != nil {
		log.Printf("Error on apis.CreatePoll: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(createdItem.ToPollResult(user.Subject))
	if err != nil {
		log.Printf("Error on apis.CreatePoll: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// DeletePoll Deletes a poll with the specified id
// @Description Deletes a poll with the specified id
// @Tags Client
// @ID DeletePoll
// @Success 200
// @Security UserAuth
// @Router /polls/{id} [delete]
func (h ApisHandler) DeletePoll(user *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resData, err := h.app.Services.GetPoll(user, id)
	if err != nil {
		log.Printf("Error on apis.GetPoll(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if resData == nil || resData.UserID == user.Subject {
		log.Printf("Error on apis.DeletePoll(%s): access denied", id)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	err = h.app.Services.DeletePoll(user, id)
	if err != nil {
		log.Printf("Error on apis.DeletePoll(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// GetPollEvents Subscribes to a poll events as SSE
// @Description  Subscribes to a poll events as SSE
// @Tags Client
// @ID GetPollEvents
// @Produce json
// @Success 200
// @Security UserAuth
// @Router /polls/{id}/events [post]
func (h ApisHandler) GetPollEvents(user *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Connection doesn't support streaming", http.StatusBadRequest)
		return
	}

	resultChan := make(chan map[string]interface{})
	closeChan := make(chan interface{})
	defer close(closeChan)
	defer fmt.Println("Closing channel.")

	go h.app.Services.SubscribeToPoll(user, id, resultChan, closeChan)

	for {
		select {
		case <-closeChan:
			close(resultChan)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			return
		case data := <-resultChan:
			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Printf("Error on apis.GetPollEvents(): %s", err)
			}
			log.Printf(string(jsonData))
			w.Write(jsonData)
			flusher.Flush()
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// VotePoll Votes a poll with the specified id
// @Description  Votes a poll with the specified id
// @Tags Client
// @ID VotePoll
// @Param data body model.PollVote true "body json"
// @Accept json
// @Produce json
// @Success 200
// @Security UserAuth
// @Router /polls/{id}/vote [post]
func (h ApisHandler) VotePoll(user *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resData, err := h.app.Services.GetPoll(user, id)
	if err != nil {
		log.Printf("Error on apis.VotePoll(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if resData == nil || !resData.UserHasAccess(user.Subject) {
		log.Printf("Error on apis.VotePoll(%s): access denied", id)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on apis.VotePoll(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.PollVote
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on apis.VotePoll(%s): %s", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Subject != item.UserID {
		log.Printf("Error on apis.VotePoll(%s): inconsistent user id", id)
		http.Error(w, "inconsistent user id", http.StatusBadRequest)
	}

	err = h.app.Services.VotePoll(user, id, item)
	if err != nil {
		log.Printf("Error on apis.VotePoll(%s): %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// StartPoll Starts an existing poll with the specified id
// @Description  Starts an existing poll with the specified id
// @Tags Client
// @ID StartPoll
// @Accept json
// @Produce json
// @Success 200
// @Security UserAuth
// @Router /polls/{id}/start [post]
func (h ApisHandler) StartPoll(user *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.app.Services.StartPoll(user, id)
	if err != nil {
		log.Printf("Error on apis.StartPoll(%s): %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// EndPoll Finishes an existing poll with the specified id
// @Description  Finishes an existing poll with the specified id
// @Tags Client
// @ID EndPoll
// @Accept json
// @Produce json
// @Success 200
// @Security UserAuth
// @Router /polls/{id}/end [post]
func (h ApisHandler) EndPoll(user *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.app.Services.EndPoll(user, id)
	if err != nil {
		log.Printf("Error on apis.EndPoll(%s): %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
