package targetgroup_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	targetgroup "github.com/loureirovinicius/cleanup/aws/service/ec2/targetGroup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEC2 is a mock of EC2 interface
type MockEC2 struct {
	mock.Mock
}

func (m *MockEC2) DescribeTargetGroups(ctx context.Context, params *elasticloadbalancingv2.DescribeTargetGroupsInput, optFns ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DescribeTargetGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*elasticloadbalancingv2.DescribeTargetGroupsOutput), args.Error(1)
}

func (m *MockEC2) DeleteTargetGroup(ctx context.Context, params *elasticloadbalancingv2.DeleteTargetGroupInput, optFns ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DeleteTargetGroupOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*elasticloadbalancingv2.DeleteTargetGroupOutput), args.Error(1)
}

func TestList(t *testing.T) {
	mockSvc := new(MockEC2)

	mockOutput := &elasticloadbalancingv2.DescribeTargetGroupsOutput{
		TargetGroups: []types.TargetGroup{
			{
				TargetGroupArn: aws.String("arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-load-balancer/12ab3c456d7e8900"),
			},
			{
				TargetGroupArn: aws.String("arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-load-balancer/12ab3c456d7e8901"),
			},
		},
	}

	// Mock AWS client response
	mockSvc.On("DescribeTargetGroups", mock.Anything, mock.Anything).Return(mockOutput, nil)

	// Instantiate the object responsible for calling the methods
	tg := &targetgroup.TargetGroup{
		API: mockSvc,
	}

	// Call the "List" function
	result, err := tg.List(context.Background())

	// Assert the results
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-load-balancer/12ab3c456d7e8900", result[0])
	assert.Equal(t, "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-load-balancer/12ab3c456d7e8901", result[1])

	// Assert that the mock expectations were met
	mockSvc.AssertExpectations(t)
}

func TestValidate(t *testing.T) {
	cases := map[string]struct {
		mockTg types.TargetGroup
		expect bool
	}{
		"not deletable targetGroup (has LBs attached to it)": {
			mockTg: types.TargetGroup{
				TargetGroupArn: aws.String("arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-load-balancer/12ab3c456d7e8900"),
				LoadBalancerArns: []string{
					"arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-load-balancer/12ab3c456d7e8900",
				},
			},
			expect: false,
		},
		"deletable eip (not associated)": {
			mockTg: types.TargetGroup{
				TargetGroupArn:   aws.String("arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-load-balancer/12ab3c456d7e8901"),
				LoadBalancerArns: []string{},
			},
			expect: true,
		},
	}

	// Loop through all test cases
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			mockSvc := new(MockEC2)

			// Instantiate the object responsible for calling the methods
			tg := &targetgroup.TargetGroup{
				API: mockSvc,
			}

			mockOutput := &elasticloadbalancingv2.DescribeTargetGroupsOutput{
				TargetGroups: []types.TargetGroup{test.mockTg},
			}

			// Mock AWS client response
			mockSvc.On("DescribeTargetGroups", mock.Anything, mock.Anything).Return(mockOutput, nil)

			// Call the "Validate" function
			result, err := tg.Validate(context.TODO(), *test.mockTg.TargetGroupArn)

			assert.NoError(t, err)
			assert.EqualValues(t, bool(test.expect), result)

			// Assert that the mock expectations were met
			mockSvc.AssertExpectations(t)

			t.Log(name)
		})
	}

}

func TestDelete(t *testing.T) {
	mockSvc := new(MockEC2)

	mockOutput := &elasticloadbalancingv2.DeleteTargetGroupOutput{}

	// Mock AWS client response
	mockSvc.On("DeleteTargetGroup", mock.Anything, mock.Anything).Return(mockOutput, nil)

	// Instantiate the object responsible for calling the methods
	tg := &targetgroup.TargetGroup{
		API: mockSvc,
	}

	// Call the "Delete" function
	err := tg.Delete(context.Background(), "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-load-balancer/12ab3c456d7e8900")

	// Assert the results
	assert.NoError(t, err)

	// Assert that the mock expectations were met
	mockSvc.AssertExpectations(t)
}
