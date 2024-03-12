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

// Docker is a reference to the Docker CLI, that can then be used to run commands.
type Docker struct {
	exePath PathString
}

// NewDocker will attempt to locate the Docker CLI, and return a Docker if successful.
// If the Docker CLI cannot be located from the PATH, then ErrNoDockerFound will be returned.
func NewDocker() (*Docker, error) {
	_path, err := exec.LookPath("docker")
	if err != nil {
		return nil, ErrNoDockerFound
	}
	return &Docker{
		exePath: Path(_path),
	}, nil
}
