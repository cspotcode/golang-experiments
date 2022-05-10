module github.com/cspotcode/golang-experiments/docker-credential

go 1.15

// https://github.com/cspotcode/docker-cli/tree/as-module
// replace github.com/docker/cli => ../docker-cli

require (
	github.com/aws/aws-sdk-go v1.36.0
	github.com/awslabs/amazon-ecr-credential-helper/ecr-login v0.0.0-20201113001948-d77edb6d2e47
	github.com/docker/cli v20.10.15+incompatible
	github.com/docker/docker v1.13.1
	github.com/docker/docker-credential-helpers v0.6.3 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/spf13/cobra v1.1.1
	golang.org/x/sys v0.0.0-20220503163025-988cb79eb6c6 // indirect
	golang.org/x/tools/gopls v0.5.4 // indirect
)
