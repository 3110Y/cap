package command_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/3110Y/cap/internal/presentation/command"
)

func TestVersionCommand_OutputsVersion(t *testing.T) {
	cmd := command.NewVersionCommand()

	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()

	require.NoError(t, err)
	// Version variable is set to "0.1.0-dev" by default.
	assert.Contains(t, buf.String(), command.Version)
}

func TestVersionCommand_ContainsCap(t *testing.T) {
	cmd := command.NewVersionCommand()

	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "cap")
}
