package main

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/golang/mock/gomock"
)

func TestEc2id(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := NewMockEC2DescribeInstancesAPI(ctrl)
	// reservation-id, instance-id, launch-time, tag:Name
	// r-03803cf1200d61cf7,i-0123456789abcdef0, 2023-01-05-12:00:00.00, sample
	// r-03803cf1200d61cf7,i-0123456789abcdef1, 2023-01-05-12:00:00.01, sample
	// r-03803cf1200d61cf7,i-0123456789abcdef2, 2023-01-05-12:00:00.02, sample2
	// r-03803cf1200d61cf7,i-00000000000abcdef, 2023-01-04-12:00:00.00, test
	mockClient.EXPECT().
		DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
			Filters: []types.Filter{
				{
					Name:   aws.String("instance-state-name"),
					Values: []string{"running"},
				},
				{
					Name:   aws.String("tag:Name"),
					Values: []string{"noexist"},
				},
			},
		}).
		Return(&ec2.DescribeInstancesOutput{}, nil).
		AnyTimes()

	mockClient.EXPECT().
		DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
			Filters: []types.Filter{
				{
					Name:   aws.String("instance-state-name"),
					Values: []string{"running"},
				},
			},
		}).
		Return(&ec2.DescribeInstancesOutput{
			Reservations: []types.Reservation{
				{
					Instances: []types.Instance{
						{
							InstanceId: aws.String("i-0123456789abcdef0"),
							LaunchTime: aws.Time(time.Date(2023, 1, 5, 12, 0, 0, 0, time.UTC)),
						},
						{
							InstanceId: aws.String("i-0123456789abcdef1"),
							LaunchTime: aws.Time(time.Date(2023, 1, 5, 12, 0, 0, 1, time.UTC)),
						},
						{
							InstanceId: aws.String("i-0123456789abcdef2"),
							LaunchTime: aws.Time(time.Date(2023, 1, 5, 12, 0, 0, 2, time.UTC)),
						},
					},
				},
				{
					Instances: []types.Instance{
						{
							InstanceId: aws.String("i-00000000000abcdef"),
							LaunchTime: aws.Time(time.Date(2023, 1, 4, 12, 0, 0, 0, time.UTC)),
						},
					},
				},
			},
		}, nil).
		AnyTimes()

	mockClient.EXPECT().
		DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
			Filters: []types.Filter{
				{
					Name:   aws.String("instance-state-name"),
					Values: []string{"running"},
				},
				{
					Name:   aws.String("tag:Name"),
					Values: []string{"test"},
				},
			},
		}).
		Return(&ec2.DescribeInstancesOutput{
			Reservations: []types.Reservation{
				{
					Instances: []types.Instance{
						{
							InstanceId: aws.String("i-00000000000abcdef"),
							LaunchTime: aws.Time(time.Date(2023, 1, 4, 12, 0, 0, 0, time.UTC)),
						},
					},
				},
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name         string
		client       EC2DescribeInstancesAPI
		instanceName string
		expect       []string
		wantErr      bool
		expectErr    string
	}{
		{
			name:         "return no instance-id",
			client:       mockClient,
			instanceName: "noexist",
			expect:       []string{},
			wantErr:      false,
			expectErr:    "",
		},
		{
			name:         "return instance-ids sorted by LaunchTime descending with no input",
			client:       mockClient,
			instanceName: "",
			expect:       []string{"i-0123456789abcdef2", "i-0123456789abcdef1", "i-0123456789abcdef0", "i-00000000000abcdef"},
			wantErr:      false,
			expectErr:    "",
		},
		{
			name:         "return instance-id",
			client:       mockClient,
			instanceName: "test",
			expect:       []string{"i-00000000000abcdef"},
			wantErr:      false,
			expectErr:    "",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ids, err := Ec2id(tt.instanceName, tt.client)
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.expectErr) {
					t.Errorf("expect NoSuchKey error, got %T", err)
				}
				return
			}
			if err != nil {
				t.Errorf("expect no error, got error: %v", err)
			}
			if len(tt.expect) != len(ids) {
				t.Errorf("expect %s, got id: %s", tt.expect, ids)
			}
			for k, v := range ids {
				if v != tt.expect[k] {
					t.Errorf("expect %s, got id: %s at index %d", tt.expect, ids, k)
				}
			}
		})
	}
}
