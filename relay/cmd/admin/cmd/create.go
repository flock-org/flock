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

package admin

import (
	"github.com/flock-org/flock/relay/pkg/api"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a relay/party",
	Long:  `Create a relay/party `,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var createRelayCmd = &cobra.Command{
	Use:   "relay",
	Short: "Create a relay",
	Long:  `Create a relay `,
	Run: func(cmd *cobra.Command, args []string) {
		api.CreateRelay()
	},
}

var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Create a user domain",
	Long:  `Create a user domain `,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		api.CreateUser(name)
	},
}

var createPartyCmd = &cobra.Command{
	Use:   "party",
	Short: "Create a party within user domain",
	Long:  `Create a party within user domain`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		user, _ := cmd.Flags().GetString("user")
		api.CreateParty(name, user)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createRelayCmd)
	createCmd.AddCommand(createUserCmd)
	createUserCmd.Flags().String("name", "", "User name.")
	createCmd.AddCommand(createPartyCmd)
	createPartyCmd.Flags().String("name", "", "Party name.")
	createPartyCmd.Flags().String("user", "", "User name associated.")
}
