package cleaner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/loureirovinicius/cleanup/helpers/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type LogOutput struct {
	Time  time.Time `json:"time"`
	Level string    `json:"level"`
	Msg   string    `json:"msg"`
}

type MockCleanable struct {
	mock.Mock
}

func (m *MockCleanable) List(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockCleanable) Validate(ctx context.Context, resource string) (bool, error) {
	args := m.Called(ctx, resource)
	return args.Bool(0), args.Error(1)
}

func (m *MockCleanable) Delete(ctx context.Context, resource string) error {
	args := m.Called(ctx, resource)
	return args.Error(0)
}

// Read output and unmarshall the JSON log into a log struct
func getLastLogLine(logs string) (string, error) {
	log := new(LogOutput)

	logLines := strings.Split(logs, "\n")
	lastLog := logLines[len(logLines)-2]

	if err := json.Unmarshal([]byte(lastLog), &log); err != nil {
		return "", errors.New("Failed to parse log output")
	}

	return log.Msg, nil
}

func TestList(t *testing.T) {
	var buf bytes.Buffer
	ctx := context.Background()
	mockService := new(MockCleanable)

	// Initialize logger and set it to output to a buffer
	logger.InitializeLogger("info", "json", &buf)

	// Test cases
	cases := map[string]struct {
		input    string
		helpers  func()
		testCase func(*testing.T, string, error)
	}{
		"Successful list all resources": {
			input: "TestService",
			helpers: func() {
				mockService.On("List", ctx).Return([]string{"res1", "res2"}, nil)
			},
			testCase: func(t *testing.T, output string, err error) {
				log, testErr := getLastLogLine(output)
				if testErr != nil {
					require.NoError(t, testErr)
				}

				assert.Equal(t, "Resources for TestService: res1, res2", log)
				assert.Nil(t, err)
			},
		},
		"List returns an error": {
			input: "TestService",
			helpers: func() {
				mockService.On("List", ctx).Return([]string{}, errors.New("list error"))
			},
			testCase: func(t *testing.T, output string, err error) {
				assert.EqualError(t, err, "error listing resources for service 'TestService': list error")
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			mockService.ExpectedCalls = nil

			test.helpers()

			err := list(ctx, mockService, test.input)
			output := buf.String()

			test.testCase(t, output, err)

			buf.Reset()
			mockService.AssertExpectations(t)
		})
	}
}

func TestValidate(t *testing.T) {
	var buf bytes.Buffer
	ctx := context.Background()
	mockService := new(MockCleanable)

	// Initialize logger and set it to output to a buffer
	logger.InitializeLogger("info", "json", &buf)

	// Test cases
	cases := map[string]struct {
		input    string
		helpers  func()
		testCase func(*testing.T, string, error)
	}{
		"Successful validation of all resources": {
			input: "TestService",
			helpers: func() {
				mockService.On("List", ctx).Return([]string{"res1"}, nil)
				mockService.On("Validate", ctx, "res1").Return(true, nil)
			},
			testCase: func(t *testing.T, output string, err error) {
				assert.Nil(t, err)
			},
		},
		"Validation returns an error": {
			input: "TestService",
			helpers: func() {
				mockService.On("List", ctx).Return([]string{"res1"}, nil)
				mockService.On("Validate", ctx, "res1").Return(false, errors.New("validation error"))
			},
			testCase: func(t *testing.T, output string, err error) {
				assert.EqualError(t, err, "error validating resource 'res1' in service 'TestService': validation error")
			},
		},
		"Resource is deletable": {
			input: "TestService",
			helpers: func() {
				mockService.On("List", ctx).Return([]string{"res1"}, nil)
				mockService.On("Validate", ctx, "res1").Return(true, nil)
			},
			testCase: func(t *testing.T, output string, err error) {
				log, testErr := getLastLogLine(output)
				if testErr != nil {
					require.NoError(t, testErr)
				}

				assert.Equal(t, "Resource 'res1' in service 'TestService' is empty and can be excluded.", log)
				assert.Nil(t, err)
			},
		},
		"Resource is not deletable": {
			input: "TestService",
			helpers: func() {
				mockService.On("List", ctx).Return([]string{"res1"}, nil)
				mockService.On("Validate", ctx, "res1").Return(false, nil)
			},
			testCase: func(t *testing.T, output string, err error) {
				log, testErr := getLastLogLine(output)
				if testErr != nil {
					require.NoError(t, testErr)
				}

				assert.Equal(t, "Resource 'res1' in service 'TestService' is not empty and cannot be excluded.", log)
				assert.Nil(t, err)
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			mockService.ExpectedCalls = nil

			test.helpers()

			err := validate(ctx, mockService, test.input)
			output := buf.String()

			test.testCase(t, output, err)

			buf.Reset()
			mockService.AssertExpectations(t)
		})
	}
}

func TestDelete(t *testing.T) {
	var buf bytes.Buffer
	ctx := context.Background()
	mockService := new(MockCleanable)

	// Initialize logger and set it to output to a buffer
	logger.InitializeLogger("info", "json", &buf)

	// Test cases
	cases := map[string]struct {
		input    string
		helpers  func()
		testCase func(*testing.T, string, error)
	}{
		"Successful deletion of all resources": {
			input: "TestService",
			helpers: func() {
				mockService.On("List", ctx).Return([]string{"res1"}, nil)
				mockService.On("Validate", ctx, "res1").Return(true, nil)
				mockService.On("Delete", ctx, "res1").Return(nil)
			},
			testCase: func(t *testing.T, output string, err error) {
				log, testErr := getLastLogLine(output)
				if testErr != nil {
					require.NoError(t, testErr)
				}

				assert.Equal(t, "Resource 'res1' in service 'TestService' has been deleted successfully.", log)
				assert.Nil(t, err)
			},
		},
		"Deletion fails for a resource": {
			input: "TestService",
			helpers: func() {
				mockService.On("List", ctx).Return([]string{"res1"}, nil)
				mockService.On("Validate", ctx, "res1").Return(true, nil)
				mockService.On("Delete", ctx, "res1").Return(errors.New("delete error"))
			},
			testCase: func(t *testing.T, output string, err error) {
				assert.EqualError(t, err, "error deleting resource 'res1' in service 'TestService': delete error")
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			mockService.ExpectedCalls = nil

			test.helpers()

			err := delete(ctx, mockService, test.input)
			output := buf.String()

			test.testCase(t, output, err)

			buf.Reset()
			mockService.AssertExpectations(t)
		})
	}
}
