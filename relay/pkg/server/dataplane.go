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
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/praveingk/openssl"
	"github.com/sirupsen/logrus"

	cutil "github.com/clusterlink-net/clusterlink/pkg/util"
	"github.com/flock-org/flock/relay/config"
	"github.com/flock-org/flock/relay/pkg/store"
)

var (
	maxDataBufferSize = 64 * 1024
)

// Server contains the declaration of the relay server
type Server struct {
	router         *chi.Mux
	parsedCertData *cutil.ParsedCertData
	states         *store.State
	logger         *logrus.Entry
	f1             *os.File
	f2             *os.File
}

type connection struct {
	party1 string
	party2 string
	tag    string
}

func (s *Server) startForwarding(srcParty string, conn1 net.Conn, sslConn1 *openssl.Conn, dstParty string, conn2 net.Conn, sslConn2 *openssl.Conn, tag string) {
	forwarder := newForwarder(conn1, conn2)
	b1, b2 := forwarder.run()
	s.logger.Infof("Forwarding finished for %s:%s(%s), bytes transferred(%d, %d)", srcParty, dstParty, tag, b1, b2)
	sslConn1.Close()
	sslConn2.Close()
	s.states.RemoveConnection(srcParty, dstParty, tag)
	s.states.RemoveConnection(dstParty, srcParty, tag)
	s.states.RemoveTLSConnection(srcParty, dstParty, tag)
	s.states.RemoveTLSConnection(dstParty, srcParty, tag)
}

func (s *Server) receiveWaitAndForward(address string, ctx *openssl.Ctx) error {
	acceptor, err := net.Listen("tcp", address)
	if err != nil {
		s.logger.Errorln("Error:", err)
		return err
	}

	for {
		tcpConn, err := acceptor.Accept()
		if err != nil {
			s.logger.Errorln("Accept error:", err)
			continue
		}
		tlsConn, err := openssl.Server(tcpConn, ctx)
		err = tlsConn.Handshake()
		if err != nil {
			s.logger.Errorf("Handshake failed: %v.", err)
			tlsConn.Close()
			continue
		}
		s.logger.Info("Accept incoming connection from ", tlsConn.RemoteAddr().String())
		reqParty, err := getPartyName(tlsConn)
		if err != nil {
			s.logger.Errorf("Failed to get party name: %v.", err)
			tlsConn.Close()
			continue
		}
		s.logger.Infof("Got connection from %s requesting access to %s", reqParty, tlsConn.GetServername())

		err = s.authorize(reqParty, tcpConn, tlsConn)
		if err != nil {
			s.logger.Errorf("Failed to authorize %s; %v", reqParty, err)
			tlsConn.Close()
			continue
		}
	}
}

// StartRelaySSLServer starts the Flock relay dataplane server which listens to connections from the user's parties (or functions)
func (s *Server) StartRelaySSLServer(port string) error {
	defer s.f1.Close()
	defer s.f2.Close()
	address := fmt.Sprintf(":%s", port)
	s.logger.Info("Starting relay... Listening to ", address, " for connections")
	ctx, err := openssl.NewCtxFromFiles(filepath.Join(config.FlockrelayCADirectory(), config.CertificateFileName),
		filepath.Join(config.FlockrelayCADirectory(), config.PrivateKeyFileName))
	if err != nil {
		s.logger.Fatal(err)
	}
	ctx.LoadVerifyLocations(config.FrCAFileRoot, "")
	ctx.SetVerifyMode(openssl.VerifyPeer)

	return s.receiveWaitAndForward(address, ctx)
}

// MonitorConnections prints the active connection periodically
func (s *Server) MonitorConnections() {
	for {
		s.logger.Infof("Active Connections : %d", s.states.Conns())
		time.Sleep(5 * time.Second)
	}
}

// NewRelay returns a new dataplane HTTP server.
func NewRelay(parsedCertData *cutil.ParsedCertData) *Server {
	s := &Server{
		router:         chi.NewRouter(),
		parsedCertData: parsedCertData,
		states:         store.GetState(),
		logger:         logrus.WithField("component", "server.relay"),
	}
	return s
}
