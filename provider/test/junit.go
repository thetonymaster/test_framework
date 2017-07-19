package test

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/jeffail/tunny"
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
	pool      *tunny.WorkPool
}

const JUnitProject = "junit"

// NewJUnit creates a new instance of a JUnit task manager
func NewJUnit(generator generator, pool *tunny.WorkPool) *JUnit {
	return &JUnit{
		Generator: generator,
		pool:      pool,
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

func (junit *JUnit) RunTask(tasks []string) error {
	containers := junit.Generator.New(JUnitProject)
	containers.Run()

	for _, task := range tasks {
		t := func() {
			log.Println("Running tests for file " + task)
			err := containers.Execute("petclinic", task)
			time.Sleep(3 * time.Second)
			fmt.Println(err)
			time.Sleep(3 * time.Second)
		}
		junit.pool.SendWorkAsync(t, nil)
		time.Sleep(1 * time.Second)

	}
	for junit.pool.NumPendingAsyncJobs() > 0 {
		time.Sleep(1 * time.Second)
	}
	containers.Kill()
	return nil
}
