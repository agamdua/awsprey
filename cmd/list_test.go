package cmd

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type FakeRequestSender struct {
	resp *ec2.DescribeInstancesOutput
}

func makeTestData() *ec2.DescribeInstancesOutput {
	testData := ec2.DescribeInstancesOutput{
		Reservations: []ec2.RunInstancesOutput{
			{
				Instances: []ec2.Instance{
					{
						Tags: []ec2.Tag{
							{
								Key:   aws.String("Name"),
								Value: aws.String("web1"),
							},
							{
								Key:   aws.String("service"),
								Value: aws.String("web"),
							},
							{
								Key:   aws.String("environment"),
								Value: aws.String("staging"),
							},
						},
					},
				},
			},
		},
	}
	return &testData
}

func (r FakeRequestSender) SendRequest(req ec2.DescribeInstancesRequest) (*ec2.DescribeInstancesOutput, error) {
	return r.resp, nil
}

func Test_filterByTag(t *testing.T) {
	golden := []string{"web1"}
	args := []string{"web:staging"}
	rs := FakeRequestSender{
		resp: makeTestData(),
	}

	actual := filterByTag(args, rs)

	if !reflect.DeepEqual(actual, golden) {
		t.Fatalf("Actual: %s, expected: %s", actual, golden)
	}
}
