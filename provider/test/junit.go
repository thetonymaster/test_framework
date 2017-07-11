package test

import (
	"errors"
	"fmt"
	"log"
)

type provider interface {
	Run() error
	Execute(target, task string) error
	Kill() error
}

type generator interface {
	New(projectName string) interface
}

// JUnit runs the JUnit tests
type JUnit struct {
	Generator    generator
	noContainers int
}

const JUnitProject = "junit"

// NewJUnit creates a new instance of a JUnit task manager
func NewJUnit(generator generator) *JUnit {
	return &JUnit{
		Generator:    generator,
		noContainers: 0,
	}
}

func (junit *JUnit) RunTask() error {
	junit.noContainers = junit.noContainers + 1
	projectName := fmt.Sprintf("%s_%d", JUnitProject, junit.noContainers)

	project := junit.Generator.New(projectName)
	err := project.Run()
	if err != nil {
		log.Println(err)
	}

	err = project.Execute(projectName, "asdf")
	if err != nil {
		log.Println(err)
		err = project.Kill()
		if err != nil {
			return errors.New("Cannot kill containers")
		}

		junit.noContainers = junit.noContainers - 1
		return errors.New("Cannot run tasks")
	}

	err = project.Kill()
	if err != nil {
		return errors.New("Cannot kill containers")
	}

	junit.noContainers = junit.noContainers - 1
	return nil
}
