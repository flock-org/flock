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

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/clusterlink-net/clusterlink/cmd/cl-adm/util"

	"github.com/flock-org/flock/relay/config"
)

func createFlockrelayCerts(entity string) {
	fmt.Printf("Creating Party %s certs for Flock relay auth.\n", entity)
	partyDirectory := config.PartyDirectory(entity)
	if err := os.MkdirAll(partyDirectory, 0755); err != nil {
		fmt.Printf("Unable to create directory :%v\n", err)
		return
	}
	err := util.CreateCertificate(&util.CertificateConfig{
		Name:              entity,
		IsClient:          true,
		DNSNames:          []string{entity},
		CAPath:            config.FrCAFile,
		CAKeyPath:         config.FrKeyFile,
		CertOutPath:       filepath.Join(partyDirectory, config.CertificateFileName),
		PrivateKeyOutPath: filepath.Join(partyDirectory, config.PrivateKeyFileName),
	})
	if err != nil {
		fmt.Printf("Unable to generate certficate/key :%v\n", err)
		return
	}
}

func createUserCerts(user, party string) {
	fmt.Printf("Creating Party %s certs for multi-party comms for user %s\n", user, party)
	partyDirectory := config.UserPartyDirectory(user, party)
	userDirectory := config.UserDirectory(user)

	if err := os.MkdirAll(partyDirectory, 0755); err != nil {
		fmt.Printf("Unable to create directory :%v\n", err)
		return
	}
	err := util.CreateCertificate(&util.CertificateConfig{
		Name:              party,
		IsClient:          true,
		IsServer:          true,
		DNSNames:          []string{party},
		CAPath:            filepath.Join(userDirectory, config.UserCAFile),
		CAKeyPath:         filepath.Join(userDirectory, config.UserKeyFile),
		CertOutPath:       filepath.Join(partyDirectory, config.CertificateFileName),
		PrivateKeyOutPath: filepath.Join(partyDirectory, config.PrivateKeyFileName),
	})
	if err != nil {
		fmt.Printf("Unable to generate certficate/key :%v\n", err)
		return
	}
}

func CreateUser(name string) {
	fmt.Printf("Creating CA Cert for user %s.\n", name)
	userDirectory := config.UserDirectory(name)
	if err := os.MkdirAll(userDirectory, 0755); err != nil {
		fmt.Printf("Unable to create directory :%v\n", err)
		return
	}
	err := util.CreateCertificate(&util.CertificateConfig{
		Name:              name,
		IsCA:              true,
		CertOutPath:       filepath.Join(userDirectory, config.UserCAFile),
		PrivateKeyOutPath: filepath.Join(userDirectory, config.UserKeyFile),
	})
	if err != nil {
		fmt.Printf("Unable to generate CA certficate :%v\n", err)
		return
	}
}

func CreateParty(party string, user string) {
	// Create certs using Flock relay's CA cert
	createFlockrelayCerts(party)
	createUserCerts(user, party)
}

func CreateRelay() {
	fmt.Printf("Creating Flock relay CA Cert.\n")
	if err := os.MkdirAll(config.BaseDirectory(), 0755); err != nil {
		fmt.Printf("Unable to create directory :%v\n", err)
		return
	}
	err := util.CreateCertificate(&util.CertificateConfig{
		Name:              config.FlockrelayServerName,
		IsCA:              true,
		CertOutPath:       config.FrCAFile,
		PrivateKeyOutPath: config.FrKeyFile,
	})
	if err != nil {
		fmt.Printf("Unable to generate CA certficate :%v\n", err)
		return
	}
	fmt.Printf("Generating Certs/Key using CA.\n")

	flockrelayDirectory := config.FlockrelayCADirectory()
	if err := os.MkdirAll(flockrelayDirectory, 0755); err != nil {
		fmt.Printf("Unable to create directory :%v\n", err)
		return
	}
	err = util.CreateCertificate(&util.CertificateConfig{
		Name:              config.FlockrelayServerName,
		IsServer:          true,
		IsClient:          true,
		DNSNames:          []string{config.FlockrelayServerName},
		CAPath:            config.FrCAFileRoot,
		CAKeyPath:         config.FrKeyFile,
		CertOutPath:       filepath.Join(flockrelayDirectory, config.CertificateFileName),
		PrivateKeyOutPath: filepath.Join(flockrelayDirectory, config.PrivateKeyFileName),
	})
	if err != nil {
		fmt.Printf("Unable to generate certficate/key :%v\n", err)
		return
	}
}
