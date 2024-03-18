package modmake_docker

import (
	"context"
	"fmt"
	"testing"

	. "github.com/saylorsolutions/modmake"
	"github.com/stretchr/testify/assert"
)

func TestDockerRun_Run(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Run("some-image:latest").
		Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "docker run some-image:latest", err.Error())
}

func TestDockerRun_Run_Name(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Run("some-image:latest").
		Name("some-container").
		Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "docker run --name=some-container some-image:latest", err.Error())
}

func TestDockerRun_Run_Hostname(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Run("some-image:latest").
		SetHostname("some-container").
		Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "docker run -h some-container some-image:latest", err.Error())
}

func TestDockerRun_Run_WorkingDir(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Run("some-image:latest").
		WorkingDirectory("/app").
		Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "docker run -w /app some-image:latest", err.Error())
}

func TestDockerRun_Run_Detached(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Run("some-image:latest").
		Detached().
		Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "docker run -d some-image:latest", err.Error())
}

func TestDockerRun_Run_Interactive(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Run("some-image:latest").
		InteractiveTerminal().
		Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "docker run -it some-image:latest", err.Error())
}

func TestDockerRun_Run_Privileged(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Run("some-image:latest").
		PrivilegedContainer().
		Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "docker run --privileged some-image:latest", err.Error())
}

func TestDockerRun_Run_ReadOnly(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Run("some-image:latest").
		ReadOnlyFS().
		Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "docker run --read-only some-image:latest", err.Error())
}

func TestDockerRun_Run_Remove(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Run("some-image:latest").
		RemoveAfterExit().
		Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "docker run --rm some-image:latest", err.Error())
}

func TestDockerRun_Run_RestartPolicy(t *testing.T) {
	ctx := context.Background()
	tests := map[string]struct {
		policy   RestartPolicy
		expected string
	}{
		"Never restart which is the default": {
			policy:   RestartNever,
			expected: "docker run some-image:latest",
		},
		"Restart on failure": {
			policy:   RestartOnFailure,
			expected: "docker run --restart=on-failure some-image:latest",
		},
		"Restart unless stopped": {
			policy:   RestartUnlessStopped,
			expected: "docker run --restart=unless-stopped some-image:latest",
		},
		"Restart always": {
			policy:   RestartAlways,
			expected: "docker run --restart=always some-image:latest",
		},
		"Invalid policy": {
			policy:   RestartPolicy("something-else"),
			expected: "missing required parameter: unknown restart policy 'something-else'",
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			err := Docker().Dry().Run("some-image:latest").
				SetRestartPolicy(tc.policy).
				Run(ctx)
			assert.Error(t, err)
			assert.Equal(t, tc.expected, err.Error())
		})
	}
}

func TestDockerRun_Run_RestartRetries(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Run("some-image:latest").
		SetRestartRetries(5).
		Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "docker run --restart=on-failure:5 some-image:latest", err.Error())
}

func TestDockerRun_Run_Network(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Run("some-image:latest").
		ConnectNetwork("somenet").
		Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "docker run --network=somenet some-image:latest", err.Error())
}

func TestDockerRun_Run_EnvVars(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Run("some-image:latest").
		SetEnvVar("c", "d").
		SetEnvVar("a", "b").
		Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "docker run -e c=d -e a=b some-image:latest", err.Error())
}

func TestDockerRun_Run_Ports(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Run("some-image:latest").
		PublishPort(8080, 8081).
		Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "docker run -p 8080:8081 some-image:latest", err.Error())
}

func TestDockerRun_Run_Volumes(t *testing.T) {
	hostPath, _ := Path("./mount").Abs()
	containerPath := Path("/app/mount")
	ctx := context.Background()
	err := Docker().Dry().Run("some-image:latest").
		VolumeMount(hostPath, containerPath).
		Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("docker run -v %s:%s some-image:latest", hostPath.String(), containerPath.ToSlash()), err.Error())
}

func TestDockerRun_Run_Args(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Run("some-image:latest", "cmd", "arg1").Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "docker run some-image:latest cmd arg1", err.Error())
}
