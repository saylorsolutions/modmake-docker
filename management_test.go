package modmake_docker

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocker_RemoveImage(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().RemoveImage("some-image:latest").Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "dry run: docker rmi some-image:latest", err.Error())
}

func TestDocker_RemoveImage_Force(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().RemoveImage("some-image:latest").Force().Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "dry run: docker rmi -f some-image:latest", err.Error())
}

func TestDocker_RemoveContainer(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().RemoveContainer("some-container").Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "dry run: docker rm some-container", err.Error())
}

func TestDocker_RemoveContainer_Force(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().RemoveContainer("some-container").Force().Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "dry run: docker rm -f some-container", err.Error())
}

func TestDocker_Start(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Start("some-container").Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "dry run: docker start some-container", err.Error())
}

func TestDocker_Stop(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Stop("some-container").Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "dry run: docker stop some-container", err.Error())
}
