/**
 * Copyright (c) 2015 Intel Corporation
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
	"net/http"
)

type ServerError struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

func respondWithError(w *http.ResponseWriter, status int, errorMsg string) {
	(*w).WriteHeader(status)
	log.Errorf(errorMsg)
	msg := ServerError{
		Status: status,
		Error:  errorMsg,
	}
	payload, err := json.Marshal(msg)
	if err == nil {
		(*w).Write(payload)
	}
}
