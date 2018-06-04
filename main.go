package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// Version of the CLI
const CLIVersion = "0.2.1"

func main() {
	cfg, err := external.LoadDefaultAWSConfig()

	if err != nil {
		panic(err)
	}

	cfg.Region = endpoints.UsEast1RegionID

	ec2svc := ec2.New(cfg)

	var versionFlag = flag.Bool("version", false, "Displays current version.")

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

	flag.Parse()

	if *versionFlag {
		fmt.Println(CLIVersion)
		os.Exit(0)
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
		filterRunningState := "instance-state-code"

		serviceFilters := []ec2.Filter{
			{
				Name:   &tagFilterService,
				Values: []string{filterValues[0]},
			},
			{
				Name:   &tagFilterEnvironment,
				Values: []string{filterValues[1]},
			},
			{
				Name:   &filterRunningState,
				Values: []string{"16"}, // 16: "running" state
			},
		}

		if len(os.Args) > 3 && os.Args[3] == "with" {
			switch os.Args[4] {
			case "tag":
				tag := os.Args[5]
				extraFilters := strings.Split(tag, ":")
				key := "tag:" + extraFilters[0]
				extraServiceFilters := ec2.Filter{
					Name:   &key,
					Values: []string{extraFilters[1]},
				}
				serviceFilters = append(serviceFilters, extraServiceFilters)

			default:
				errorMessage := fmt.Sprintf(
					"%s is not a supported input. Choose: 'tag'.",
					os.Args[4],
				)
				fmt.Println(errorMessage)
				os.Exit(1)
			}
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
					}
				}
			}
		}

		sort.Strings(instanceNames)
		for _, instance := range instanceNames {
			fmt.Println(instance)
		}
	}
}
