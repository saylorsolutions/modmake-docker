package modmake_docker

import (
	"context"
	"errors"
	"testing"

	. "github.com/saylorsolutions/modmake"
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

func TestDockerRef_Exec(t *testing.T) {
	d := Docker().Dry()
	isDryRunResult(t, d.Exec("some-container", "echo", "Hello!").Task(),
		"docker exec -i some-container echo Hello!")
	isDryRunResult(t, d.Exec("some-container", "blah").InteractiveTerm().Task(),
		"docker exec -i -t some-container blah")
	isDryRunResult(t, d.Exec("some-container", "blah").Detached().Task(),
		"docker exec -d some-container blah")
	isDryRunResult(t, d.Exec("some-container", "blah").Detached().Privileged().Task(),
		"docker exec -d --privileged some-container blah")
	isDryRunResult(t, d.Exec("some-container", "blah").User("bob:admin").Task(),
		"docker exec -i -u bob:admin some-container blah")
	isDryRunResult(t, d.Exec("some-container", "blah").WorkingDir("/app").Task(),
		"docker exec -i -w /app some-container blah")
}

func TestPullRetagPush(t *testing.T) {
	d := Docker().Dry()
	isDryRunResult(t, d.Login("some-host.com").Username("bob").Password(F("${SOME_SECRET_VAR:secret}")).Task(),
		"docker login -u bob -p secret some-host.com")
	isDryRunResult(t, d.Pull("some-host.com/my-image:1"),
		"docker pull some-host.com/my-image:1",
	)
	isDryRunResult(t,
		d.Tag("some-host.com/my-image:1", "some-host.com/my-image:latest"),
		"docker tag some-host.com/my-image:1 some-host.com/my-image:latest",
	)
	isDryRunResult(t,
		d.Push("some-host.com/my-image:latest"),
		"docker push some-host.com/my-image:latest",
	)
}

func isDryRunResult(t *testing.T, task Task, command string) {
	err := task.Run(context.Background())
	assert.Error(t, err)
	var dryRunResult *DryRunResult
	ok := errors.As(err, &dryRunResult)
	assert.True(t, ok, "Should have been a dry run error")
	assert.Equal(t, "dry run: "+command, err.Error())
}
