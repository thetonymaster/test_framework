package test

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

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

func (junit JUnit) GetFiles(searchDir string) []string {
	fileList := []string{}
	pattern := "(.+?)((Tests.java))"

	filepath.Walk(searchDir, func(filePath string, f os.FileInfo, err error) error {
		match, _ := regexp.MatchString(pattern, filePath)
		if match {
			name := strings.TrimSuffix(path.Base(filePath), filepath.Ext(filePath))
			fileList = append(fileList, name)
		}
		return nil
	})
	return fileList
}

func (junit *JUnit) RunTask() error {
	containers := junit.Generator.New(JUnitProject)
	containers.Run()
	containers.Scale(map[string]int{
		"petclinic": 2,
	})
	containers.Kill()
	return nil
}
