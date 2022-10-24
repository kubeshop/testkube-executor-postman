package newman

import (
	"encoding/json"
	"os"

	"github.com/kelseyhightower/envconfig"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/executor"
	"github.com/kubeshop/testkube/pkg/executor/content"
	"github.com/kubeshop/testkube/pkg/executor/secret"
	"github.com/kubeshop/testkube/pkg/tmp"
)

// Params ...
type Params struct {
	Endpoint        string // RUNNER_ENDPOINT
	AccessKeyID     string // RUNNER_ACCESSKEYID
	SecretAccessKey string // RUNNER_SECRETACCESSKEY
	Location        string // RUNNER_LOCATION
	Token           string // RUNNER_TOKEN
	Ssl             bool   // RUNNER_SSL
	ScrapperEnabled bool   // RUNNER_SCRAPPERENABLED
	GitUsername     string // RUNNER_GITUSERNAME
	GitToken        string // RUNNER_GITTOKEN
}

func NewNewmanRunner() *NewmanRunner {
	var params Params
	err := envconfig.Process("runner", &params)
	if err != nil {
		panic(err.Error())
	}

	return &NewmanRunner{
		Params:  params,
		Fetcher: content.NewFetcher(""),
	}
}

// NewmanRunner struct for newman based runner
type NewmanRunner struct {
	Params  Params
	Fetcher content.ContentFetcher
}

// Run runs particular test content on top of newman binary
func (r *NewmanRunner) Run(execution testkube.Execution) (result testkube.ExecutionResult, err error) {
	if r.Params.GitUsername != "" && r.Params.GitToken != "" {
		if execution.Content != nil && execution.Content.Repository != nil {
			execution.Content.Repository.Username = r.Params.GitUsername
			execution.Content.Repository.Token = r.Params.GitToken
		}
	}

	path, err := r.Fetcher.Fetch(execution.Content)
	if err != nil {
		return result, err
	}

	if !execution.Content.IsFile() {
		return result, testkube.ErrTestContentTypeNotFile
	}

	envManager := secret.NewEnvManagerWithVars(execution.Variables)
	// write params to tmp file
	envReader, err := NewEnvFileReader(execution.Variables, execution.VariablesFile, envManager.GetEnvs())
	if err != nil {
		return result, err
	}
	envpath, err := tmp.ReaderToTmpfile(envReader)
	if err != nil {
		return result, err
	}

	tmpName := tmp.Name() + ".json"

	args := []string{
		"run", path, "-e", envpath, "--reporters", "cli,json", "--reporter-json-export", tmpName,
	}
	args = append(args, execution.Args...)

	runPath := ""
	if execution.Content.Repository != nil {
		runPath = execution.Content.Repository.WorkingDir
	}

	// we'll get error here in case of failed test too so we treat this as
	// starter test execution with failed status
	out, err := executor.Run(runPath, "newman", envManager, args...)

	out = envManager.Obfuscate(out)

	// try to get json result even if process returned error (could be invalid test)
	newmanResult, nerr := r.GetNewmanResult(tmpName, out)
	// convert newman result to OpenAPI struct
	result = MapMetadataToResult(newmanResult)

	// catch errors if any
	if err != nil {
		return result.Err(err), nil
	}

	if nerr != nil {
		return result.Err(nerr), nil
	}

	return result, nil
}

func (r NewmanRunner) GetNewmanResult(tmpName string, out []byte) (newmanResult NewmanExecutionResult, err error) {
	newmanResult.Output = string(out)

	// parse JSON output of newman test
	bytes, err := os.ReadFile(tmpName)
	if err != nil {
		return newmanResult, err
	}

	err = json.Unmarshal(bytes, &newmanResult.Metadata)
	if err != nil {
		return newmanResult, err
	}

	return
}
