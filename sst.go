package sst

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type BucketResources struct {
	BucketName string
}

func Bucket(name string) *BucketResources {
	return fromEnvironment(&BucketResources{}, name)
}

type EventBusResources struct {
	EventBusName string
}

func EventBus(name string) *EventBusResources {
	return fromEnvironment(&EventBusResources{}, name)
}

type FunctionResources struct {
	FunctionName string
}

func Function(name string) *FunctionResources {
	return fromEnvironment(&FunctionResources{}, name)
}

type QueueResources struct {
	QueueUrl string
}

func Queue(name string) *QueueResources {
	return fromEnvironment(&QueueResources{}, name)
}

type TopicResources struct {
	TopicArn string
}

func Topic(name string) *TopicResources {
	return fromEnvironment(&TopicResources{}, name)
}

type RDSResources struct {
	ClusterArn          string
	SecretArn           string
	DefaultDatabaseName string
}

func RDS(name string) *RDSResources {
	return fromEnvironment(&RDSResources{}, name)
}

type TableResources struct {
	TableName string
}

func Table(name string) *TableResources {
	return fromEnvironment(&TableResources{}, name)
}

func fromEnvironment[T any](dest *T, name string) *T {
	v := reflect.ValueOf(dest).Elem()

	constructName := toConstructName(v.Type().Name())

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)

		propName := toPropName(v.Type().Field(i).Name)
		envVar := fmt.Sprintf("SST_%s_%s_%s", constructName, propName, normaliseId(name))

		value := os.Getenv(envVar)
		if value == "" {
			return nil
		}

		f.SetString(value)
	}

	return dest
}

func toConstructName(s string) string {
	return strings.TrimSuffix(s, "Resources")
}

func toPropName(s string) string {
	return strings.ToLower(s[0:1]) + s[1:]
}

func normaliseId(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}
