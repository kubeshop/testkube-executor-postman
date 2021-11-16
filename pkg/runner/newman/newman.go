package newman

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/executor"
	"github.com/kubeshop/testkube/pkg/tmp"
)

func NewNewmanRunner() *NewmanRunner {
	return &NewmanRunner{}
}

// NewmanRunner struct for newman based runner
type NewmanRunner struct {
}

// Run runs particular script content on top of newman binary
func (r *NewmanRunner) Run(execution testkube.Execution) (result testkube.ExecutionResult, err error) {

	input := strings.NewReader(execution.ScriptContent)

	path, err := tmp.ReaderToTmpfile(input)
	if err != nil {
		return result, err
	}

	// write params to tmp file
	envReader, err := NewEnvFileReader(execution.Params)
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

	// we'll get error here in case of failed test too so we treat this as
	// starter test execution with failed status
	out, err := executor.Run("", "newman", args...)

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

	// parse JSON output of newman script
	bytes, err := ioutil.ReadFile(tmpName)
	if err != nil {
		return newmanResult, err
	}

	err = json.Unmarshal(bytes, &newmanResult.Metadata)
	if err != nil {
		return newmanResult, err
	}

	return
}
