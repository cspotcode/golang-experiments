module github.com/cspotcode/golang-experiments/docker-credential-cli

go 1.15

replace github.com/docker/cli => ../docker-cli

// replace cli => ../docker-cli/cli

require (
	// require github.com/docker/cli v19.03.0-dev
	// require github.com/docker/cli/cli/config v0.0.0-20200915230204-cd8016b6bcc5
	// require github.com/docker/cli v0.0.0
	// cli v0.0.0
	github.com/docker/cli v0.0.0-00010101000000-000000000000
	github.com/docker/docker v1.13.1 // indirect
	github.com/docker/docker-credential-helpers v0.6.3 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.7.0 // indirect
)
