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

package web

import (
	"fmt"
	"log"
	"net/http"
	"polls/core"
	"polls/core/model"
	"polls/driver/web/rest"
	"polls/utils"
	"strings"

	"github.com/rokwire/core-auth-library-go/tokenauth"

	"github.com/casbin/casbin"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Adapter entity
type Adapter struct {
	host          string
	port          string
	auth          *Auth
	authorization *casbin.Enforcer

	apisHandler         rest.ApisHandler
	adminApisHandler    rest.AdminApisHandler
	internalApisHandler rest.InternalApisHandler

	app *core.Application
}

// @title Polls Building Block v2 API
// @description RoRewards Building Block API Documentation.
// @version 1.0.21
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost
// @BasePath /content
// @schemes https

// @securityDefinitions.apikey InternalApiAuth
// @in header (add INTERNAL-API-KEY with correct value as a header)
// @name Authorization

// @securityDefinitions.apikey AdminUserAuth
// @in header (add Bearer prefix to the Authorization value)
// @name Authorization

// @securityDefinitions.apikey AdminGroupAuth
// @in header
// @name GROUP

// Start starts the module
func (we Adapter) Start() {

	router := mux.NewRouter().StrictSlash(true)
	router.Use()

	subrouter := router.PathPrefix("/polls").Subrouter()
	subrouter.PathPrefix("/doc/ui").Handler(we.serveDocUI())
	subrouter.HandleFunc("/doc", we.serveDoc)
	subrouter.HandleFunc("/version", we.wrapFunc(we.apisHandler.Version)).Methods("GET")

	// handle apis
	apiRouter := subrouter.PathPrefix("/api").Subrouter()

	// Client APIs
	apiRouter.HandleFunc("/polls", we.userAuthWrapFunc(we.apisHandler.GetPolls)).Methods("GET")
	apiRouter.HandleFunc("/polls", we.userAuthWrapFunc(we.apisHandler.CreatePoll)).Methods("POST")
	apiRouter.HandleFunc("/polls/{id}", we.userAuthWrapFunc(we.apisHandler.GetPoll)).Methods("GET")
	apiRouter.HandleFunc("/polls/{id}", we.userAuthWrapFunc(we.apisHandler.UpdatePoll)).Methods("PUT")
	apiRouter.HandleFunc("/polls/{id}", we.userAuthWrapFunc(we.apisHandler.DeletePoll)).Methods("DELETE")
	apiRouter.HandleFunc("/polls/{id}/events", we.userAuthWrapFunc(we.apisHandler.GetPollEvents)).Methods("GET")
	apiRouter.HandleFunc("/polls/{id}/vote", we.userAuthWrapFunc(we.apisHandler.VotePoll)).Methods("PUT")
	apiRouter.HandleFunc("/polls/{id}/start", we.userAuthWrapFunc(we.apisHandler.StartPoll)).Methods("PUT")
	apiRouter.HandleFunc("/polls/{id}/end", we.userAuthWrapFunc(we.apisHandler.EndPoll)).Methods("PUT")
	apiRouter.HandleFunc("/surveys/{id}", we.userAuthWrapFunc(we.apisHandler.GetSurvey)).Methods("GET")
	apiRouter.HandleFunc("/surveys", we.userAuthWrapFunc(we.apisHandler.CreateSurvey)).Methods("POST")
	apiRouter.HandleFunc("/surveys/{id}", we.userAuthWrapFunc(we.apisHandler.UpdateSurvey)).Methods("PUT")
	apiRouter.HandleFunc("/surveys/{id}", we.userAuthWrapFunc(we.apisHandler.DeleteSurvey)).Methods("DELETE")
	apiRouter.HandleFunc("/survey-responses/{id}", we.userAuthWrapFunc(we.apisHandler.GetSurveyResponse)).Methods("GET")
	apiRouter.HandleFunc("/survey-responses", we.userAuthWrapFunc(we.apisHandler.GetSurveyResponses)).Methods("GET")
	apiRouter.HandleFunc("/survey-responses", we.userAuthWrapFunc(we.apisHandler.CreateSurveyResponse)).Methods("POST")
	apiRouter.HandleFunc("/survey-responses/{id}", we.userAuthWrapFunc(we.apisHandler.UpdateSurveyResponse)).Methods("PUT")
	apiRouter.HandleFunc("/survey-responses/{id}", we.userAuthWrapFunc(we.apisHandler.DeleteSurveyResponse)).Methods("DELETE")

	// handle admin apis
	adminRouter := apiRouter.PathPrefix("/admin").Subrouter()

	adminRouter.HandleFunc("/surveys/{id}", we.adminAuthWrapFunc(we.adminApisHandler.GetSurvey)).Methods("GET")
	adminRouter.HandleFunc("/surveys", we.adminAuthWrapFunc(we.adminApisHandler.CreateSurvey)).Methods("POST")
	adminRouter.HandleFunc("/surveys/{id}", we.adminAuthWrapFunc(we.adminApisHandler.UpdateSurvey)).Methods("PUT")
	adminRouter.HandleFunc("/surveys/{id}", we.adminAuthWrapFunc(we.adminApisHandler.DeleteSurvey)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":"+we.port, router))
}

func (we Adapter) serveDoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("access-control-allow-origin", "*")
	http.ServeFile(w, r, "./driver/web/docs/gen/def.yaml")
}

func (we Adapter) serveDocUI() http.Handler {
	url := fmt.Sprintf("%s/polls/doc", we.host)
	return httpSwagger.Handler(httpSwagger.URL(url))
}

func (we Adapter) wrapFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		handler(w, req)
	}
}

type apiKeysAuthFunc = func(http.ResponseWriter, *http.Request)

func (we Adapter) apiKeyOrTokenWrapFunc(handler apiKeysAuthFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		// apply core token check
		coreAuth, _ := we.auth.coreAuth.Check(req)
		if coreAuth {
			handler(w, req)
			return
		}

		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

type authFunc = func(*model.User, http.ResponseWriter, *http.Request)

func (we Adapter) userAuthWrapFunc(handler authFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		coreAuth, user := we.auth.coreAuth.Check(req)
		if coreAuth && user != nil && !user.Claims.Anonymous {
			handler(user, w, req)
			return
		}
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

// TODO: Switch to Core BB model for auth
func (we Adapter) adminAuthWrapFunc(handler authFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		obj := req.URL.Path // the resource that is going to be accessed.
		act := req.Method   // the operation that the user performs on the resource.

		coreAuth, user := we.auth.coreAuth.Check(req)
		if coreAuth {
			permissions := strings.Split(user.Claims.Permissions, ",")

			HasAccess := false
			for _, s := range permissions {
				HasAccess = we.authorization.Enforce(s, obj, act)
				if HasAccess {
					break
				}
			}
			if HasAccess {
				handler(user, w, req)
				return
			}
			log.Printf("Access control error - Core Subject: %s is trying to apply %s operation for %s\n", user.Claims.Subject, act, obj)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

type internalAPIKeyAuthFunc = func(http.ResponseWriter, *http.Request)

func (we Adapter) internalAPIKeyAuthWrapFunc(handler internalAPIKeyAuthFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		apiKeyAuthenticated := we.auth.internalAuth.check(w, req)

		if apiKeyAuthenticated {
			handler(w, req)
		} else {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

// NewWebAdapter creates new WebAdapter instance
func NewWebAdapter(host string, port string, app *core.Application, tokenAuth *tokenauth.TokenAuth, config *model.Config) Adapter {
	auth := NewAuth(app, config, tokenAuth)
	authorization := casbin.NewEnforcer("driver/web/authorization_model.conf", "driver/web/authorization_policy.csv")

	apisHandler := rest.NewApisHandler(app, config)
	adminApisHandler := rest.NewAdminApisHandler(app, config)
	internalApisHandler := rest.NewInternalApisHandler(app, config)
	return Adapter{
		host:                host,
		port:                port,
		auth:                auth,
		authorization:       authorization,
		apisHandler:         apisHandler,
		adminApisHandler:    adminApisHandler,
		internalApisHandler: internalApisHandler,
		app:                 app,
	}
}

// AppListener implements core.ApplicationListener interface
type AppListener struct {
	adapter *Adapter
}
