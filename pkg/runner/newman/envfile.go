package newman

import (
	"bytes"
	"encoding/json"
	"io"
	"time"
)

func NewEnvFileReader(m map[string]string, paramsFile string, secrets []string) (io.Reader, error) {
	envFile := NewEnvFileFromMap(m)

	if paramsFile != "" {
		// create env structure from passed params file
		envFromParamsFile, err := NewEnvFileFromString(paramsFile)
		if err != nil {
			return nil, err
		}
		envFile.PrependParams(envFromParamsFile)
	}

	for _, secret := range secrets {
		// create env structure from passed secret
		envFromSecret, err := NewEnvFileFromString(secret)
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

func NewEnvFileFromMap(m map[string]string) (envFile EnvFile) {
	envFile.ID = "executor-env-file"
	envFile.Name = "executor-env-file"
	envFile.PostmanVariableScope = "environment"
	envFile.PostmanExportedAt = time.Now()
	envFile.PostmanExportedUsing = "Postman/7.34.0"

	for k, v := range m {
		envFile.Values = append(envFile.Values, Value{Key: k, Value: v, Enabled: true})
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
