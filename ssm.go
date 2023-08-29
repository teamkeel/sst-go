package sst

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/samber/lo"
)

const (
	fetchFromSsmPlaceholder = "__FETCH_FROM_SSM__"
	fetchFromSecretPrefix   = "__FETCH_FROM_SECRET__:"
)

var (
	mutex     sync.Mutex
	ssmValues map[string]string
	ssmClient SSMClient
)

type SSMClient interface {
	GetParameters(context.Context, *ssm.GetParametersInput, ...func(*ssm.Options)) (*ssm.GetParametersOutput, error)
}

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}

	ssmClient = ssm.NewFromConfig(cfg)
}

func SetSSMClient(c SSMClient) {
	mutex.Lock()
	defer mutex.Unlock()
	ssmValues = nil
	ssmClient = c
}

func fetchValuesFromSSM(ctx context.Context) (map[string]string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if ssmValues != nil {
		return ssmValues, nil
	}

	secretVars := []*envVar{}

	for _, kv := range os.Environ() {
		parts := strings.SplitN(kv, "=", 2)
		key, value := parts[0], parts[1]

		switch {
		case value == fetchFromSsmPlaceholder:
			v := parseEnvName(key)
			if v.constructName == "Secret" {
				secretVars = append(secretVars, parseEnvName(key))
			}
		}
	}

	// No params - we're done
	if len(secretVars) == 0 {
		ssmValues = map[string]string{}
		return ssmValues, nil
	}

	params, invalid, err := loadSecrets(ctx, lo.Map(secretVars, func(v *envVar, _ int) string {
		return buildSsmPath(v)
	}))
	if err != nil {
		return nil, err
	}

	values := map[string]string{}
	for _, p := range params {
		v := parseSsmPath(*p.Name)
		values[v.constructID] = *p.Value
	}

	// For any missing values try fallbacks
	if len(invalid) > 0 {
		params, invalid, err := loadSecrets(ctx, lo.Map(invalid, func(name string, _ int) string {
			v := parseSsmPath(name)
			return buildSsmFallbackPath(v)
		}))
		if err != nil {
			return nil, err
		}

		// If we still have missing params that's an error
		if len(invalid) > 0 {
			missing := lo.Map(invalid, func(s string, _ int) string {
				v := parseSsmFallbackPath(s)
				return v.constructID
			})
			return nil, fmt.Errorf("the following secrets are not set in the %s stage: %s", os.Getenv("SST_STAGE"), strings.Join(missing, ", "))
		}

		values = map[string]string{}
		for _, p := range params {
			v := parseSsmFallbackPath(*p.Name)
			values[v.constructID] = *p.Value
		}
	}

	ssmValues = values
	return ssmValues, nil
}

func loadSecrets(ctx context.Context, names []string) ([]types.Parameter, []string, error) {
	chunks := [][]string{}
	for i := 0; i < len(names); i += 10 {
		chunk := names[i:int(math.Min(float64(i+10), float64(len(names))))]
		chunks = append(chunks, chunk)
	}

	params := []types.Parameter{}
	invalid := []string{}

	for _, chunk := range chunks {
		out, err := ssmClient.GetParameters(ctx, &ssm.GetParametersInput{
			Names:          chunk,
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			return nil, nil, err
		}

		params = append(params, out.Parameters...)
		invalid = append(invalid, out.InvalidParameters...)
	}

	return params, invalid, nil
}

type envVar struct {
	constructName string
	constructID   string
	propName      string
}

func parseEnvName(env string) *envVar {
	parts := strings.SplitN(env, "_", 4)
	return &envVar{
		constructName: parts[1],
		propName:      parts[2],
		constructID:   parts[3],
	}
}

func parseSsmPath(path string) *envVar {
	prefix := ssmPrefix()
	path = strings.TrimPrefix(path, prefix)
	parts := strings.Split(path, "/")
	return &envVar{
		constructName: parts[0],
		constructID:   parts[1],
		propName:      parts[2],
	}
}

func parseSsmFallbackPath(path string) *envVar {
	parts := strings.Split(path, "/")
	return &envVar{
		constructName: parts[4],
		constructID:   parts[5],
		propName:      parts[6],
	}
}

func buildSsmPath(data *envVar) string {
	return fmt.Sprintf(`%s%s/%s/%s`, ssmPrefix(), data.constructName, data.constructID, data.propName)
}

func buildSsmFallbackPath(data *envVar) string {
	return fmt.Sprintf(
		`/sst/%s/.fallback/%s/%s/%s`,
		os.Getenv("SST_APP"), data.constructName, data.constructID, data.propName,
	)
}

func ssmPrefix() string {
	return os.Getenv("SST_SSM_PREFIX")
}
