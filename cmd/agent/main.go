package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kubeshop/testkube-executor-postman/pkg/runner/newman"
	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
)

func main() {

	args := os.Args
	if len(args) == 1 {
		fmt.Println("missing input argument")
		os.Exit(1)
	}

	script := args[1]

	e := testkube.Execution{}
	json.Unmarshal([]byte(script), &e)
	runner := newman.NewNewmanRunner()
	result, err := runner.Run(e)
	fmt.Println(result)
	fmt.Printf("$$$%s$$$", e.Id)
}
