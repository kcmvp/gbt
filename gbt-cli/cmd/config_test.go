package cmd_test

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os/exec"
	"testing"
)

type InitializrTestSuite struct {
	BasicTestSuite
}

func (s *InitializrTestSuite) TestInit() {
	cmd := exec.Command("go", "run", "github.com/kcmvp/gbt/gbt-cli init")
	stderr := bytes.NewBuffer(nil)
	stdout := bytes.NewBuffer(nil)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	require.NoError(s.T(), cmd.Run())
}

func TestInitializrSuite(t *testing.T) {
	suite.Run(t, new(RootTestSuit))
}
