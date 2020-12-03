module github.com/cspotcode/golang-experiments/docker-credential

go 1.15

replace github.com/docker/cli => ../docker-cli

// replace cli => ../docker-cli/cli

require (
	github.com/aws/aws-sdk-go v1.36.0
	github.com/awslabs/amazon-ecr-credential-helper/ecr-login v0.0.0-20201113001948-d77edb6d2e47 // indirect
	github.com/docker/cli v20.10.0-rc2+incompatible
	github.com/docker/docker v1.13.1
	github.com/docker/docker-credential-helpers v0.6.3 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/spf13/cobra v1.1.1
	golang.org/x/sys v0.0.0-20201202213521-69691e467435 // indirect
	golang.org/x/tools/gopls v0.5.4 // indirect
)
