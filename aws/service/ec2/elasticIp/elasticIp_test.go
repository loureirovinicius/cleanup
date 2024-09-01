package elasticip_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	elasticip "github.com/loureirovinicius/cleanup/aws/service/ec2/elasticIp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEC2 is a mock of EC2 interface
type MockEC2 struct {
	mock.Mock
}

func (m *MockEC2) DescribeAddresses(ctx context.Context, params *ec2.DescribeAddressesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeAddressesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeAddressesOutput), args.Error(1)
}

func (m *MockEC2) ReleaseAddress(ctx context.Context, params *ec2.ReleaseAddressInput, optFns ...func(*ec2.Options)) (*ec2.ReleaseAddressOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.ReleaseAddressOutput), args.Error(1)
}

func TestList(t *testing.T) {
	mockSvc := new(MockEC2)

	mockOutput := &ec2.DescribeAddressesOutput{
		Addresses: []types.Address{
			{
				AllocationId: aws.String("eipalloc-00a12b30"),
			},
			{
				AllocationId: aws.String("eipalloc-00a12b31"),
			},
		},
	}

	// Mock AWS client response
	mockSvc.On("DescribeAddresses", mock.Anything, mock.Anything).Return(mockOutput, nil)

	// Instantiate the object responsible for calling the methods
	eip := &elasticip.ElasticIP{
		API: mockSvc,
	}

	// Call the "List" function
	result, err := eip.List(context.Background())

	// Assert the results
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "eipalloc-00a12b30", result[0])
	assert.Equal(t, "eipalloc-00a12b31", result[1])

	// Assert that the mock expectations were met
	mockSvc.AssertExpectations(t)
}

func TestValidate(t *testing.T) {
	cases := map[string]struct {
		mockEip types.Address
		expect  bool
	}{
		"not deletable eip (associated)": {
			mockEip: types.Address{
				AllocationId:  aws.String("eipalloc-00a12b30"),
				AssociationId: aws.String("eipassoc-ab12c345"),
			},
			expect: false,
		},
		"deletable eip (not associated)": {
			mockEip: types.Address{
				AllocationId: aws.String("eipalloc-00a12b30"),
			},
			expect: true,
		},
	}

	// Loop through all test cases
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			mockSvc := new(MockEC2)

			eip := &elasticip.ElasticIP{
				API: mockSvc,
			}

			mockOutput := &ec2.DescribeAddressesOutput{
				Addresses: []types.Address{test.mockEip},
			}

			mockSvc.On("DescribeAddresses", mock.Anything, mock.Anything).Return(mockOutput, nil)

			result, err := eip.Validate(context.TODO(), *test.mockEip.AllocationId)

			assert.NoError(t, err)
			assert.EqualValues(t, bool(test.expect), result)

			mockSvc.AssertExpectations(t)

			t.Log(name)
		})
	}

}

func TestDelete(t *testing.T) {
	mockSvc := new(MockEC2)

	mockOutput := &ec2.ReleaseAddressOutput{}

	// Mock AWS client response
	mockSvc.On("ReleaseAddress", mock.Anything, mock.Anything).Return(mockOutput, nil)

	// Instantiate the object responsible for calling the methods
	eip := &elasticip.ElasticIP{
		API: mockSvc,
	}

	// Call the "Delete" function
	err := eip.Delete(context.Background(), "eipalloc-00a12b30")

	// Assert the results
	assert.NoError(t, err)

	// Assert that the mock expectations were met
	mockSvc.AssertExpectations(t)
}
