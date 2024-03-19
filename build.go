package modmake_docker

import (
	"context"
	"fmt"
	"strings"
	"time"

	. "github.com/saylorsolutions/modmake"
)

// Build provides a method and options for building a Docker image.
func (d *DockerRef) Build(image string, context PathString) *DockerBuild {
	image = strings.TrimSpace(image)
	if len(image) == 0 {
		return &DockerBuild{err: fmt.Errorf("%w: missing image", ErrRequiredParam)}
	}
	return &DockerBuild{d: d, image: image, context: context}
}

type DockerBuild struct {
	d         *DockerRef
	err       error
	image     string
	context   PathString
	buildFile PathString
	buildArgs []string
	labels    []string
}

// BuildArg sets a build argument for this image build.
// These are distinct from environment variables.
func (b *DockerBuild) BuildArg(key, value string) *DockerBuild {
	if b.err != nil {
		return b
	}
	b.buildArgs = append(b.buildArgs, fmt.Sprintf("%s=%s", key, value))
	return b
}

// BuildFile specifies a Dockerfile for this build.
// This is useful if the Dockerfile isn't named "Dockerfile" in the working context of the build.
func (b *DockerBuild) BuildFile(filePath PathString) *DockerBuild {
	if b.err != nil {
		return b
	}
	b.buildFile = filePath
	return b
}

// Label will set a metadata label in the built image.
func (b *DockerBuild) Label(key, val string) *DockerBuild {
	if b.err != nil {
		return b
	}
	b.labels = append(b.labels, fmt.Sprintf("%s=%s", key, val))
	return b
}

// LabelBuildTimestamp will set a build time timestamp in the built image.
func (b *DockerBuild) LabelBuildTimestamp() *DockerBuild {
	if b.err != nil {
		return b
	}
	return b.Label("buildTimestamp", time.Now().Format(time.RFC3339))
}

func (b *DockerBuild) Task() Task {
	return b.Run
}

func (b *DockerBuild) Run(ctx context.Context) error {
	if b.err != nil {
		return b.err
	}
	var (
		args    = []string{"build", "-t", b.image}
		chdir   PathString
		doChdir bool
	)
	if b.context.IsDir() {
		chdir = b.context
		doChdir = true
	}
	if len(b.buildFile.String()) > 0 {
		if doChdir {
			rel, err := b.context.Rel(b.buildFile)
			if err != nil {
				return fmt.Errorf("unable to make a relative path from '%s' to '%s'", b.context, b.buildFile)
			}
			args = append(args, "-f", rel.String())
		} else {
			args = append(args, "-f", b.buildFile.String())
		}
	}

	for _, arg := range b.buildArgs {
		args = append(args, "--build-arg", arg)
	}

	for _, label := range b.labels {
		args = append(args, "--label", label)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		cmd := b.d.Command(append(args, ".")...)
		if doChdir {
			return Chdir(chdir, cmd).Run(ctx)
		}
		return cmd.Run(ctx)
	}
}
