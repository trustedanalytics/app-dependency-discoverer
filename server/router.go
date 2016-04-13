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
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/auth"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
)

const (
	apiVersion = "v1"
)

var (
	discoverURLPattern = fmt.Sprintf("/%v/discover/:rootGUID", apiVersion)
)

type router struct {
	m *martini.ClassicMartini
}

// ServeHTTP logs all requests and dispatches to the appropriate handler
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if dump, err := httputil.DumpRequest(req, true); err != nil {
		log.Tracef("Cannot log incoming request: %v", err)
	} else {
		log.Tracef(string(dump))
	}
	w.Header().Set("Content-Type", "application/json")
	r.m.ServeHTTP(w, req)
}

func Start(config Config) {
	m := martini.Classic()
	m.Use(auth.Basic(GetEnvVarAsString("AUTH_USER", ""), GetEnvVarAsString("AUTH_PASS", "")))

	handlers := Handlers{}
	m.Get(discoverURLPattern, handlers.Discover)

	r := &router{m}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	errCh := make(chan error, 1)

	go func() {
		address := fmt.Sprintf("%v:%v", config.CFEnv.Host, config.CFEnv.Port)
		log.Infof("starting: %v", address)
		errCh <- http.ListenAndServe(address, r)
	}()

	select {
	case err := <-errCh:
		log.Errorf("error: %v", err)
	case sig := <-sigCh:
		var _ = sig
		log.Info("done")
	}
}
