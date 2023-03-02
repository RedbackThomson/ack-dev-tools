// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	ackConfigPath string
)

const (
	appName      = "ack-dev"
	appShortDesc = "ack-dev - manage ACK controllers and repositories"
	appLongDesc  = `ack-dev

A tool to manage ACK controllers and repositories`
)

func init() {
	rootCmd.PersistentFlags().StringVar(&ackConfigPath, "config-file", defaultConfigPath, "ackdev configuration file path")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(ensureCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(versionCmd)
}

var rootCmd = &cobra.Command{
	Use:           "ackdev",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "A tool to manage ACK repositories, CRDs, development tools and testing",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
