package cmd_test

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os/exec"
	"testing"
)

type RootTestSuit struct {
	BasicTestSuite
}

func (s *RootTestSuit) TestNotInRoot() {
	cmd := exec.Command("go", "run", "github.com/kcmvp/gbt/gbt-cli")
	stderr := bytes.NewBuffer(nil)
	stdout := bytes.NewBuffer(nil)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	require.NoError(s.T(), cmd.Run())
}
func TestRootSuit(t *testing.T) {
	suite.Run(t, new(RootTestSuit))
}
