package elasticblockstorage_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	elasticblockstorage "github.com/loureirovinicius/cleanup/aws/service/ec2/elasticBlockStorage"
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

func (m *MockEC2) DescribeVolumes(ctx context.Context, input *ec2.DescribeVolumesInput, opts ...func(*ec2.Options)) (*ec2.DescribeVolumesOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*ec2.DescribeVolumesOutput), args.Error(1)
}

func (m *MockEC2) DeleteVolume(ctx context.Context, input *ec2.DeleteVolumeInput, optFns ...func(*ec2.Options)) (*ec2.DeleteVolumeOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*ec2.DeleteVolumeOutput), args.Error(1)
}

func TestList(t *testing.T) {
	mockSvc := new(MockEC2)

	mockOutput := &ec2.DescribeVolumesOutput{
		Volumes: []types.Volume{
			{
				VolumeId: aws.String("vol-1234567890abcdef0"),
			},
			{
				VolumeId: aws.String("vol-1234567890abcdef1"),
			},
		},
	}

	// Mock AWS client response
	mockSvc.On("DescribeVolumes", mock.Anything, mock.Anything).Return(mockOutput, nil)

	// Instantiate the object responsible for calling the methods
	ebs := &elasticblockstorage.ElasticBlockStorage{
		API: mockSvc,
	}

	// Call the "List" function
	result, err := ebs.List(context.Background())

	// Assert the results
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "vol-1234567890abcdef0", result[0])
	assert.Equal(t, "vol-1234567890abcdef1", result[1])

	// Assert that the mock expectations were met
	mockSvc.AssertExpectations(t)
}

func TestValidate(t *testing.T) {
	// Scenarios to test
	cases := map[string]struct {
		mockVolume types.Volume
		expect     bool
	}{
		"not deletable volume (tagged)": {
			mockVolume: types.Volume{
				VolumeId: aws.String("vol-1234567890abcdef0"),
				State:    types.VolumeStateInUse,
				Tags: []types.Tag{
					{
						Key:   aws.String("cleanup-ignore"),
						Value: aws.String("true"),
					},
				},
			},
			expect: false,
		},

		"not deletable volume (not tagged, but in-use)": {
			mockVolume: types.Volume{
				VolumeId: aws.String("vol-1234567890abcdef1"),
				State:    types.VolumeStateInUse,
			},
			expect: false,
		},

		"deletable volume (not tagged and state equals available)": {
			mockVolume: types.Volume{
				VolumeId: aws.String("vol-1234567890abcdef2"),
				State:    types.VolumeStateAvailable,
			},
			expect: true,
		},
	}

	// Loop through all test cases
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			mockSvc := new(MockEC2)

			// Instantiate the object responsible for calling the methods
			ebs := &elasticblockstorage.ElasticBlockStorage{
				API: mockSvc,
			}

			mockOutput := &ec2.DescribeVolumesOutput{
				Volumes: []types.Volume{test.mockVolume},
			}

			// Mock AWS client response
			mockSvc.On("DescribeVolumes", mock.Anything, mock.Anything).Return(mockOutput, nil)

			// Call the "Validate" function
			result, err := ebs.Validate(context.TODO(), *test.mockVolume.VolumeId)

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

	mockOutput := &ec2.DeleteVolumeOutput{}

	// Mock AWS client response
	mockSvc.On("DeleteVolume", mock.Anything, mock.Anything).Return(mockOutput, nil)

	// Instantiate the object responsible for calling the methods
	ebs := &elasticblockstorage.ElasticBlockStorage{
		API: mockSvc,
	}

	// Call the "Delete" function
	err := ebs.Delete(context.Background(), "vol-1234567890abcdef0")

	// Assert the results
	assert.NoError(t, err)

	// Assert that the mock expectations were met
	mockSvc.AssertExpectations(t)
}
