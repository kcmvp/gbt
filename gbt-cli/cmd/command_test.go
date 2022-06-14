package cmd

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type CommandTestSuit struct {
	BasicTestSuite
}

func (s *CommandTestSuit) TestNotInRoot() {
	os.Chdir("cmd")
	cmd := exec.Command("go", "run", "github.com/kcmvp/gbt/gbt-cli", "init")
	stderr := bytes.NewBuffer(nil)
	stdout := bytes.NewBuffer(nil)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	assert.Error(s.T(), cmd.Run())
	verified := false
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		msg := scanner.Text()
		if strings.HasPrefix(msg, "Error:") {
			assert.True(s.T(), strings.Contains(msg, not_in_root.Error()))
			verified = true
		}
	}
	assert.True(s.T(), verified)
}

func (s *CommandTestSuit) TestRoot() {

	assert.Equal(s.T(), 123, 123, "they should be equal")
	cmd := exec.Command("go", "run", "github.com/kcmvp/gbt/gbt-cli")
	stderr := bytes.NewBuffer(nil)
	stdout := bytes.NewBuffer(nil)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	assert.NoError(s.T(), cmd.Run())
}

func (s *CommandTestSuit) TestInit() {
	cmd := exec.Command("go", "run", "github.com/kcmvp/gbt/gbt-cli", "init")
	stderr := bytes.NewBuffer(nil)
	stdout := bytes.NewBuffer(nil)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	require.NoError(s.T(), cmd.Run())
}

func TestCommandSuit(t *testing.T) {
	suite.Run(t, new(CommandTestSuit))
}
