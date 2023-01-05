//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=mock_$GOFILE
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2DescribeInstancesAPI interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

func GetInstances(c context.Context, api EC2DescribeInstancesAPI, input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return api.DescribeInstances(c, input)
}

func NewAwsClient() (EC2DescribeInstancesAPI, error){
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Fprintln(os.Stderr, "configuration error")
		return nil, err
	}
	return ec2.NewFromConfig(cfg), nil
}

func Ec2id(name string, client EC2DescribeInstancesAPI) (string, error) {
	var params = &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []string{name},
			},
		}
	}

	result, err := GetInstances(context.TODO(), client, params)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Got an error retrieving information about your Amazon EC2 instance")
		return "", err
	}

	if len(result.Reservations) == 0 {
		return "", nil
	}

	// TODO jmespath.Searchで書き換えるとシンプルになる？
	var filteredInstance = result.Reservations[0].Instances[0]
	for _, v := range result.Reservations {
		for _, instance := range v.Instances {
			if string(filteredInstance.State.Name) == "running" && string(instance.State.Name) != "running" {
				continue
			}
			if filteredInstance.LaunchTime.After(*instance.LaunchTime) {
				continue
			}

			filteredInstance = instance
		}
	}

	return *filteredInstance.InstanceId, nil
}
