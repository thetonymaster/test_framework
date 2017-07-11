package main

import (
	"github.com/thetonymaster/test_catridge/provider/container"
	"github.com/thetonymaster/test_catridge/provider/test"
)

func main() {
	containerProvider := container.NewDockerComposeGenerator([]string{""})
	jUnitTestProvider := test.NewJUnit(containerProvider)
}
