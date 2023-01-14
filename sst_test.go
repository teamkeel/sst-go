package sst_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teamkeel/sst-go"
)

func TestFunctionNotPresent(t *testing.T) {
	assert.Nil(t, sst.Function("doesntexist"))
}

func TestBucketNotPresent(t *testing.T) {
	assert.Nil(t, sst.Bucket("doesntexist"))
}

func TestQueueNotPresent(t *testing.T) {
	assert.Nil(t, sst.Queue("doesntexist"))
}

func TestTopicNotPresent(t *testing.T) {
	assert.Nil(t, sst.Topic("doesntexist"))
}

func TestEventBusNotPresent(t *testing.T) {
	assert.Nil(t, sst.EventBus("doesntexist"))
}

func TestRDSBusNotPresent(t *testing.T) {
	assert.Nil(t, sst.RDS("doesntexist"))
}

func TestFunction(t *testing.T) {
	t.Setenv("SST_Function_functionName_MyFunction", "the-function-name")

	r := sst.Function("MyFunction")
	assert.NotNil(t, r)
	assert.Equal(t, r.FunctionName, "the-function-name")
}

func TestBucket(t *testing.T) {
	t.Setenv("SST_Bucket_bucketName_MyBucket", "the-bucket-name")

	r := sst.Bucket("MyBucket")
	assert.NotNil(t, r)
	assert.Equal(t, r.BucketName, "the-bucket-name")
}

func TestQueue(t *testing.T) {
	t.Setenv("SST_Queue_queueUrl_MyQueue", "the-queue-url")

	r := sst.Queue("MyQueue")
	assert.NotNil(t, r)
	assert.Equal(t, r.QueueUrl, "the-queue-url")
}

func TestTopic(t *testing.T) {
	t.Setenv("SST_Topic_topicArn_MyTopic", "the-topic-arn")

	r := sst.Topic("MyTopic")
	assert.NotNil(t, r)
	assert.Equal(t, r.TopicArn, "the-topic-arn")
}

func TestEventBus(t *testing.T) {
	t.Setenv("SST_EventBus_eventBusName_MyEventBus", "the-event-bus-name")

	r := sst.EventBus("MyEventBus")
	assert.NotNil(t, r)
	assert.Equal(t, r.EventBusName, "the-event-bus-name")
}

func TestRDS(t *testing.T) {
	t.Setenv("SST_RDS_clusterArn_MyDatabase", "the-cluster-arn")
	t.Setenv("SST_RDS_secretArn_MyDatabase", "the-secret-arn")
	t.Setenv("SST_RDS_defaultDatabaseName_MyDatabase", "the-database-name")

	r := sst.RDS("MyDatabase")
	assert.NotNil(t, r)
	assert.Equal(t, r.ClusterArn, "the-cluster-arn")
	assert.Equal(t, r.SecretArn, "the-secret-arn")
	assert.Equal(t, r.DefaultDatabaseName, "the-database-name")
}

func TestTable(t *testing.T) {
	t.Setenv("SST_Table_tableName_MyTable", "the-table-name")

	r := sst.Table("MyTable")
	assert.NotNil(t, r)
	assert.Equal(t, r.TableName, "the-table-name")
}
