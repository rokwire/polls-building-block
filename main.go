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

package main

import (
	"github.com/rokwire/core-auth-library-go/authservice"
	"github.com/rokwire/core-auth-library-go/tokenauth"
	"github.com/rokwire/logging-library-go/logs"
	"log"
	"os"
	"polls/core"
	"polls/core/model"
	cacheadapter "polls/driven/cache"
	"polls/driven/groups"
	"polls/driven/notifications"
	storage "polls/driven/storage"
	driver "polls/driver/web"
	"strings"
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

	remoteConfig := authservice.RemoteAuthDataLoaderConfig{
		AuthServicesHost: coreBBHost,
	}

	serviceLoader, err := authservice.NewRemoteAuthDataLoader(remoteConfig, []string{"core", "notifications", "groups"}, logs.NewLogger("polls-v2", &logs.LoggerOpts{}))
	if err != nil {
		log.Fatalf("Error initializing auth service: %v", err)
	}

	authService, err := authservice.NewAuthService("polls-v2", serviceURL, serviceLoader)
	if err != nil {
		log.Fatalf("Error initializing auth service: %v", err)
	}

	tokenAuth, err := tokenauth.NewTokenAuth(true, authService, nil, nil)
	if err != nil {
		log.Fatalf("Error intitializing token auth: %v", err)
	}

	// Notifications service reg
	notificationsServiceReg, err := authService.GetServiceReg("notifications")
	if err != nil {
		log.Fatalf("error finding notifications service reg: %s", err)
	}

	// Groups service reg
	groupsServiceReg, err := authService.GetServiceReg("groups")
	if err != nil {
		log.Fatalf("error finding notifications service reg: %s", err)
	}

	config := &model.Config{
		MongoDBAuth:       mongoDBAuth,
		MongoDBName:       mongoDBName,
		MongoTimeout:      mongoTimeout,
		InternalAPIKey:    internalAPIKey,
		CoreBBHost:        coreBBHost,
		PollServiceURL:    serviceURL,
		UiucOrgID:         uiucOrgID,
		GroupsHost:        groupsServiceReg.Host,
		NotificationsHost: notificationsServiceReg.Host,
	}

	storageAdapter := storage.NewStorageAdapter(config)
	err = storageAdapter.Start()
	if err != nil {
		log.Fatal("Cannot start the mongoDB adapter - " + err.Error())
	} else {
		log.Printf("Storage started")
	}

	notificationsAdapter := notifications.NewNotificationsAdapter(config)

	groupsAdapter := groups.NewGroupsAdapter(config)

	defaultCacheExpirationSeconds := getEnvKey("DEFAULT_CACHE_EXPIRATION_SECONDS", false)
	cacheAdapter := cacheadapter.NewCacheAdapter(defaultCacheExpirationSeconds)

	// application
	application := core.NewApplication(Version, Build, storageAdapter, cacheAdapter, notificationsAdapter, groupsAdapter)
	application.Start()

	webAdapter := driver.NewWebAdapter(host, port, application, tokenAuth, config)

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
