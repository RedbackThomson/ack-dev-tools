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

import "github.com/spf13/cobra"

var (
	optSoakService string
)

func init() {
	soakCmd.PersistentFlags().StringVar(&optSoakService, "service", "", "service going to be running tests on the cluster")

	soakCmd.AddCommand(soakBootstrapCmd)
	soakCmd.AddCommand(soakInstallCmd)
}

var soakCmd = &cobra.Command{
	Use:   "soak",
	Args:  cobra.NoArgs,
	Short: "Commands related to the soak tests",
}
