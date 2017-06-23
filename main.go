package main

import (
	"log"
	"os"

	"github.com/thetonymaster/test_catridge/provider/container"
)

type provider interface {
	Run() error
}

func main() {
	args := os.Args[1:]

	projectName := args[0]
	files := args[1:]

	provider := container.NewDockerCompose(projectName, files)
	provisionContainers(provider)

}

func provisionContainers(pr provider) {
	err := pr.Run()
	log.Fatal(err)
}
