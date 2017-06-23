package container

import (
	"context"
	"log"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
)

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

// Run runs the docker compose files
func (dc DockerCompose) Run() error {
	return dc.project.Up(context.Background(), options.Up{})
}
