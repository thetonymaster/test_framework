package container

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
)

type Provider interface {
	Execute(target, task string) error
	Kill() error
	Run() error
	Scale(containers map[string]int) error
}

type Container struct {
	Provider
}

// DockerComposeGenerator generates
type DockerComposeGenerator struct {
	Paths []string
}

// NewDockerComposeGenerator generates new DockerCompose instances
func NewDockerComposeGenerator(paths []string) *DockerComposeGenerator {
	return &DockerComposeGenerator{
		Paths: paths,
	}
}

// New returns a new container provider with docker compose
func (gen DockerComposeGenerator) New(projectName string) *Container {
	compose := NewDockerCompose(projectName, gen.Paths)
	return &Container{
		compose,
	}
}

// DockerCompose runs docker compose tasks
type DockerCompose struct {
	project project.APIProject
}

// NewDockerCompose creates a
func NewDockerCompose(projectName string, files []string) *DockerCompose {
	context := &ctx.Context{
		Context: project.Context{
			ComposeFiles: files,
			ProjectName:  projectName,
		},
	}
	project, err := docker.NewProject(context, nil)

	if err != nil {
		log.Fatal(err)
	}

	return &DockerCompose{
		project: project,
	}
}

func (dc DockerCompose) Scale(containers map[string]int) error {
	return dc.project.Scale(context.TODO(), 300, containers)
}

// Run runs the docker compose files
func (dc DockerCompose) Run() error {
	err := dc.project.Up(context.Background(), options.Up{})
	if err != nil {
		return err
	}
	err = dc.project.Kill(context.Background(), "SIGKILL", "petclinic")
	time.Sleep(5 * time.Second)

	return err
}

// Kill kills all the containers related to the project
func (dc DockerCompose) Kill() error {
	return dc.project.Down(context.TODO(), options.Down{
		RemoveImages:  "local",
		RemoveOrphans: true,
		RemoveVolume:  true,
	})
}

// Execute Executes a task in a container
func (dc DockerCompose) Execute(target, task string) error {
	out, errs := dc.runTask(target, task)
	select {
	case err := <-errs:
		return err
	case status := <-out:
		if status != 0 {
			return fmt.Errorf("Container exited with status %d", status)
		}
	}
	return nil
}

func (dc DockerCompose) runTask(target, task string) (<-chan int, <-chan error) {
	out := make(chan int, 1)
	errs := make(chan error, 1)

	go func() {
		log.Println("Running tests for: " + task)
		status, err := dc.project.Run(context.TODO(),
			"petclinic",
			[]string{"./mvnw", "test", fmt.Sprintf("-Dtest=%s", task)},
			options.Run{
				Detached: false,
			})

		if err != nil {
			errs <- err
		} else {
			out <- status
		}
		close(out)
		close(errs)
	}()

	return out, errs
}
