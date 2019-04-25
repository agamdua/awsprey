// Copyright © 2019 Agam Dua
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
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List EC2 instances in the format of <service>:<environment>",
	Long: `
Info:

This command gives you the capability to receive a list of EC2
instance names (as per the Name key in AWS) by specifying a
key:value pair coorespondig to values of AWS tags:

	* service
	* environment


For example:

	$ awsprey list <service value>:<environment value>
	$ awsprey list web:staging
		`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := external.LoadDefaultAWSConfig()

		helpHowTo := "For more details:\n\t awsprey list --help\n"

		if len(args) != 1 {
			fmt.Println(
				"[ERROR] An argument MUST be passed to list in the form of <service name>:<environment name>")
			fmt.Println(helpHowTo)
			os.Exit(1)
		}

		searchString := args[0]

		if !strings.Contains(searchString, ":") {
			fmt.Println(
				"[ERROR] Must use a colon (:) to separate the service name and environment name.")
			fmt.Println(helpHowTo)
			os.Exit(1)
		}

		filterValues := strings.Split(searchString, ":")

		if err != nil {
			panic(err)
		}

		cfg.Region = endpoints.UsEast1RegionID

		ec2svc := ec2.New(cfg)
		tagFilterService := "tag:service"
		tagFilterEnvironment := "tag:environment"

		serviceFilters := []ec2.Filter{
			{
				Name:   &tagFilterService,
				Values: []string{filterValues[0]},
			},
			{
				Name:   &tagFilterEnvironment,
				Values: []string{filterValues[1]},
			},
		}

		input := ec2.DescribeInstancesInput{
			Filters: serviceFilters,
		}

		req := ec2svc.DescribeInstancesRequest(&input)
		resp, err := req.Send()

		if err != nil {
			panic(err)
		}

		// TODO: when the whole function is refactored out from `main` we should
		//return the list of instance names
		var instanceNames []string

		for _, reservation := range resp.Reservations {
			for _, instance := range reservation.Instances {
				for _, tag := range instance.Tags {
					if *tag.Key == "Name" {
						instanceNames = append(instanceNames, *tag.Value)
					}
				}
			}
		}

		sort.Strings(instanceNames)
		for _, instance := range instanceNames {
			fmt.Println(instance)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
