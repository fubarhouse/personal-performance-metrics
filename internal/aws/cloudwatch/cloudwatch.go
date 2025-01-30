package cloudwatch

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/stretchr/testify/mock"
)

// MockCloudWatchClient is a mock of the CloudWatch client
type MockCloudWatchClient struct {
	mock.Mock
}

// PutMetricData mocks the PutMetricData method
func (m *MockCloudWatchClient) PutMetricData(ctx context.Context, params *cloudwatch.PutMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.PutMetricDataOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*cloudwatch.PutMetricDataOutput), args.Error(1)
}
