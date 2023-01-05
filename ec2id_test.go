package main

import (
	"context"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func TestEc2id(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := NewMockEC2DescribeInstancesAPI(ctrl)
	mockClient.EXPECT().
		DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
			Filters: []types.Filter{
				{
					Name: aws.String("tag:Name"),
					Values: []string{"test"},
				},
			},
		}).
		Return(&ec2.DescribeInstancesOutput{}, nil).
		AnyTimes()

	cases := []struct {
		name string
		client EC2DescribeInstancesAPI
		expect string
		wantErr bool
		expectErr string
	}{
		{
			name: "return no instance-id",
			client: mockClient,
			expect: "",
			wantErr: false,
			expectErr: "",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			id, err := Ec2id("test", tt.client)
			if tt.wantErr {
				if (!strings.Contains(err.Error(), tt.expectErr)) {
					t.Errorf("expect NoSuchKey error, got %T", err)
				}
				return
			}
			if err != nil {
				t.Errorf("expect no error, got error: %v", err)
			}
			if expected := tt.expect; id != expected {
				t.Errorf("expect no output, got id: %s", id)
			}
		})
	}
}
