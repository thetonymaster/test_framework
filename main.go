package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jeffail/tunny"
	"github.com/thetonymaster/test_catridge/configuration"
	"github.com/thetonymaster/test_catridge/provider/container"
	"github.com/thetonymaster/test_catridge/provider/test"
)

func main() {

	args := os.Args[1:]

	conf, err := configuration.Read(args[0])
	if err != nil {
		log.Fatal(err)
	}

	pool, _ := tunny.CreatePool(conf.Containers.Limit, func(f interface{}) interface{} {
		input, _ := f.(func())
		input()
		return nil
	}).Open()
	defer pool.Close()

	dir, _ := filepath.Abs(filepath.Dir(conf.Tests["junit"].Path))
	containerProvider := container.NewDockerComposeGenerator([]string{conf.Tests["junit"].Path})
	jUnitTestProvider := test.NewJUnit(containerProvider, conf.Tests["junit"].Target, pool)
	tasks := jUnitTestProvider.GetFiles(dir + "/src/test/")
	jUnitTestProvider.RunTask(tasks)

}
