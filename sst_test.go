package sst_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/stretchr/testify/assert"
	"github.com/teamkeel/sst-go"
)

func TestFunctionDoesNotExist(t *testing.T) {
	_, err := sst.Function(context.Background(), "doesntexist")
	assert.Error(t, err)
}

func TestBucketDoesNotExist(t *testing.T) {
	_, err := sst.Bucket(context.Background(), "doesntexist")
	assert.Error(t, err)
}

func TestQueueDoesNotExist(t *testing.T) {
	_, err := sst.Queue(context.Background(), "doesntexist")
	assert.Error(t, err)
}

func TestTopicDoesNotExist(t *testing.T) {
	_, err := sst.Topic(context.Background(), "doesntexist")
	assert.Error(t, err)
}

func TestEventBusDoesNotExist(t *testing.T) {
	_, err := sst.EventBus(context.Background(), "doesntexist")
	assert.Error(t, err)
}

func TestRDSBusDoesNotExist(t *testing.T) {
	_, err := sst.RDS(context.Background(), "doesntexist")
	assert.Error(t, err)
}

func TestFunction(t *testing.T) {
	t.Setenv("SST_Function_functionName_MyFunction", "the-function-name")

	r, err := sst.Function(context.Background(), "MyFunction")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, r.FunctionName, "the-function-name")
}

func TestNormaliseName(t *testing.T) {
	t.Setenv("SST_Function_functionName_my_function", "the-function-name")

	r, err := sst.Function(context.Background(), "my-function")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, r.FunctionName, "the-function-name")
}

func TestBucket(t *testing.T) {
	t.Setenv("SST_Bucket_bucketName_MyBucket", "the-bucket-name")

	r, err := sst.Bucket(context.Background(), "MyBucket")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, r.BucketName, "the-bucket-name")
}

func TestQueue(t *testing.T) {
	t.Setenv("SST_Queue_queueUrl_MyQueue", "the-queue-url")

	r, err := sst.Queue(context.Background(), "MyQueue")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, r.QueueUrl, "the-queue-url")
}

func TestTopic(t *testing.T) {
	t.Setenv("SST_Topic_topicArn_MyTopic", "the-topic-arn")

	r, err := sst.Topic(context.Background(), "MyTopic")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, r.TopicArn, "the-topic-arn")
}

func TestEventBus(t *testing.T) {
	t.Setenv("SST_EventBus_eventBusName_MyEventBus", "the-event-bus-name")

	r, err := sst.EventBus(context.Background(), "MyEventBus")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, r.EventBusName, "the-event-bus-name")
}

func TestRDS(t *testing.T) {
	t.Setenv("SST_RDS_clusterArn_MyDatabase", "the-cluster-arn")
	t.Setenv("SST_RDS_secretArn_MyDatabase", "the-secret-arn")
	t.Setenv("SST_RDS_defaultDatabaseName_MyDatabase", "the-database-name")

	r, err := sst.RDS(context.Background(), "MyDatabase")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, r.ClusterArn, "the-cluster-arn")
	assert.Equal(t, r.SecretArn, "the-secret-arn")
	assert.Equal(t, r.DefaultDatabaseName, "the-database-name")
}

func TestTable(t *testing.T) {
	t.Setenv("SST_Table_tableName_MyTable", "the-table-name")

	r, err := sst.Table(context.Background(), "MyTable")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, r.TableName, "the-table-name")
}

func TestSecretDoesNotExist(t *testing.T) {
	sst.SetSSMClient(nil)
	_, err := sst.Secret(context.Background(), "MY_SECRET")
	assert.Error(t, err)
}

type mockSSM struct {
	responseIndex int
	responses     []*ssm.GetParametersOutput
}

func (m *mockSSM) GetParameters(ctx context.Context, r *ssm.GetParametersInput, opts ...func(*ssm.Options)) (*ssm.GetParametersOutput, error) {
	resp := m.responses[m.responseIndex]
	m.responseIndex++
	return resp, nil
}

func TestSecret(t *testing.T) {
	t.Setenv("SST_Secret_value_MY_SECRET", "__FETCH_FROM_SSM__")
	t.Setenv("SST_SSM_PREFIX", "/sst/my-app/prod/")

	ssmClient := &mockSSM{
		responses: []*ssm.GetParametersOutput{
			{
				Parameters: []types.Parameter{
					{
						Name:  aws.String("/sst/my-app/prod/Secret/MY_SECRET/value"),
						Value: aws.String("my secret value"),
					},
				},
			},
		},
	}
	sst.SetSSMClient(ssmClient)

	value, err := sst.Secret(context.Background(), "MY_SECRET")
	assert.NoError(t, err)

	assert.Equal(t, "my secret value", value)
}

func TestSecretFallback(t *testing.T) {
	t.Setenv("SST_Secret_value_MY_SECRET", "__FETCH_FROM_SSM__")
	t.Setenv("SST_SSM_PREFIX", "/sst/my-app/prod/")
	t.Setenv("SST_APP", "my-app")

	ssmClient := &mockSSM{
		responses: []*ssm.GetParametersOutput{
			{
				InvalidParameters: []string{"/sst/my-app/prod/Secret/MY_SECRET/value"},
			},
			{
				Parameters: []types.Parameter{
					{
						Name:  aws.String("/sst/my-app/.fallback/Secret/MY_SECRET/value"),
						Value: aws.String("my fallback secret value"),
					},
				},
			},
		},
	}
	sst.SetSSMClient(ssmClient)

	value, err := sst.Secret(context.Background(), "MY_SECRET")
	assert.NoError(t, err)

	assert.Equal(t, "my fallback secret value", value)
}

func TestSecretMissing(t *testing.T) {
	t.Setenv("SST_Secret_value_MY_SECRET", "__FETCH_FROM_SSM__")
	t.Setenv("SST_Secret_value_MY_OTHER_SECRET", "__FETCH_FROM_SSM__")
	t.Setenv("SST_SSM_PREFIX", "/sst/my-app/prod/")
	t.Setenv("SST_APP", "my-app")
	t.Setenv("SST_STAGE", "prod")

	ssmClient := &mockSSM{
		responses: []*ssm.GetParametersOutput{
			{
				InvalidParameters: []string{
					"/sst/my-app/prod/Secret/MY_SECRET/value",
					"/sst/my-app/prod/Secret/MY_OTHER_SECRET/value",
				},
			},
			{
				InvalidParameters: []string{
					"/sst/my-app/prod/Secret/MY_SECRET/value",
					"/sst/my-app/prod/Secret/MY_OTHER_SECRET/value",
				},
			},
		},
	}
	sst.SetSSMClient(ssmClient)

	_, err := sst.Secret(context.Background(), "MY_SECRET")
	assert.Error(t, err)
	assert.Equal(t, "the following secrets are not set in the prod stage: MY_SECRET, MY_OTHER_SECRET", err.Error())
}
