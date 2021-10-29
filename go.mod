module github.com/kubeshop/testkube-executor-postman

go 1.16

// replace github.com/kubeshop/testkube-operator v0.1.1 => ../testkube-operator
// replace github.com/kubeshop/testkube v0.5.59 => ../testkube

require (
	github.com/kelseyhightower/envconfig v1.4.0
	// use beta for now until we merge everuthing together with job executors
	github.com/kubeshop/testkube v0.6.4-beta1
	github.com/stretchr/testify v1.7.0
)
