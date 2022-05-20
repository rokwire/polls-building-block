package web

import (
	"github.com/rokwire/core-auth-library-go/tokenauth"
	"log"
	"net/http"
	"polls/core"
	"polls/core/model"
)

// CoreAuth implementation
type CoreAuth struct {
	app       *core.Application
	tokenAuth *tokenauth.TokenAuth
}

// NewCoreAuth creates new CoreAuth
func NewCoreAuth(app *core.Application, tokenAuth *tokenauth.TokenAuth) *CoreAuth {
	auth := CoreAuth{app: app, tokenAuth: tokenAuth}
	return &auth
}

// Check checks the request contains a valid Core access token
func (ca CoreAuth) Check(r *http.Request) (bool, *model.User) {

	claims, err := ca.tokenAuth.CheckRequestTokens(r)
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
	}

	return false, nil
}
