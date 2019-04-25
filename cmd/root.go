// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var showVersion bool

// rootCmd represents the root command
var rootCmd = &cobra.Command{
	Use:   "awsprey",
	Short: "AWSprey is a way to inspect your AWS infrastructure",
	Long:  ``,
	Run:   run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Display the current version")
}

func run(ccmd *cobra.Command, args []string) {
	if showVersion {
		fmt.Println(CurrentVersion)
	} else {
		ccmd.HelpFunc()(ccmd, args)
	}
}
