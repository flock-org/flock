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

package core

import (
	"github.com/clusterlink-net/clusterlink/pkg/util"
	"github.com/sirupsen/logrus"

	"github.com/flock-org/flock/relay/pkg/server"
)

var clog = logrus.WithField("component", "relay-core")

var localhost = "127.0.0.1"

// Relay struct defines the properties of the relay
type Relay struct {
	url      string
	DPServer *server.Server
}

// StartRelay starts the main function of the relay
func (r *Relay) StartRelay(parsedCertData *util.ParsedCertData, port string) error {
	r.DPServer = server.NewRelay(parsedCertData)
	// Start a routine to print active connections periodically
	go r.DPServer.MonitorConnections()
	// Start the main relay server
	err := r.DPServer.StartRelaySSLServer(port)

	return err
}

// Init initializes the relay
func (r *Relay) Init(ip, port string, loglevel logrus.Level) {
	r.url = ip + ":" + port
	clog.Logger.SetLevel(loglevel)
	clog.Logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}
