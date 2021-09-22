package newman

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kubeshop/kubtest/pkg/api/kubtest"
	"github.com/kubeshop/kubtest/pkg/process"
	"github.com/kubeshop/kubtest/pkg/tmp"
)

func NewNewmanRunner() *NewmanRunner {
	return &NewmanRunner{}
}

// NewmanRunner struct for newman based runner
type NewmanRunner struct {
}

// Run runs particular script content on top of newman binary
func (r *NewmanRunner) Run(execution kubtest.Execution) (result kubtest.ExecutionResult) {

	input := strings.NewReader(execution.ScriptContent)

	path, err := tmp.ReaderToTmpfile(input)
	if err != nil {
		return result.Err(err)
	}

	// write params to tmp file
	envReader, err := NewEnvFileReader(execution.Params)
	if err != nil {
		return result.Err(err)
	}
	envpath, err := tmp.ReaderToTmpfile(envReader)
	if err != nil {
		return result.Err(err)
	}

	var newmanResult NewmanExecutionResult

	tmpName := tmp.Name() + ".json"
	out, err := process.LoggedExecuteInDir("", os.Stdout, "newman", "run", path, "-e", envpath, "--reporters", "cli,json", "--reporter-json-export", tmpName)
	if err != nil {
		return result.Err(err)
	}

	newmanResult.Output = string(out)

	// parse JSON output of newman script
	bytes, err := ioutil.ReadFile(tmpName)
	if err != nil {
		return result.Err(err)
	}

	err = json.Unmarshal(bytes, &newmanResult.Metadata)
	if err != nil {
		return result.Err(fmt.Errorf("parsing results metadata error: %w", err))
	}

	// convert newman result to OpenAPI struct
	res := MapMetadataToResult(newmanResult)

	return res
}
