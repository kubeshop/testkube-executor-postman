package main

import (
	"os"

	"github.com/kubeshop/testkube-executor-postman/pkg/runner/newman"
	"github.com/kubeshop/testkube/pkg/runner/agent"
)

func main() {
	agent.Run(newman.NewNewmanRunner(), os.Args)
}
