package modmake_docker

import (
	"errors"
	"os/exec"

	. "github.com/saylorsolutions/modmake"
)

var (
	ErrNoDockerFound = errors.New("unable to locate docker executable")
	ErrRequiredParam = errors.New("missing required parameter")
)

// DockerRef is a reference to the DockerRef CLI, that can then be used to run commands.
type DockerRef struct {
	exePath PathString
}

// Docker will attempt to locate the Docker CLI, and return a DockerRef if successful.
// If the Docker CLI cannot be located from the PATH, then ErrNoDockerFound will be returned.
func Docker() *DockerRef {
	_path, err := exec.LookPath("docker")
	if err != nil {
		panic(ErrNoDockerFound)
	}
	return &DockerRef{
		exePath: Path(_path),
	}
}
