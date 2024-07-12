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

package main

import (
	"log"
	"os"

	"polls/core"
	"polls/core/model"
	cacheadapter "polls/driven/cache"
	corebb "polls/driven/core"
	"polls/driven/groups"
	"polls/driven/notifications"
	storage "polls/driven/storage"
	driver "polls/driver/web"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/rokwire/core-auth-library-go/v2/authservice"
	"github.com/rokwire/core-auth-library-go/v2/sigauth"
	"github.com/rokwire/logging-library-go/v2/logs"
)

var (
	// Version : version of this executable
	Version string
	// Build : build date of this executable
	Build string
)

func main() {
	if len(Version) == 0 {
		Version = "dev"
	}

	serviceID := "polls-v2"

	loggerOpts := logs.LoggerOpts{SuppressRequests: logs.NewStandardHealthCheckHTTPRequestProperties("/polls/version")}
	logger := logs.NewLogger(serviceID, &loggerOpts)

	port := getEnvKey("PORT", true)

	internalAPIKey := getEnvKey("INTERNAL_API_KEY", true)

	//mongoDB adapter
	mongoDBAuth := getEnvKey("MONGO_AUTH", true)
	mongoDBName := getEnvKey("MONGO_DATABASE", true)
	mongoTimeout := getEnvKey("MONGO_TIMEOUT", false)

	// web adapter
	host := getEnvKey("HOST", true)
	coreBBHost := getEnvKey("CORE_BB_HOST", true)
	serviceURL := getEnvKey("POLL_SERVICE_URL", true)
	uiucOrgID := getEnvKey("UIUC_ORG_ID", true)

	authService := authservice.AuthService{
		ServiceID:   serviceID,
		ServiceHost: serviceURL,
		FirstParty:  true,
		AuthBaseURL: coreBBHost,
	}

	serviceRegLoader, err := authservice.NewRemoteServiceRegLoader(&authService, []string{"auth"})
	if err != nil {
		log.Fatalf("Error initializing remote service registration loader: %v", err)
	}

	serviceRegManager, err := authservice.NewServiceRegManager(&authService, serviceRegLoader)
	if err != nil {
		log.Fatalf("Error initializing service registration manager: %v", err)
	}

	//core adapter
	serviceAccountID := getEnvKey("POLLS_SERVICE_ACCOUNT_ID", false)
	privKeyRaw := getEnvKey("POLLS_PRIV_KEY", true)
	privKeyRaw = strings.ReplaceAll(privKeyRaw, "\\n", "\n")
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privKeyRaw))
	if err != nil {
		log.Fatalf("Error parsing priv key: %v", err)
	}
	signatureAuth, err := sigauth.NewSignatureAuth(privKey, serviceRegManager, false)
	if err != nil {
		log.Fatalf("Error initializing signature auth: %v", err)
	}

	serviceAccountLoader, err := authservice.NewRemoteServiceAccountLoader(&authService, serviceAccountID, signatureAuth)
	if err != nil {
		log.Fatalf("Error initializing remote service account loader: %v", err)
	}

	serviceAccountManager, err := authservice.NewServiceAccountManager(&authService, serviceAccountLoader)
	if err != nil {
		log.Fatalf("Error initializing service account manager: %v", err)
	}
	config := &model.Config{
		MongoDBAuth:    mongoDBAuth,
		MongoDBName:    mongoDBName,
		MongoTimeout:   mongoTimeout,
		InternalAPIKey: internalAPIKey,
		CoreBBHost:     coreBBHost,
		PollServiceURL: serviceURL,
		UiucOrgID:      uiucOrgID,
	}

	storageAdapter := storage.NewStorageAdapter(config, logger)
	err = storageAdapter.Start()
	if err != nil {
		log.Fatal("Cannot start the mongoDB adapter - " + err.Error())
	} else {
		log.Printf("Storage started")
	}

	//notifications BB adapter
	appID := getEnvKey("POLLS_APP_ID", true)
	orgID := getEnvKey("POLLS_ORG_ID", true)
	notificationHost := getEnvKey("POLLS_NOTIFICATIONS_BB_HOST", true)
	notificationsBBAdapter := notifications.NewNotificationsAdapter(notificationHost, internalAPIKey, appID, orgID)

	groupsAdapter := groups.NewGroupsAdapter(config)

	defaultCacheExpirationSeconds := getEnvKey("DEFAULT_CACHE_EXPIRATION_SECONDS", false)
	cacheAdapter := cacheadapter.NewCacheAdapter(defaultCacheExpirationSeconds)

	//core adapter
	coreAdapter := corebb.NewCoreAdapter(coreBBHost, orgID, appID, serviceAccountManager)

	// application
	application := core.NewApplication(Version, Build, storageAdapter, cacheAdapter, notificationsBBAdapter,
		groupsAdapter, serviceID, coreAdapter, logger)
	application.Start()

	webAdapter := driver.NewWebAdapter(host, port, application, config, serviceRegManager, logger)

	webAdapter.Start()
}

func getEnvKeyAsList(key string, required bool) []string {
	stringValue := getEnvKey(key, required)

	// it is comma separated format
	stringListValue := strings.Split(stringValue, ",")
	if len(stringListValue) == 0 && required {
		log.Fatalf("missing or empty env var: %s", key)
	}

	return stringListValue
}

func getEnvKey(key string, required bool) string {
	// get from the environment
	value, exist := os.LookupEnv(key)
	if !exist {
		if required {
			log.Fatal("No provided environment variable for " + key)
		} else {
			log.Printf("No provided environment variable for " + key)
		}
	}
	return value
}
