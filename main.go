package main

import (
	"fmt"

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

	dryRun := false
	tagFilterService := "tag:service"
	tagFilterEnvironment := "tag:environment"

	serviceFilters := []ec2.Filter{
		{
			Name:   &tagFilterService,
			Values: []string{""},
		},
		{
			Name:   &tagFilterEnvironment,
			Values: []string{""},
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

	for _, instance := range resp.Reservations[0].Instances {
		for _, tag := range instance.Tags {
			if *tag.Key == "Name" {
				instanceNames = append(instanceNames, *tag.Value)
				fmt.Println(*tag.Value)
			}
		}
	}
}
