package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2DescribeInstancesAPI interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

func GetInstances(c context.Context, api EC2DescribeInstancesAPI, input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return api.DescribeInstances(c, input)
}

func Ec2id(name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Fprintln(os.Stderr, "configuration error")
		return err
	}
	client := ec2.NewFromConfig(cfg)
	result, err := GetInstances(context.TODO(), client, &ec2.DescribeInstancesInput{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Got an error retrieving information about your Amazon EC2 instance")
		return err
	}

	if len(result.Reservations) == 0 {
		fmt.Fprintf(os.Stderr, "error: reservations length=%d\n", len(result.Reservations))
		return fmt.Errorf("no instance exist")
	}
	if len(result.Reservations[0].Instances) == 0 {
		fmt.Fprintf(os.Stderr, "error: instances length=%d\n", len(result.Reservations[0].Instances))
		return fmt.Errorf("no instance exist")
	}

	fmt.Println(*result.Reservations[0].Instances[0].InstanceId)
	return nil
}
