package elasticnetworkinterface_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	elasticnetworkinterface "github.com/loureirovinicius/cleanup/aws/service/ec2/elasticNetworkInterface"
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

func (m *MockEC2) DescribeNetworkInterfaces(ctx context.Context, params *ec2.DescribeNetworkInterfacesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeNetworkInterfacesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeNetworkInterfacesOutput), args.Error(1)
}

func (m *MockEC2) DeleteNetworkInterface(ctx context.Context, params *ec2.DeleteNetworkInterfaceInput, optFns ...func(*ec2.Options)) (*ec2.DeleteNetworkInterfaceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteNetworkInterfaceOutput), args.Error(1)
}

func TestList(t *testing.T) {
	mockSvc := new(MockEC2)

	mockOutput := &ec2.DescribeNetworkInterfacesOutput{
		NetworkInterfaces: []types.NetworkInterface{
			{
				NetworkInterfaceId: aws.String("eni-e1ab23a0"),
			},
			{
				NetworkInterfaceId: aws.String("eni-e1ab23a1"),
			},
		},
	}

	// Mock AWS client response
	mockSvc.On("DescribeNetworkInterfaces", mock.Anything, mock.Anything).Return(mockOutput, nil)

	// Instantiate the object responsible for calling the methods
	eni := &elasticnetworkinterface.ElasticNetworkInterface{
		API: mockSvc,
	}

	// Call the "List" function
	result, err := eni.List(context.Background())

	// Assert the results
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "eni-e1ab23a0", result[0])
	assert.Equal(t, "eni-e1ab23a1", result[1])

	// Assert that the mock expectations were met
	mockSvc.AssertExpectations(t)
}

func TestValidate(t *testing.T) {
	cases := map[string]struct {
		mockEni types.NetworkInterface
		expect  bool
	}{
		"not deletable eni (status is not 'available')": {
			mockEni: types.NetworkInterface{
				NetworkInterfaceId: aws.String("eni-e1ab23a0"),
				Status:             types.NetworkInterfaceStatusAssociated,
			},
			expect: false,
		},
		"deletable eni (status is 'available')": {
			mockEni: types.NetworkInterface{
				NetworkInterfaceId: aws.String("eni-e1ab23a1"),
				Status:             types.NetworkInterfaceStatusAvailable,
			},
			expect: true,
		},
	}

	// Loop through all test cases
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			mockSvc := new(MockEC2)

			// Instantiate the object responsible for calling the methods
			eni := &elasticnetworkinterface.ElasticNetworkInterface{
				API: mockSvc,
			}

			mockOutput := &ec2.DescribeNetworkInterfacesOutput{
				NetworkInterfaces: []types.NetworkInterface{test.mockEni},
			}

			// Mock AWS client response
			mockSvc.On("DescribeNetworkInterfaces", mock.Anything, mock.Anything).Return(mockOutput, nil)

			// Call the "Validate" function
			result, err := eni.Validate(context.TODO(), *test.mockEni.NetworkInterfaceId)

			assert.NoError(t, err)
			assert.EqualValues(t, bool(test.expect), result)

			mockSvc.AssertExpectations(t)

			t.Log(name)
		})
	}

}

func TestDelete(t *testing.T) {
	mockSvc := new(MockEC2)

	mockOutput := &ec2.DeleteNetworkInterfaceOutput{}

	// Mock AWS client response
	mockSvc.On("DeleteNetworkInterface", mock.Anything, mock.Anything).Return(mockOutput, nil)

	// Instantiate the object responsible for calling the methods
	eni := &elasticnetworkinterface.ElasticNetworkInterface{
		API: mockSvc,
	}

	// Call the "Delete" function
	err := eni.Delete(context.Background(), "eni-e1ab23a0")

	// Assert the results
	assert.NoError(t, err)

	// Assert that the mock expectations were met
	mockSvc.AssertExpectations(t)
}
