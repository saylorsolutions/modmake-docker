package modmake_docker

import (
	"context"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDockerBuild_Build(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Build("some-image:latest", ".").Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "dry run: docker build -t some-image:latest .", err.Error())
}

func TestDockerBuild_Build_NoContext(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Build("some-image:latest", "").Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "dry run: docker build -t some-image:latest .", err.Error())
}

func TestDockerBuild_Build_Context(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Build("some-image:latest", "./test_ctx").Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "dry run: docker build -t some-image:latest .", err.Error())
}

func TestDockerBuild_BuildArg(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Build("some-image:latest", ".").
		BuildArg("c", "d").
		BuildArg("a", "b").
		Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "dry run: docker build -t some-image:latest --build-arg c=d --build-arg a=b .", err.Error())
}

func TestDockerBuild_BuildFileRelative(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Build("some-image:latest", "./test_ctx").
		BuildFile("./test_ctx/Dockerfile").
		Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "dry run: docker build -t some-image:latest -f Dockerfile .", err.Error())
}

func TestDockerBuild_SetLabel(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Build("some-image:latest", ".").
		Label("c", "d").
		Label("a", "b").
		Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "dry run: docker build -t some-image:latest --label c=d --label a=b .", err.Error())
}

func TestDockerBuild_LabelBuildTimestamp(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Build("some-image:latest", ".").
		LabelBuildTimestamp().
		Run(ctx)
	assert.Error(t, err)
	pat := regexp.MustCompile(`dry run: docker build -t some-image:latest --label buildTimestamp=.+ \.`)
	assert.True(t, pat.MatchString(err.Error()), "Should match: '%s'", err.Error())
}
