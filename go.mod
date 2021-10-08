module github.com/kubeshop/testkube-executor-postman

go 1.16

// replace github.com/kubeshop/testkube-operator v0.1.1 => ../testkube-operator
// replace github.com/kubeshop/testkube v0.5.22 => ../testkube

require (
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kubeshop/testkube v0.5.43
	github.com/stretchr/testify v1.7.0
)
