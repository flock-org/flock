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
	"encoding/json"
	"net"

	"github.com/praveingk/openssl"

	"github.com/flock-org/flock/relay/pkg/api"
)

func (s *Server) authorize(srcParty string, tcpConn net.Conn, tlsConn *openssl.Conn) error {
	authReq := api.AuthReq{}
	buf := make([]byte, 512)

	n, err := tlsConn.Read(buf)
	if err != nil {
		s.logger.Errorf("Failed to read auth request: %s", err)
		return err
	}
	err = json.Unmarshal(buf[:n], &authReq)
	if err != nil {
		s.logger.Errorf("Failed to unmarshal auth request: %s", err)
		return err
	}

	dstConn, err := s.states.GetConnection(authReq.DestParty, srcParty, authReq.Tag)
	if err != nil {
		//s.logger.Infof("Destination party doesnt have an active connection, Waiting")
		//s.logger.Infof("Storing the tcp connection for %s:%s (%s) comms", srcParty, authReq.DestParty, authReq.Tag)
		s.states.StoreTLSConnection(srcParty, authReq.DestParty, authReq.Tag, tlsConn)
		err = s.states.StoreConnection(srcParty, authReq.DestParty, authReq.Tag, tcpConn)
		if err != nil {
			return err
		}
		// go func() {
		// 	connKey := connection{party1: srcParty, party2: authReq.DestParty, tag: authReq.Tag}
		// 	s.timelineMutex.Lock()
		// 	s.relayTimeline[connKey] = &timeline{party1Incoming: incomingConn[tcpConn.RemoteAddr().String()],
		// 		party1Handshake: partyHandshake, party1Auth: time.Now()}
		// 	s.timelineMutex.Unlock()
		// }()
		return nil
	}
	//party2Auth := time.Now()

	destTLSConn, err := s.states.GetTLSConnection(authReq.DestParty, srcParty, authReq.Tag)
	if err != nil {
		return err
	}

	//s.logger.Infof("Ending the TLS Connections(%s, %s, %s) and start TCP forwarding", authReq.DestParty, srcParty, authReq.Tag)
	err = s.sendReady(tlsConn, api.Ready{Mode: api.TLSModeServer})
	if err != nil {
		tlsConn.Close()
		destTLSConn.Close()
		s.states.RemoveConnection(srcParty, authReq.DestParty, authReq.Tag)
		s.states.RemoveConnection(authReq.DestParty, srcParty, authReq.Tag)
		s.states.RemoveTLSConnection(srcParty, authReq.DestParty, authReq.Tag)
		s.states.RemoveTLSConnection(authReq.DestParty, srcParty, authReq.Tag)
		return err
	}
	// To synchronize the TLS connections, we wait for the ACK and proceed to next server
	err = s.sendReady(destTLSConn, api.Ready{Mode: api.TLSModeClient})
	if err != nil {
		tlsConn.Close()
		destTLSConn.Close()
		s.states.RemoveConnection(srcParty, authReq.DestParty, authReq.Tag)
		s.states.RemoveConnection(authReq.DestParty, srcParty, authReq.Tag)
		s.states.RemoveTLSConnection(srcParty, authReq.DestParty, authReq.Tag)
		s.states.RemoveTLSConnection(authReq.DestParty, srcParty, authReq.Tag)
		return err
	}
	go s.startForwarding(srcParty, tcpConn, tlsConn, authReq.DestParty, dstConn, destTLSConn, authReq.Tag)
	// go func() {
	// 	connKey := connection{party1: authReq.DestParty, party2: srcParty, tag: authReq.Tag}
	// 	s.timelineMutex.Lock()
	// 	s.relayTimeline[connKey].SetupConn = time.Now()
	// 	s.relayTimeline[connKey].party2Auth = party2Auth
	// 	s.relayTimeline[connKey].party2Incoming = incomingConn[tcpConn.RemoteAddr().String()]
	// 	s.relayTimeline[connKey].party2Handshake = partyHandshake
	// 	s.timelineMutex.Unlock()
	// }()

	return nil
}

func (s *Server) sendReady(conn *openssl.Conn, ready api.Ready) error {
	buf := make([]byte, 512)

	readyData, err := json.Marshal(ready)
	if err != nil {
		s.logger.Errorf("Failed to marshal ready response: %v.", err)
	}
	_, err = conn.Write(readyData)
	if err != nil {
		s.logger.Errorf("Failed to write send ready: %s", err)
		return err
	}
	conn.CloseSSL()
	_, err = conn.Read(buf)
	//s.logger.Infof("Ready and closed SSL")
	return nil
}
