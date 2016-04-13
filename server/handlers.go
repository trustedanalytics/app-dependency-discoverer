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

package server

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	"github.com/go-martini/martini"
	"github.com/trustedanalytics/app-dependency-discoverer/graph"
	"net/http"
)

type Handlers struct{}

func (*Handlers) Discover(w http.ResponseWriter, r *http.Request, params martini.Params) {
	if _, ok := params["rootGUID"]; !ok {
		respondWithError(&w, http.StatusBadRequest, "No root GUID provided")
		return
	}

	api := graph.NewGraphAPI()
	result, err := api.Discover(params["rootGUID"])
	if err != nil {
		respondWithError(&w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Debugf("Sent: %v", result)

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(result)
}
