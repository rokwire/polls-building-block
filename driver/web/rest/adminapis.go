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
	"io"
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

	data, err := io.ReadAll(r.Body)
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

	data, err := io.ReadAll(r.Body)
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

// GetAlertContacts Retrieves all alert contacts
// @Description Retrieves all alert contacts
// @Tags Admin
// @ID GetAlertContacts
// @Accept json
// @Produce json
// @Success 200 {object} model.AlertContact
// @Failure 401
// @Security UserAuth
// @Router /alert-contacts [get]
func (h AdminApisHandler) GetAlertContacts(user *model.User, w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	resData, err := h.app.Services.GetAlertContacts(user)
	if err != nil {
		log.Printf("Error on apis.GetAlertContact(%s): %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resData == nil {
		log.Printf("Error on apis.GetAlertContact(%s): not found", id)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	data, err := json.Marshal(resData)
	if err != nil {

		log.Printf("Error on apis.GetAlertContact(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetAlertContact Retrieves a alert contact by id
// @Description Retrieves a alert contact by id
// @Tags Admin
// @ID GetAlertContact
// @Accept json
// @Produce json
// @Success 200 {object} model.AlertContact
// @Failure 401
// @Security UserAuth
// @Router /alert-contacts/{id} [get]
func (h AdminApisHandler) GetAlertContact(user *model.User, w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	resData, err := h.app.Services.GetAlertContact(user, id)
	if err != nil {
		log.Printf("Error on apis.GetAlertContact(%s): %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resData == nil {
		log.Printf("Error on apis.GetAlertContact(%s): not found", id)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	data, err := json.Marshal(resData)
	if err != nil {

		log.Printf("Error on apis.GetAlertContact(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// CreateAlertContact Create a new alert contact
// @Description Create a new alert contact
// @Tags Admin
// @ID CreateAlertContact
// @Param data body model.AlertContact true "body json"
// @Accept json
// @Success 200 {object} model.AlertContact
// @Security UserAuth
// @Router /alert-contacts [post]
func (h AdminApisHandler) CreateAlertContact(user *model.User, w http.ResponseWriter, r *http.Request) {

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on apis.CreateAlertContact: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.AlertContact
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on apis.CreateAlertContact: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdItem, err := h.app.Services.CreateAlertContact(user, item)
	if err != nil {
		log.Printf("Error on apis.CreateAlertContact: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(createdItem)
	if err != nil {
		log.Printf("Error on apis.CreateAlertContact: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// UpdateAlertContact Updates an alert contact with the specified id
// @Description Updates an alert contact with the specified id
// @Tags Admin
// @ID UpdateAlertContact
// @Param data body model.AlertContact true "body json"
// @Accept json
// @Produce json
// @Success 200 {object} model.AlertContact
// @Failure 401
// @Security UserAuth
// @Router /alert-contacts/{id} [put]
func (h AdminApisHandler) UpdateAlertContact(user *model.User, w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	data, err := io.ReadAll(r.Body)
	if err != nil {

		log.Printf("Error on apis.UpdateAlertContacts(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.AlertContact
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on apis.UpdateAlertContact(%s): %s", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.app.Services.UpdateAlertContact(user, id, item)
	if err != nil {
		log.Printf("Error on apis.UpdateAlertContact(%s): %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// DeleteAlertContact Deletes an alert contact with the specified id
// @Description Deletes an alert contact with the specified id
// @Tags Admin
// @ID DeleteAlertContact
// @Success 200
// @Security UserAuth
// @Router /alert-contact/{id} [delete]
func (h AdminApisHandler) DeleteAlertContact(user *model.User, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.app.Services.DeleteAlertContact(user, id)
	if err != nil {
		log.Printf("Error on apis.DeleteAlertContact(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
