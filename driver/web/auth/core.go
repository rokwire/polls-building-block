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
	"net/http"
	"polls/core"
	"polls/core/model"

	"github.com/rokwire/core-auth-library-go/v2/authservice"
	"github.com/rokwire/core-auth-library-go/v2/tokenauth"
)

// CoreAuth implementation
type CoreAuth struct {
	app       *core.Application
	tokenAuth *tokenauth.TokenAuth
}

// Check checks the request contains a valid Core access token
func (ca CoreAuth) Check(r *http.Request) (bool, *model.User) {

	/*	claims, err := ca.tokenAuth.CheckRequestTokens(r)
		if err != nil {
			log.Printf("error validate token: %s", err)
			return false, nil
		}

		if claims != nil {
			if claims.Valid() == nil {
				token, _, _ := tokenauth.GetRequestTokens(r)
				if len(token) > 0 {
					return true, &model.User{
						Token:  token,
						Claims: *claims,
					}
				}

			}
		} */

	return false, nil
}

// NewCoreAuth creates new CoreAuth
func NewCoreAuth(app *core.Application, authService *authservice.AuthService) *CoreAuth {

	//TODO
	//tokenAuth, err := tokenauth.NewTokenAuth(true, authService, nil, nil)

	auth := CoreAuth{app: app}
	return &auth
}
