package loadbalancer_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	loadbalancer "github.com/loureirovinicius/cleanup/aws/service/ec2/loadBalancer"
	"github.com/loureirovinicius/cleanup/helpers/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEC2 is a mock of EC2 interface
type MockEC2 struct {
	mock.Mock
}

func init() {
	logger.InitializeLogger("info", "json")
}

func (m *MockEC2) DescribeLoadBalancers(ctx context.Context, params *elasticloadbalancingv2.DescribeLoadBalancersInput, optFns ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DescribeLoadBalancersOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*elasticloadbalancingv2.DescribeLoadBalancersOutput), args.Error(1)
}

func (m *MockEC2) DescribeListeners(ctx context.Context, params *elasticloadbalancingv2.DescribeListenersInput, optFns ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DescribeListenersOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*elasticloadbalancingv2.DescribeListenersOutput), args.Error(1)
}

func (m *MockEC2) DeleteLoadBalancer(ctx context.Context, params *elasticloadbalancingv2.DeleteLoadBalancerInput, optFns ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DeleteLoadBalancerOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*elasticloadbalancingv2.DeleteLoadBalancerOutput), args.Error(1)
}

func TestList(t *testing.T) {
	mockSvc := new(MockEC2)

	mockOutput := &elasticloadbalancingv2.DescribeLoadBalancersOutput{
		LoadBalancers: []types.LoadBalancer{
			{
				LoadBalancerArn: aws.String("arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-load-balancer/12ab3c456d7e8900"),
			},
			{
				LoadBalancerArn: aws.String("arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-load-balancer/12ab3c456d7e8901"),
			},
		},
	}

	// Mock AWS client response
	mockSvc.On("DescribeLoadBalancers", mock.Anything, mock.Anything).Return(mockOutput, nil)

	// Instantiate the object responsible for calling the methods
	lb := &loadbalancer.LoadBalancer{
		API: mockSvc,
	}

	// Call the "List" function
	result, err := lb.List(context.Background())

	// Assert the results
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-load-balancer/12ab3c456d7e8900", result[0])
	assert.Equal(t, "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-load-balancer/12ab3c456d7e8901", result[1])

	// Assert that the mock expectations were met
	mockSvc.AssertExpectations(t)
}

func TestValidate(t *testing.T) {
	cases := map[string]struct {
		mockElb []types.Listener
		expect  bool
	}{
		"not deletable elb (has a listener)": {
			mockElb: []types.Listener{
				{
					ListenerArn:     aws.String("arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/test-load-balancer/12ab3c456d7e8900/12ab3c456d7e8900"),
					LoadBalancerArn: aws.String("arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-load-balancer/12ab3c456d7e8900"),
				},
			},
			expect: false,
		},
		"deletable elb (no listeners)": {
			mockElb: []types.Listener{},
			expect:  true,
		},
	}

	// Loop through all test cases
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			mockSvc := new(MockEC2)

			// Instantiate the object responsible for calling the methods
			lb := &loadbalancer.LoadBalancer{
				API: mockSvc,
			}

			mockOutput := &elasticloadbalancingv2.DescribeListenersOutput{
				Listeners: test.mockElb,
			}

			// Mock AWS client response
			mockSvc.On("DescribeListeners", mock.Anything, mock.Anything).Return(mockOutput, nil)

			// Call the "Validate" function
			// Here the LB ARN doesn't matter
			result, err := lb.Validate(context.TODO(), "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-load-balancer/12ab3c456d7e8900")

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

	mockOutput := &elasticloadbalancingv2.DeleteLoadBalancerOutput{}

	// Mock AWS client response
	mockSvc.On("DeleteLoadBalancer", mock.Anything, mock.Anything).Return(mockOutput, nil)

	// Instantiate the object responsible for calling the methods
	lb := &loadbalancer.LoadBalancer{
		API: mockSvc,
	}

	// Call the "Delete" function
	err := lb.Delete(context.Background(), "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-load-balancer/12ab3c456d7e8900")

	// Assert the results
	assert.NoError(t, err)

	// Assert that the mock expectations were met
	mockSvc.AssertExpectations(t)
}
