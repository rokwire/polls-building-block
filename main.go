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

	"github.com/rokwire/rokwire-building-block-sdk-go/services/core/auth"
	"github.com/rokwire/rokwire-building-block-sdk-go/services/core/auth/keys"
	"github.com/rokwire/rokwire-building-block-sdk-go/services/core/auth/sigauth"
	"github.com/rokwire/rokwire-building-block-sdk-go/utils/logging/logs"
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

	// Groups BB Host
	groupsBBHost := getEnvKey("POLLS_GROUPS_BB_HOST", true)

	// Notifications BB Host
	notificationsBBHost := getEnvKey("POLLS_NOTIFICATIONS_BB_HOST", true)

	authService := auth.Service{
		ServiceID:   serviceID,
		ServiceHost: serviceURL,
		FirstParty:  true,
		AuthBaseURL: coreBBHost,
	}

	serviceRegLoader, err := auth.NewRemoteServiceRegLoader(&authService, []string{"auth"})
	if err != nil {
		log.Fatalf("Error initializing remote service registration loader: %v", err)
	}

	serviceRegManager, err := auth.NewServiceRegManager(&authService, serviceRegLoader, !strings.HasPrefix(serviceURL, "http://localhost"))
	if err != nil {
		log.Fatalf("Error initializing service registration manager: %v", err)
	}

	//core adapter
	var serviceAccountManager *auth.ServiceAccountManager

	serviceAccountID := getEnvKey("POLLS_SERVICE_ACCOUNT_ID", false)
	privKeyRaw := getEnvKey("POLLS_PRIV_KEY", true)
	privKeyRaw = strings.ReplaceAll(privKeyRaw, "\\n", "\n")
	privKey, err := keys.NewPrivKey(keys.PS256, privKeyRaw)
	if err != nil {
		logger.Errorf("Error parsing priv key: %v", err)
	} else if serviceAccountID == "" {
		logger.Errorf("Missing service account id")
	} else {
		signatureAuth, err := sigauth.NewSignatureAuth(privKey, serviceRegManager, false, false)
		if err != nil {
			logger.Fatalf("Error initializing signature auth: %v", err)
		}

		serviceAccountLoader, err := auth.NewRemoteServiceAccountLoader(&authService, serviceAccountID, signatureAuth)
		if err != nil {
			logger.Fatalf("Error initializing remote service account loader: %v", err)
		}

		serviceAccountManager, err = auth.NewServiceAccountManager(&authService, serviceAccountLoader)
		if err != nil {
			logger.Fatalf("Error initializing service account manager: %v", err)
		}
	}
	config := &model.Config{
		MongoDBAuth:       mongoDBAuth,
		MongoDBName:       mongoDBName,
		MongoTimeout:      mongoTimeout,
		InternalAPIKey:    internalAPIKey, // pragma: allowlist secret
		CoreBBHost:        coreBBHost,
		PollServiceURL:    serviceURL,
		UiucOrgID:         uiucOrgID,
		GroupsHost:        groupsBBHost,
		NotificationsHost: notificationsBBHost,
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
	notificationsBBAdapter := notifications.NewNotificationsAdapter(notificationsBBHost, internalAPIKey, appID, orgID)

	groupsAdapter := groups.NewGroupsAdapter(config, serviceAccountManager)

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
			log.Print("No provided environment variable for " + key)
		}
	}
	return value
}
