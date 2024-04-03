package modmake_docker

import (
	"context"
	"fmt"
	"strings"

	. "github.com/saylorsolutions/modmake"
)

// Run creates a new [DockerRun] instance, used to run a container.
func (d *DockerRef) Run(image string, args ...string) *DockerRun {
	anyBlankPanic(strmap{"image": &image})
	r := &DockerRun{
		d:             d,
		image:         image,
		args:          args,
		restartPolicy: RestartNever,
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
	r.detached = true
	return r
}

// SetEnvVar will set an environment variable in the container when it's run.
func (r *DockerRun) SetEnvVar(key, val string) *DockerRun {
	anyBlankPanic(strmap{"key": &key, "val": &val})
	_val := fmt.Sprintf("%s=%s", strings.TrimSpace(key), strings.TrimSpace(val))
	r.env = append(r.env, _val)
	return r
}

// PublishPort will publish a port, mapping a host port to a container port.
// These ports don't have to match.
func (r *DockerRun) PublishPort(host, container int) *DockerRun {
	if host < 1 || container < 1 {
		panicf("invalid port value '%d:%d'", host, container)
		return r
	}
	r.portMappings = append(r.portMappings, fmt.Sprintf("%d:%d", host, container))
	return r
}

// SetHostname will set a host name to be used in the running container.
// Other containers may reference this container (assuming they're on the same network) using this host name.
// By default, the container can be referenced by its container name.
func (r *DockerRun) SetHostname(host string) *DockerRun {
	anyBlankPanic(strmap{"host": &host})
	r.hostname = host
	return r
}

// InteractiveTerminal will allow the specified command to be interactive, allocating a terminal to do so.
func (r *DockerRun) InteractiveTerminal() *DockerRun {
	r.interactive = true
	return r
}

// Name will set the container name when started.
func (r *DockerRun) Name(name string) *DockerRun {
	anyBlankPanic(strmap{"name": &name})
	r.name = name
	return r
}

// ConnectNetwork will allow this container to communicate using the named network.
func (r *DockerRun) ConnectNetwork(network string) *DockerRun {
	anyBlankPanic(strmap{"network": &network})
	r.networkConn = network
	return r
}

// PrivilegedContainer will run this container in "privileged" mode.
// This should not be used unless specifically required by an image constraint, as it opens up the container host to possible security vulnerabilities.
func (r *DockerRun) PrivilegedContainer() *DockerRun {
	r.privileged = true
	return r
}

// ReadOnlyFS will place the container's file system in read only mode.
func (r *DockerRun) ReadOnlyFS() *DockerRun {
	r.readOnly = true
	return r
}

// RestartPolicy is used with DockerRun to establish a restart policy for a given container.
type RestartPolicy string

const (
	RestartNever         RestartPolicy = "no"             // RestartNever will prevent the running container from restarting itself. This is the default.
	RestartOnFailure     RestartPolicy = "on-failure"     // RestartOnFailure will perpetually restart the container when it exits with a non-zero exit code.
	RestartUnlessStopped RestartPolicy = "unless-stopped" // RestartUnlessStopped will perpetually restart the container until it's explicitly stopped.
	RestartAlways        RestartPolicy = "always"         // RestartAlways will always restart the container.
)

var knownRestartPolicies = map[RestartPolicy]struct{}{
	RestartNever:         {},
	RestartOnFailure:     {},
	RestartUnlessStopped: {},
	RestartAlways:        {},
}

// SetRestartPolicy will set a [RestartPolicy] for the running container.
func (r *DockerRun) SetRestartPolicy(policy RestartPolicy) *DockerRun {
	polStr := string(policy)
	anyBlankPanic(strmap{"policy": &polStr})
	policy = RestartPolicy(polStr)
	if _, ok := knownRestartPolicies[policy]; !ok {
		panicf("unknown restart policy '%s'", policy)
		return r
	}
	if policy != RestartNever {
		r.removeAfterExit = false
	}
	r.restartPolicy = policy
	return r
}

// SetRestartRetries will set the [RestartOnFailure] policy, with the given number of retries.
func (r *DockerRun) SetRestartRetries(retries int) *DockerRun {
	if retries < 1 {
		panicf("invalid retries '%d'", retries)
		return r
	}
	r.restartPolicy = RestartPolicy(fmt.Sprintf("%s:%d", string(RestartOnFailure), retries))
	r.removeAfterExit = false
	return r
}

// RemoveAfterExit will cause the container to be removed when it's stopped.
// This is mutually exclusive with a [RestartPolicy].
func (r *DockerRun) RemoveAfterExit() *DockerRun {
	r.removeAfterExit = true
	r.restartPolicy = RestartNever
	return r
}

// VolumeMount will mount a host path at a container path to allow reading/writing in the host file system.
func (r *DockerRun) VolumeMount(hostPath, containerPath PathString) *DockerRun {
	cps := containerPath.String()
	anyBlankPanic(strmap{"containerPath": &cps})
	absHostPath, err := hostPath.Abs()
	if err != nil {
		panicf("failed to get absolute path for host bind mount: %v", err)
		return r
	}
	r.bindMounts = append(r.bindMounts, fmt.Sprintf("%s:%s", absHostPath.String(), containerPath.ToSlash()))
	return r
}

// WorkingDirectory will set the working directory for the container entry point command.
func (r *DockerRun) WorkingDirectory(containerPath PathString) *DockerRun {
	r.workingDir = containerPath.ToSlash()
	return r
}

func (r *DockerRun) Task() Task {
	return r.Run
}

func (r *DockerRun) Run(ctx context.Context) error {
	args := []string{"run"}
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
		return r.d.Command(args...).Run(ctx)
	}
}
