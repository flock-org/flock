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

package relay

import (
	"fmt"
	"path/filepath"

	"github.com/clusterlink-net/clusterlink/pkg/util"
	"github.com/flock-org/flock/relay/config"
	relay "github.com/flock-org/flock/relay/pkg/core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rel relay.Relay

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start command starts the flock relay on the specified port",
	Long:  `Start command starts the flow relay on the specific port`,
	Run: func(cmd *cobra.Command, args []string) {
		ip, _ := cmd.Flags().GetString("ip")
		port, _ := cmd.Flags().GetString("port")
		debug, _ := cmd.Flags().GetBool("debug")
		ll := logrus.InfoLevel
		if debug == true {
			ll = logrus.DebugLevel
		}
		rel.Init(ip, port, ll)

		relayDirectory := config.FlockrelayCADirectory()

		// parse TLS files
		parsedCertData, err := util.ParseTLSFiles(config.FrCAFileRoot,
			filepath.Join(relayDirectory, config.CertificateFileName),
			filepath.Join(relayDirectory, config.PrivateKeyFileName))
		if err != nil {
			fmt.Printf("Unable to parse TLS files: %v", err)
			return
		}

		rel.StartRelay(parsedCertData, port)

		// TODO: Start API Server which integrates with the application provider to hand out certificates

		// apiServer := server.NewAPIServer(parsedCertData)
		// apiServer.StartFlockAPIServer()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().String("ip", "", "Optional IP address to bind the flock relay")
	startCmd.Flags().String("port", "9000", "Port to bind the flock relay (default:9000)")
	startCmd.Flags().Bool("debug", false, "Debug mode with verbose prints")
}
