package main

import (
	"fmt"

	docker_config "github.com/docker/cli/cli/config"
)

func main() {
	fmt.Println(docker_config.Dir())
}
