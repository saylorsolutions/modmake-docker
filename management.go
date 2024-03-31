package modmake_docker

import (
	"context"
	"fmt"
	"strings"

	. "github.com/saylorsolutions/modmake"
)

// RemoveImage will attempt to remove an image from the local repository.
// Use [DockerRemoveImage.Force] to allow removing currently referenced images.
func (d *DockerRef) RemoveImage(image string) *DockerRemoveImage {
	image = strings.TrimSpace(image)
	if len(image) == 0 {
		return &DockerRemoveImage{err: fmt.Errorf("%w: missing image name/hash", ErrRequiredParam)}
	}
	return &DockerRemoveImage{
		d:     d,
		image: image,
	}
}

type DockerRemoveImage struct {
	d     *DockerRef
	err   error
	image string
	force bool
}

// Force will force removal of the referenced Docker image.
func (r *DockerRemoveImage) Force() *DockerRemoveImage {
	if r.err != nil {
		return r
	}
	r.force = true
	return r
}

func (r *DockerRemoveImage) Task() Task {
	return r.Run
}

func (r *DockerRemoveImage) Run(ctx context.Context) error {
	if r.err != nil {
		return r.err
	}
	args := []string{"rmi"}
	if r.force {
		args = append(args, "-f")
	}
	return r.d.Command(append(args, r.image)...).Run(ctx)
}

// RemoveContainer will attempt to remove a container.
// Use [DockerRemoveContainer.Force] to force remove, which can remove running containers.
func (d *DockerRef) RemoveContainer(name string) *DockerRemoveContainer {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return &DockerRemoveContainer{err: fmt.Errorf("%w: missing container name", ErrRequiredParam)}
	}
	return &DockerRemoveContainer{
		d:    d,
		name: name,
	}
}

type DockerRemoveContainer struct {
	d     *DockerRef
	err   error
	name  string
	force bool
}

// Force will force removal of the referenced Docker container.
func (r *DockerRemoveContainer) Force() *DockerRemoveContainer {
	if r.err != nil {
		return r
	}
	r.force = true
	return r
}

func (r *DockerRemoveContainer) Task() Task {
	return r.Run
}

func (r *DockerRemoveContainer) Run(ctx context.Context) error {
	if r.err != nil {
		return r.err
	}
	args := []string{"rm"}
	if r.force {
		args = append(args, "-f")
	}
	return r.d.Command(append(args, r.name)...).Run(ctx)
}

// Stop will attempt to stop a running container with the given name.
func (d *DockerRef) Stop(name string) Task {
	return d.Command("stop", name).Run
}

// Start will attempt to start a container with the given name.
func (d *DockerRef) Start(name string) Task {
	return d.Command("start", name).Run
}

func (d *DockerRef) Exec(containerName string, cmdAndArgs ...string) Task {
	cmd := append([]string{"exec", "-it", containerName}, cmdAndArgs...)
	return d.Command(cmd...)
}

func (d *DockerRef) Login(host string) Task {
	return d.Command("login", host)
}

func (d *DockerRef) Pull(imageAndTag string) Task {
	return d.Command("pull", imageAndTag)
}

func (d *DockerRef) Tag(currentTag, newTag string) Task {
	return d.Command("tag", currentTag, newTag)
}

func (d *DockerRef) Push(imageAndTag string) Task {
	return d.Command("push", imageAndTag)
}
