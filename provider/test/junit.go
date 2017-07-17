package test

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/thetonymaster/test_catridge/provider/container"
)

type provider interface {
	Run() error
	Execute(target, task string) error
	Scale(containers map[string]int) error
	Kill() error
}

type generator interface {
	New(projectName string) *container.Container
}

// JUnit runs the JUnit tests
type JUnit struct {
	Generator generator
}

const JUnitProject = "junit"

// NewJUnit creates a new instance of a JUnit task manager
func NewJUnit(generator generator) *JUnit {
	return &JUnit{
		Generator: generator,
	}
}

func (junit JUnit) GetFiles(path string) {
	files, _ := filepath.Glob("*")
	fmt.Println(files)
}

func (junit *JUnit) RunTask() error {
	containers := junit.Generator.New(JUnitProject)
	containers.Run()
	containers.Scale(map[string]int{
		"petclinic": 2,
	})
	time.Sleep(time.Second * 30)
	containers.Kill()
	return nil
}
