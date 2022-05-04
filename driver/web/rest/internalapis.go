package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"polls/core"
	"polls/core/model"
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
