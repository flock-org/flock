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

package config

import "path/filepath"

const (
	//FlockrelayServerName is the servername used in flockrelay router
	FlockrelayServerName = "flockrelay"
	//CertsRootDirectory is the directory for storing all certs it can be changed to /etc/ssl/certs
	CertsRootDirectory = "certs"
	//CertsDirectory is the directory for storing all certs it can be changed to /etc/ssl/certs
	CertsDirectory = "certs"
	// FrCAFile is the path to the certificate authority file.
	FrCAFile = CertsRootDirectory + "/flockrelay-ca.pem"
	// FrCAFileRoot is the path to the certificate authority file.
	FrCAFileRoot = CertsDirectory + "/flockrelay-ca.pem"
	// FrKeyFile is the path to the private-key file.
	FrKeyFile = CertsDirectory + "/flockrelay-key.pem"

	// PrivateKeyFileName is the filename used by private key files.
	PrivateKeyFileName = "key.pem"
	// CertificateFileName is the filename used by certificate files.
	CertificateFileName = "cert.pem"

	// UserCAFile is the path to CA cert of the user
	UserCAFile = "user-ca.pem"
	// UserKeyFile is the private key file
	UserKeyFile = "user-key.pem"
)

// BaseDirectory returns the base path of the fabric certificates.
func BaseDirectory() string {
	return CertsRootDirectory
}

// BaseDirectoryRelay returns the base path of the fabric certificates.
func BaseDirectoryRelay() string {
	return CertsDirectory
}

// PartyDirectory returns the base path for a specific party within the relay's domain for auth.
func PartyDirectory(party string) string {
	return filepath.Join(BaseDirectory(), party)
}

// UserDirectory returns the base path for a specific party.
func UserDirectory(user string) string {
	return filepath.Join(BaseDirectory(), user)
}

// UserPartyDirectory returns the base path for a specific party within a user's domain.
func UserPartyDirectory(user string, party string) string {
	return filepath.Join(BaseDirectory(), user, party)
}

// FlockrelayDirectory returns the base path for a relay
func FlockrelayDirectory() string {
	return filepath.Join(BaseDirectory(), FlockrelayServerName)
}

// FlockrelayCADirectory returns the base path for a relay
func FlockrelayCADirectory() string {
	return filepath.Join(BaseDirectoryRelay(), FlockrelayServerName)
}
