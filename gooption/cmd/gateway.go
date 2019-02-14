// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
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

package cmd

import (
	"os"

	"github.com/gooption-io/gooption/proto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	module            string
	tcpPort, httpPort string
)

// gatewayCmd represents the gateway command
var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "Starts http reverse proxy",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := gooption.ServeEuropeanOptionPricerServerGateway(tcpPort, httpPort); err != nil {
			logrus.Errorln(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(gatewayCmd)

	gatewayCmd.Flags().StringVar(&httpPort, "address", ":8081", "proxy url")
	gatewayCmd.Flags().StringVar(&tcpPort, "grpc-address", "0.0.0.0:50051", "grpc server address")
}
