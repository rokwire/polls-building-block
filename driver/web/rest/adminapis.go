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

package rest

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"polls/core"
	"polls/core/model"

	"github.com/gorilla/mux"
)

// AdminApisHandler handles the rest Admin APIs implementation
type AdminApisHandler struct {
	app    *core.Application
	config *model.Config
}

// GetSurvey Retrieves a Survey by id
// @Description Retrieves a Survey by id
// @Tags Admin
// @ID GetSurvey
// @Accept json
// @Produce json
// @Success 200 {object} model.Survey
// @Failure 401
// @Security UserAuth
// @Router /surveys/{id} [get]
func (h AdminApisHandler) GetSurvey(user *model.User, w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	resData, err := h.app.Services.GetSurvey(user, id)
	if err != nil {
		log.Printf("Error on apis.GetSurvey(%s): %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resData == nil {
		log.Printf("Error on apis.GetSurvey(%s): not found", id)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	data, err := json.Marshal(resData)
	if err != nil {

		log.Printf("Error on apis.GetSurvey(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// CreateSurvey Create a new survey
// @Description Create a new survey
// @Tags Admin
// @ID CreateSurvey
// @Param data body model.Survey true "body json"
// @Accept json
// @Success 200 {object} model.Survey
// @Security UserAuth
// @Router /surveys [post]
func (h AdminApisHandler) CreateSurvey(user *model.User, w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on apis.CreateSurvey: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.Survey
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on apis.CreateSurvey: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdItem, err := h.app.Services.CreateSurvey(user, item, true)
	if err != nil {
		log.Printf("Error on apis.CreateSurvey: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(createdItem)
	if err != nil {
		log.Printf("Error on apis.CreateSurvey: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// UpdateSurvey Updates a survey type with the specified id
// @Description Updates a survey type with the specified id
// @Tags Admin
// @ID UpdateSurvey
// @Param data body model.Survey true "body json"
// @Accept json
// @Produce json
// @Success 200 {object} model.Survey
// @Failure 401
// @Security UserAuth
// @Router /surveys/{id} [put]
func (h AdminApisHandler) UpdateSurvey(user *model.User, w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {

		log.Printf("Error on apis.UpdateSurvey(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.Survey
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on apis.UpdateSurvey(%s): %s", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.app.Services.UpdateSurvey(user, item, id, true)
	if err != nil {
		log.Printf("Error on apis.UpdateSurvey(%s): %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// DeleteSurvey Deletes a survey with the specified id
// @Description Deletes a survey with the specified id
// @Tags Admin
// @ID DeleteSurvey
// @Success 200
// @Security UserAuth
// @Router /surveys/{id} [delete]
func (h AdminApisHandler) DeleteSurvey(user *model.User, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.app.Services.DeleteSurvey(user, id, true)
	if err != nil {
		log.Printf("Error on apis.DeleteSurvey(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
