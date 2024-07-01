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

package api

// TLSMode represents the role the party must take for E2E
type TLSMode int

const (
	// TLSModeServer is server
	TLSModeServer TLSMode = 1
	// TLSModeClient is client
	TLSModeClient = 2
)

type UserReq struct {
	// UserName is a unique name associated with the user e.g. alice
	UserName string
	// Parties is the count of the parties/functions to be deployed
	Parties string
}

type UserResp struct {
	// RelayTarget specifies the destionation of the Relays
	RelayTarget string
	PartyInfos  []PartyInfo
}

type PartyInfo struct {
	// PartyIDs specifies the UUIDs to be used for the Parties/functions
	PartyID string
	// Certificates are sent over separately as octets with partyID-cert.pem & pertyID-key.pem
}

// UserSpec contains all the party attributes and access group
type UserSpec struct {
	Party       []string
	AccessGroup string
}

// AuthReq contains the access authorization message sent by a party to the relay
type AuthReq struct {
	DestParty string
	Tag       string // Optional if establishing a specific connection using a tag
}

// Ready contains the message that is sent to party when the connection is ready
type Ready struct {
	Mode TLSMode
}
