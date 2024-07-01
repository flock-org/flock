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

package signing

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/bnb-chain/tss-lib/common"
	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/ecdsa/signing"
	"github.com/bnb-chain/tss-lib/tss"
	"github.com/flock-org/flock/internal/networking"
)

// Conducts multi-party ECDSA signing for n = numParties and t = numThreshold
func SigningParty(partyInt int, numParties int, numThreshold int, comm networking.Communicator, key keygen.LocalPartySaveData, msg *big.Int) string {
	// if err := logger.SetLogLevel("tss-lib", "debug"); err != nil {
	// 	panic(err)
	// }
	msgQueue = make(map[string][]tss.ParsedMessage)

	// Index of next message we must send
	next = 0

	// Number of messages of current type we have processed
	numProcessed = 0

	// Number of other parties
	numOtherParties = numParties - 1

	// Number of bytes received
	totalBytesRead := 0

	// Number of bytes sent
	totalBytesSent := 0

	// List of party IDs
	partyIDs := GetParticipantPartyIDs(numParties)
	ctx := tss.NewPeerContext(partyIDs)
	thisPartyID := partyIDs[partyInt]
	otherPartyIDs := partyIDs.Exclude(thisPartyID)

	// Channels
	errCh := make(chan *tss.Error, 1)
	outCh := make(chan tss.Message, 1)
	endCh := make(chan common.SignatureData, 1)

	// Init the party
	startTime := time.Now()
	var endTime time.Time
	params := tss.NewParameters(tss.S256(), ctx, thisPartyID, numParties, numThreshold)
	party := signing.NewLocalParty(msg, params, key, outCh, endCh).(*signing.LocalParty)
	go func() {
		if err := party.Start(); err != nil {
			errCh <- err
		}
	}()

	go func() {
		select {
		case err := <-errCh:
			log.Fatalf("Error: %s", err)
		}
	}()

	for {
		// Send outgoing messages
		select {
		case msg := <-outCh:
			dest := msg.GetTo()
			if dest == nil { // broadcast!
				for _, partyID := range otherPartyIDs {
					bytesSent := sendTSSMessage(msg, party, *partyID, comm, errCh)
					totalBytesSent += bytesSent
				}
				for _, partyID := range otherPartyIDs {
					recv_msg, bytesRead := recvTSSMessage(party, *partyID, comm, true)
					totalBytesRead += bytesRead
					if recv_msg.Type() == msg.Type() {
						go party.Update(recv_msg)
					} else {
						log.Fatalf("Message received has type %s whereas message sent has type %s", recv_msg.Type(), msg.Type())
					}
				}
			} else { // point-to-point!
				if dest[0].Index == msg.GetFrom().Index {
					log.Fatalf("party %d tried to send a message to itself (%d)", dest[0].Index, msg.GetFrom().Index)
				}
				bytesSent := sendTSSMessage(msg, party, *dest[0], comm, errCh)
				totalBytesSent += bytesSent
				recv_msg, bytesRead := recvTSSMessage(party, *dest[0], comm, false)
				totalBytesRead += bytesRead
				if recv_msg.Type() == msg.Type() {
					go party.Update(recv_msg)
				} else {
					log.Fatalf("Message received has type %s whereas message sent has type %s", recv_msg.Type(), msg.Type())
				}
			}
		case save := <-endCh:
			endTime = time.Now()
			response := map[string]interface{}{
				"signature":          fmt.Sprintf("%x", save.Signature),
				"signing_bytes_read": totalBytesRead,
				"signing_bytes_sent": totalBytesSent,
				"signing_time":       endTime.Sub(startTime).String(),
			}

			responseJson, err := json.Marshal(response)
			if err != nil {
				log.Fatalf("Failed to convert JSON: %v", err)
			}

			return string(responseJson)
		}
	}
}
