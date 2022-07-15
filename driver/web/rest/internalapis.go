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
	"log"
	"net/http"
	"polls/core"
	"polls/core/model"

	"github.com/gorilla/mux"
)

// InternalApisHandler handles the rest internal APIs implementation
type InternalApisHandler struct {
	app    *core.Application
	config *model.Config
}

// GetGroupPolls Retrieves poll id to group id mapping
// @Description  Retrieves poll id to group id mapping
// @Tags Client
// @ID GetGroupPolls
// @Produce json
// @Success 200
// @Security UserAuth
// @Router /polls/{id}/vote [post]
func (h InternalApisHandler) GetGroupPolls(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	flterGroupPolls := true
	groupPolls, err := h.app.Services.GetPolls(nil, model.PollsFilter{GroupPolls: &flterGroupPolls}, false)
	if err != nil {
		log.Printf("Error on internalapis.GetGroupPolls(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	groupMapping := map[string]string{}
	if len(groupPolls) > 0 {
		for _, poll := range groupPolls {
			groupMapping[poll.ID.Hex()] = poll.GroupID
		}
	}

	resData, err := json.Marshal(groupMapping)
	if err != nil {
		log.Printf("Error on internalapis.GetGroupPolls(): %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resData)
}
