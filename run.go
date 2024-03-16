package modmake_docker

import (
	"context"
	"fmt"
	"strings"

	. "github.com/saylorsolutions/modmake"
)

// Create a new [DockerRun] instance.
func (d *DockerRef) Run(image string, args ...string) *DockerRun {
	r := &DockerRun{
		d:             d,
		image:         image,
		args:          args,
		restartPolicy: RestartNever,
	}
	if len(image) == 0 {
		r.err = fmt.Errorf("%w: missing image", ErrRequiredParam)
	}
	return r
}

// DockerRun encapsulates a "docker run" command.
// It's created from a Docker instance.
type DockerRun struct {
	d               *DockerRef
	name            string
	image           string
	args            []string
	err             error
	hostname        string
	workingDir      string
	networkConn     string
	detached        bool
	interactive     bool
	privileged      bool
	readOnly        bool
	removeAfterExit bool
	restartPolicy   RestartPolicy
	env             []string
	portMappings    []string
	bindMounts      []string
}

// Detached runs the container detached, printing the container ID instead of writing logs to STDOUT.
func (r *DockerRun) Detached() *DockerRun {
	if r.err != nil {
		return r
	}
	r.detached = true
	return r
}

func (r *DockerRun) SetEnvVar(key, val string) *DockerRun {
	if r.err != nil {
		return r
	}
	_val := fmt.Sprintf("%s=%s", strings.TrimSpace(key), strings.TrimSpace(val))
	r.env = append(r.env, _val)
	return r
}

func (r *DockerRun) PublishPort(host, container int) *DockerRun {
	if r.err != nil {
		return r
	}
	if host < 1 || container < 1 {
		r.err = fmt.Errorf("%w: invalid port value '%d:%d'", ErrRequiredParam, host, container)
		return r
	}
	r.portMappings = append(r.portMappings, fmt.Sprintf("%d:%d", host, container))
	return r
}

func (r *DockerRun) SetHostname(host string) *DockerRun {
	if r.err != nil {
		return r
	}
	host = strings.TrimSpace(host)
	if len(host) == 0 {
		r.err = fmt.Errorf("%w: missing host name '%s'", ErrRequiredParam, host)
		return r
	}
	r.hostname = host
	return r
}

func (r *DockerRun) InteractiveTerminal() *DockerRun {
	if r.err != nil {
		return r
	}
	r.interactive = true
	return r
}

func (r *DockerRun) Name(name string) *DockerRun {
	if r.err != nil {
		return r
	}
	r.name = name
	return r
}

func (r *DockerRun) ConnectNetwork(network string) *DockerRun {
	if r.err != nil {
		return r
	}
	_network := strings.TrimSpace(network)
	if len(_network) == 0 {
		r.err = fmt.Errorf("%w: invalid network param '%s'", ErrRequiredParam, network)
		return r
	}
	r.networkConn = _network
	return r
}

func (r *DockerRun) PrivilegedContainer() *DockerRun {
	if r.err != nil {
		return r
	}
	r.privileged = true
	return r
}

func (r *DockerRun) ReadOnlyFS() *DockerRun {
	if r.err != nil {
		return r
	}
	r.readOnly = true
	return r
}

// RestartPolicy is used with DockerRun to establish a restart policy for a given container.
type RestartPolicy string

const (
	RestartNever         RestartPolicy = "no"
	RestartOnFailure     RestartPolicy = "on-failure"
	RestartUnlessStopped RestartPolicy = "unless-stopped"
	RestartAlways        RestartPolicy = "always"
)

var knownRestartPolicies = map[RestartPolicy]struct{}{
	RestartNever:         {},
	RestartOnFailure:     {},
	RestartUnlessStopped: {},
	RestartAlways:        {},
}

func (r *DockerRun) SetRestartPolicy(policy RestartPolicy) *DockerRun {
	if r.err != nil {
		return r
	}
	if _, ok := knownRestartPolicies[policy]; !ok {
		r.err = fmt.Errorf("%w: unknown restart policy '%s'", ErrRequiredParam, policy)
		return r
	}
	if policy != RestartNever {
		r.removeAfterExit = false
	}
	r.restartPolicy = policy
	return r
}

func (r *DockerRun) SetRestartRetries(retries int) *DockerRun {
	if r.err != nil {
		return r
	}
	if retries < 1 {
		r.err = fmt.Errorf("%w: invalid retries '%d'", ErrRequiredParam, retries)
		return r
	}
	r.restartPolicy = RestartPolicy(fmt.Sprintf("%s:%d", string(RestartOnFailure), retries))
	return r
}

func (r *DockerRun) RemoveAfterExit() *DockerRun {
	if r.err != nil {
		return r
	}
	r.removeAfterExit = true
	r.restartPolicy = RestartNever
	return r
}

func (r *DockerRun) VolumeMount(hostPath, containerPath PathString) *DockerRun {
	if r.err != nil {
		return r
	}
	absHostPath, err := hostPath.Abs()
	if err != nil {
		r.err = fmt.Errorf("failed to get absolute path for host bind mount: %w", err)
		return r
	}
	r.bindMounts = append(r.bindMounts, fmt.Sprintf("%s:%s", absHostPath.String(), containerPath.ToSlash()))
	return r
}

func (r *DockerRun) WorkingDirectory(containerPath PathString) *DockerRun {
	if r.err != nil {
		return r
	}
	r.workingDir = containerPath.ToSlash()
	return r
}

func (r *DockerRun) Task() Task {
	return r.Run
}

func (r *DockerRun) Run(ctx context.Context) error {
	if r.err != nil {
		return r.err
	}
	args := []string{r.d.exePath.String(), "run"}
	if len(r.name) > 0 {
		args = append(args, "--name="+r.name)
	}
	if len(r.hostname) > 0 {
		args = append(args, "-h", r.hostname)
	}
	if len(r.workingDir) > 0 {
		args = append(args, "-w", r.workingDir)
	}
	if r.detached {
		args = append(args, "-d")
	}
	if r.interactive {
		args = append(args, "-it")
	}
	if r.privileged {
		args = append(args, "--privileged")
	}
	if r.readOnly {
		args = append(args, "--read-only")
	}
	if r.removeAfterExit {
		args = append(args, "--rm")
	} else if r.restartPolicy != RestartNever {
		args = append(args, "--restart="+string(r.restartPolicy))
	}
	if len(r.networkConn) > 0 {
		args = append(args, "--network="+r.networkConn)
	}

	for _, env := range r.env {
		args = append(args, "-e", env)
	}

	for _, port := range r.portMappings {
		args = append(args, "-p", port)
	}

	for _, bind := range r.bindMounts {
		args = append(args, "-v", bind)
	}

	args = append(args, r.image)
	if len(r.args) > 0 {
		args = append(args, r.args...)
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return Exec(args...).CaptureStdin().Run(ctx)
	}
}
