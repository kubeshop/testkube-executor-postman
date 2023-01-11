package main

import (
	"fmt"
	"os"

	"github.com/kubeshop/testkube-executor-postman/pkg/runner/newman"
	"github.com/kubeshop/testkube/pkg/executor/agent"
	"github.com/kubeshop/testkube/pkg/ui"
)

func main() {
	r, err := newman.NewNewmanRunner()
	if err != nil {
		panic(fmt.Errorf("%s could not run Postman tests: %w", ui.IconCross, err))
	}
	agent.Run(r, os.Args)
}
