package main

import (
	"os"
	"path/filepath"

	"github.com/jeffail/tunny"
	"github.com/thetonymaster/test_catridge/provider/container"
	"github.com/thetonymaster/test_catridge/provider/test"
)

func main() {

	pool, _ := tunny.CreatePool(1, func(f interface{}) interface{} {
		input, _ := f.(func())
		input()
		return nil
	}).Open()
	defer pool.Close()

	args := os.Args[1:]

	dir, _ := filepath.Abs(filepath.Dir(args[0]))
	containerProvider := container.NewDockerComposeGenerator(args)
	jUnitTestProvider := test.NewJUnit(containerProvider, pool)
	tasks := jUnitTestProvider.GetFiles(dir + "/src/test/")
	jUnitTestProvider.RunTask(tasks)

}
