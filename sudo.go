//go:build !windows

package modmake_docker

import (
	"os/user"
	"strings"
	"sync"
)

var (
	groupsOnce sync.Once
	addSudo    bool
)

func (d *DockerRef) sudoPrefix() string {
	groupsOnce.Do(func() {
		userDetails, err := user.Current()
		if err != nil {
			panicf("failed to get current user details: %v", err)
		}
		groups, err := userDetails.GroupIds()
		if err != nil {
			panicf("failed to get group IDs: %v", err)
		}
		hasDocker := false
		for _, group := range groups {
			if strings.TrimSpace(group) == "docker" {
				hasDocker = true
				break
			}
		}
		if !hasDocker {
			addSudo = true
		}
	})
	if addSudo {
		return "sudo"
	}
	return ""
}
