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

package corebb

import (
	"github.com/rokwire/core-auth-library-go/v2/authservice"
	"github.com/rokwire/logging-library-go/v2/logs"
)

// Adapter implements the Core interface
type Adapter struct {
	logger                logs.Logger
	coreURL               string
	serviceAccountManager *authservice.ServiceAccountManager

	appID string
	orgID string
}

// NewCoreAdapter creates a new adapter for Core API
func NewCoreAdapter(coreURL string, orgID string, appID string, serviceAccountManager *authservice.ServiceAccountManager) *Adapter {
	return &Adapter{coreURL: coreURL, appID: appID, orgID: orgID, serviceAccountManager: serviceAccountManager}
}
