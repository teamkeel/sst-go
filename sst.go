package sst

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
)

const paramPrefix = "SST_Parameter_value_"

// Secrets returns all Config.Secret bindings as a map.
func Secrets(ctx context.Context) (map[string]string, error) {
	return fetchValuesFromSSM(ctx)
}

// Secret returns a single secret by name.
func Secret(ctx context.Context, name string) (string, error) {
	values, err := fetchValuesFromSSM(ctx)
	if err != nil {
		return "", err
	}

	v, ok := values[name]
	if !ok {
		return "", fmt.Errorf("no secret set with name %s", name)
	}

	return v, nil
}

// Parameters returns all Config.Parameter bindings as a map
func Parameters(ctx context.Context) map[string]string {
	params := map[string]string{}
	for _, kv := range os.Environ() {
		parts := strings.SplitN(kv, "=", 1)
		key, value := parts[0], parts[1]
		if strings.HasPrefix(key, paramPrefix) {
			params[strings.TrimPrefix(key, paramPrefix)] = value
		}
	}
	return params
}

// Parameter returns a single parameter by name
func Parameter(ctx context.Context, name string) (string, error) {
	v, ok := os.LookupEnv(fmt.Sprintf("%s%s", paramPrefix, name))
	if !ok {
		return "", fmt.Errorf("parameter %s is not set", name)
	}
	return v, nil
}

type BucketResources struct {
	BucketName string
}

func Bucket(ctx context.Context, name string) (*BucketResources, error) {
	return fromEnvironment(ctx, &BucketResources{}, name)
}

type EventBusResources struct {
	EventBusName string
}

func EventBus(ctx context.Context, name string) (*EventBusResources, error) {
	return fromEnvironment(ctx, &EventBusResources{}, name)
}

type FunctionResources struct {
	FunctionName string
}

func Function(ctx context.Context, name string) (*FunctionResources, error) {
	return fromEnvironment(ctx, &FunctionResources{}, name)
}

type QueueResources struct {
	QueueUrl string
}

func Queue(ctx context.Context, name string) (*QueueResources, error) {
	return fromEnvironment(ctx, &QueueResources{}, name)
}

type TopicResources struct {
	TopicArn string
}

func Topic(ctx context.Context, name string) (*TopicResources, error) {
	return fromEnvironment(ctx, &TopicResources{}, name)
}

type RDSResources struct {
	ClusterArn          string
	SecretArn           string
	DefaultDatabaseName string
}

func RDS(ctx context.Context, name string) (*RDSResources, error) {
	return fromEnvironment(ctx, &RDSResources{}, name)
}

type TableResources struct {
	TableName string
}

func Table(ctx context.Context, name string) (*TableResources, error) {
	return fromEnvironment(ctx, &TableResources{}, name)
}

func fromEnvironment[T any](ctx context.Context, dest *T, name string) (*T, error) {
	v := reflect.ValueOf(dest).Elem()

	constructName := toConstructName(v.Type().Name())

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)

		propName := toPropName(v.Type().Field(i).Name)
		envVar := fmt.Sprintf("SST_%s_%s_%s", constructName, propName, normaliseId(name))

		value := os.Getenv(envVar)
		if value == "" {
			return nil, fmt.Errorf("required env var %s not set for %s %s", envVar, constructName, normaliseId(name))
		}

		if strings.HasPrefix(value, fetchFromSecretPrefix) {
			var err error
			value, err = Secret(ctx, strings.TrimPrefix(value, fetchFromSecretPrefix))
			if err != nil {
				return nil, err
			}
		}

		f.SetString(value)
	}

	return dest, nil
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
