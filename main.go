/**
 * Copyright (c) 2016 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package main app-dependency-discoverer API
//
// This application can discover dependencies of application stack with provided root GUID.
//
//     Version: 0.2.1
//
// swagger:meta
package main

import (
	log "github.com/cihub/seelog"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/trustedanalytics/app-dependency-discoverer/logging"
	"github.com/trustedanalytics/app-dependency-discoverer/server"
)

func main() {
	logging.Initialize()

	cfEnv, err := cfenv.Current()
	if err != nil {
		log.Warnf("CF Env vars gathering failed with error [%v]. Running locally, probably.", err)
	}

	config := server.Config{}
	config.Initialize(cfEnv)

	server.Start(config)
}
