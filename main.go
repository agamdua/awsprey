package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func main() {
	cfg, err := external.LoadDefaultAWSConfig()

	if err != nil {
		panic(err)
	}

	cfg.Region = endpoints.UsEast1RegionID

	ec2svc := ec2.New(cfg)

	listCommand := flag.NewFlagSet("list", flag.ExitOnError)

	helpMessage :=
		`Oh no - you forgot to specify what exactly to do!
			
Must provide a subcommand out of the following options: list

Example usage:
	$ awsprey list web:staging
`

	if len(os.Args) <= 1 {
		fmt.Println(helpMessage)

		os.Exit(1)
	}

	switch os.Args[1] {
	case "list":
		listCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if listCommand.Parsed() {
		if len(os.Args) <= 2 {
			// shouldn't have to do this, error handling should be at parse time
			fmt.Println(helpMessage)
		}
		listArg := os.Args[2]

		filterValues := strings.Split(listArg, ":")

		dryRun := false

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
			DryRun:  &dryRun,
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
						fmt.Println(*tag.Value)
					}
				}
			}
		}
	}
}
