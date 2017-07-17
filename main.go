package main

import (
	"os"

	"github.com/thetonymaster/test_catridge/provider/container"
	"github.com/thetonymaster/test_catridge/provider/test"
)

func main() {

	args := os.Args[1:]
	containerProvider := container.NewDockerComposeGenerator(args)
	jUnitTestProvider := test.NewJUnit(containerProvider)
	jUnitTestProvider.RunTask()
	jUnitTestProvider.GetFiles("Test_Framework_Application")

}
