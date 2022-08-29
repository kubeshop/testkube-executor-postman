package newman

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/executor/secret"
)

func NewEnvFileReader(m map[string]testkube.Variable, paramsFile string, secretEnvs []string) (io.Reader, error) {
	envFile := NewEnvFileFromVariablesMap(m)

	if paramsFile != "" {
		// create env structure from passed params file
		envFromParamsFile, err := NewEnvFileFromString(paramsFile)
		if err != nil {
			return nil, err
		}
		envFile.PrependParams(envFromParamsFile)
	}

	for _, secretEnv := range secretEnvs {
		// create env structure from passed secret
		envFromSecret, err := NewEnvFileFromString(secretEnv)
		if err != nil {
			return nil, err
		}
		envFile.PrependParams(envFromSecret)
	}

	b, err := json.Marshal(envFile)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), err
}

func NewEnvFileFromVariablesMap(m map[string]testkube.Variable) (envFile EnvFile) {
	envFile.ID = "executor-env-file"
	envFile.Name = "executor-env-file"
	envFile.PostmanVariableScope = "environment"
	envFile.PostmanExportedAt = time.Now()
	envFile.PostmanExportedUsing = "Postman/9.15.13"

	var vars []string
	for _, v := range m {
		if v.Type_ == testkube.VariableTypeSecret && v.SecretRef != nil {
			vars = append(vars, v.Name)
		}
	}

	manager := secret.NewEnvManager()
	secretVars := manager.GetVars(vars)
	for _, v := range m {
		value := v.Value
		if secretValue, ok := secretVars[v.Name]; ok {
			value = secretValue
		}

		envFile.Values = append(envFile.Values, Value{Key: v.Name, Value: value, Enabled: true})
	}

	return
}

func NewEnvFileFromString(f string) (envFile EnvFile, err error) {
	err = json.Unmarshal([]byte(f), &envFile)
	return
}

type EnvFile struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	Values               []Value   `json:"values"`
	PostmanVariableScope string    `json:"_postman_variable_scope"`
	PostmanExportedAt    time.Time `json:"_postman_exported_at"`
	PostmanExportedUsing string    `json:"_postman_exported_using"`
}

// Prepend params adds Values from EnvFile on the beginning of array
func (e *EnvFile) PrependParams(from EnvFile) {
	vals := from.Values
	vals = append(vals, e.Values...)
	e.Values = vals
}

type Value struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}
