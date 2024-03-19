package modmake_docker

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocker(t *testing.T) {
	ctx := context.Background()
	err := Docker().Dry().Command("do", "command").Run(ctx)
	assert.Error(t, err)
	assert.Equal(t, "dry run: docker do command", err.Error())
}
