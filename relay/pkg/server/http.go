// Copyright 2024 The Flock Authors.
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

package server

import (
	"fmt"

	"github.com/clusterlink-net/clusterlink/pkg/util"
	cutil "github.com/clusterlink-net/clusterlink/pkg/util"
	"github.com/clusterlink-net/clusterlink/pkg/utils/netutils"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

type APIServer struct {
	router         *chi.Mux
	parsedCertData *util.ParsedCertData
	logger         *logrus.Entry
}

// StartFlockAPIServer starts the Dataplane server
func (s *APIServer) StartFlockAPIServer() error {
	address := fmt.Sprintf(":%d", apiPort)
	s.logger.Infof("Flock API server starting at %s.", address)
	server := netutils.CreateResilientHTTPServer(address, s.router, s.parsedCertData.ServerConfig(), nil, nil, nil)

	return server.ListenAndServeTLS("", "")
}

func (s *APIServer) addAPIHandlers() {
	s.router.Route("/user", func(r chi.Router) {
		s.router.Get("/", s.addUser)
		s.router.Post("/", s.addUser)
	})
	s.router.Mount("/", s.router)
}

func NewAPIServer(parsedCertData *cutil.ParsedCertData) *APIServer {
	s := &APIServer{
		router:         chi.NewRouter(),
		parsedCertData: parsedCertData,
		logger:         logrus.WithField("component", "server.flockrelay"),
	}

	s.addAPIHandlers()

	return s
}
