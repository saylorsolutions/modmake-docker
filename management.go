package modmake_docker

import (
	"context"
	"errors"
	"strings"

	. "github.com/saylorsolutions/modmake"
)

// RemoveImage will attempt to remove an image from the local repository.
// Use [DockerRemoveImage.Force] to allow removing currently referenced images.
func (d *DockerRef) RemoveImage(image string) *DockerRemoveImage {
	anyBlankPanic(strmap{"image": &image})
	return &DockerRemoveImage{
		d:     d,
		image: image,
	}
}

type DockerRemoveImage struct {
	d     *DockerRef
	image string
	force bool
}

// Force will force removal of the referenced Docker image.
func (r *DockerRemoveImage) Force() *DockerRemoveImage {
	r.force = true
	return r
}

func (r *DockerRemoveImage) Task() Task {
	return r.Run
}

func (r *DockerRemoveImage) Run(ctx context.Context) error {
	args := []string{"rmi"}
	if r.force {
		args = append(args, "-f")
	}
	return r.d.Command(append(args, r.image)...).Run(ctx)
}

// RemoveContainer will attempt to remove a container.
// Use [DockerRemoveContainer.Force] to force remove, which can remove running containers.
func (d *DockerRef) RemoveContainer(name string) *DockerRemoveContainer {
	anyBlankPanic(strmap{"name": &name})
	return &DockerRemoveContainer{
		d:    d,
		name: name,
	}
}

type DockerRemoveContainer struct {
	d     *DockerRef
	name  string
	force bool
}

// Force will force removal of the referenced Docker container.
func (r *DockerRemoveContainer) Force() *DockerRemoveContainer {
	r.force = true
	return r
}

func (r *DockerRemoveContainer) Task() Task {
	return r.Run
}

func (r *DockerRemoveContainer) Run(ctx context.Context) error {
	args := []string{"rm"}
	if r.force {
		args = append(args, "-f")
	}
	return r.d.Command(append(args, r.name)...).Run(ctx)
}

// Stop will attempt to stop a running container with the given name.
func (d *DockerRef) Stop(name string) Task {
	anyBlankPanic(strmap{"name": &name})
	return d.Command("stop", name).Run
}

// Start will attempt to start a container with the given name.
func (d *DockerRef) Start(name string) Task {
	anyBlankPanic(strmap{"name": &name})
	return d.Command("start", name).Run
}

type DockerExec struct {
	ref                                    *DockerRef
	containerName                          string
	cmd                                    []string
	detached, interactive, tty, privileged bool
	userGroup                              string
	workingDir                             PathString
}

func (d *DockerRef) Exec(containerName string, cmd string, args ...string) *DockerExec {
	anyBlankPanic(strmap{"containerName": &containerName, "cmd": &cmd})
	return &DockerExec{ref: d, containerName: containerName, cmd: append([]string{cmd}, args...), interactive: true}
}

func (e *DockerExec) Detached() *DockerExec {
	e.detached = true
	e.interactive = false
	e.tty = false
	return e
}

func (e *DockerExec) InteractiveTerm() *DockerExec {
	e.interactive = true
	e.tty = true
	e.detached = false
	return e
}

func (e *DockerExec) Privileged() *DockerExec {
	e.privileged = true
	return e
}

func (e *DockerExec) User(userGroup string) *DockerExec {
	anyBlankPanic(strmap{"userGroup": &userGroup})
	e.userGroup = userGroup
	return e
}

func (e *DockerExec) WorkingDir(wd PathString) *DockerExec {
	if len(wd.String()) == 0 {
		panic("empty working directory path")
	}
	e.workingDir = wd
	return e
}

func (e *DockerExec) Task() Task {
	return e.Run
}

func (e *DockerExec) Run(ctx context.Context) error {
	args := []string{"exec"}
	if e.detached {
		args = append(args, "-d")
	} else {
		if e.interactive {
			args = append(args, "-i")
		}
		if e.tty {
			args = append(args, "-t")
		}
	}
	if e.privileged {
		args = append(args, "--privileged")
	}
	if len(e.userGroup) > 0 {
		args = append(args, "-u", e.userGroup)
	}
	wd := e.workingDir.ToSlash()
	if len(wd) > 0 {
		args = append(args, "-w", wd)
	}
	args = append(append(args, e.containerName), e.cmd...)
	return e.ref.Command(args...).Run(ctx)
}

type DockerLogin struct {
	ref            *DockerRef
	host, username string
	password       []byte
	readStdin      bool
}

func (d *DockerRef) Login(host string) *DockerLogin {
	anyBlankPanic(strmap{"host": &host})
	return &DockerLogin{ref: d, host: host, readStdin: true}
}

func (l *DockerLogin) Username(username string) *DockerLogin {
	username = strings.TrimSpace(username)
	if len(username) == 0 {
		return l
	}
	l.username = username
	return l
}

func (l *DockerLogin) Password(password string) *DockerLogin {
	password = strings.TrimSpace(password)
	if len(password) == 0 {
		return l
	}
	l.password = []byte(password)
	l.readStdin = false
	return l
}

func (l *DockerLogin) Task() Task {
	return l.Run
}

func (l *DockerLogin) Run(ctx context.Context) error {
	args := []string{"login"}
	if len(l.username) > 0 {
		args = append(args, "-u", l.username)
	}
	if !l.readStdin {
		if len(l.password) == 0 {
			return errors.New("missing password for login")
		}
		args = append(args, "-p", string(l.password))
	}
	args = append(args, l.host)
	return l.ref.Command(args...).Run(ctx)
}

func (d *DockerRef) Pull(imageAndTag string) Task {
	anyBlankPanic(strmap{"imageAndTag": &imageAndTag})
	return d.Command("pull", imageAndTag)
}

func (d *DockerRef) Tag(currentTag, newTag string) Task {
	anyBlankPanic(strmap{"currentTag": &currentTag, "newTag": &newTag})
	return d.Command("tag", currentTag, newTag)
}

func (d *DockerRef) Push(imageAndTag string) Task {
	anyBlankPanic(strmap{"imageAndTag": &imageAndTag})
	return d.Command("push", imageAndTag)
}
