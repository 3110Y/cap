package command_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/3110Y/cap/internal/presentation/command"
)

func TestPingCommand_OutputsPong(t *testing.T) {
	cmd := command.NewPingCommand()

	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "pong")
}
