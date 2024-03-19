package modmake_docker

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	. "github.com/saylorsolutions/modmake"
)

var (
	ErrNoDockerFound = errors.New("unable to locate docker executable")
	ErrRequiredParam = errors.New("missing required parameter")
)

var _ error = (*DryRunResult)(nil)

type DryRunResult struct {
	args []string
}

func dryRunExec(args ...string) Task {
	return func(ctx context.Context) error {
		return &DryRunResult{args: args}
	}
}

func (r *DryRunResult) Error() string {
	return fmt.Sprintf("dry run: docker " + strings.Join(r.args, " "))
}

func (r *DryRunResult) Args() []string {
	return r.args
}

// DockerRef is a reference to the DockerRef CLI, that can then be used to run commands.
type DockerRef struct {
	exePath PathString
	dryRun  bool
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

func (d *DockerRef) Dry() *DockerRef {
	d.dryRun = true
	return d
}

func (d *DockerRef) Command(args ...string) Task {
	if d.dryRun {
		return dryRunExec(args...)
	}
	if d.exePath != Path("") {
		_path, err := exec.LookPath("docker")
		if err != nil {
			return Error("%w: unable to locate docker, and this is not a dry run", ErrNoDockerFound)
		}
		d.exePath = Path(_path)
	}
	return Exec(append([]string{d.exePath.String()}, args...)...).CaptureStdin().Run
}
