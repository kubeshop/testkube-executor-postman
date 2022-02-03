module github.com/kubeshop/testkube-executor-postman

go 1.16

// replace github.com/kubeshop/testkube v0.6.4 => ../testkube

require (
	// use beta for now until we merge everything together with job executors
	github.com/kubeshop/testkube v0.8.6-beta009
	github.com/stretchr/testify v1.7.0
)
