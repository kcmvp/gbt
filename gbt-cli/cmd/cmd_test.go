package cmd

import (
	"bytes"
	"github.com/kcmvp/gbt/script"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type CmdTestSuit struct {
	suite.Suite
}

func (s *CmdTestSuit) SetupTest() {
	err := os.Chdir(script.ProjectRoot())
	require.NoError(s.T(), err)
	os.Remove(Application)

}

func TestCmdTestSuit(t *testing.T) {
	suite.Run(t, new(CmdTestSuit))
}

func (s *CmdTestSuit) TestConfigCmd() {
	cmd := NewRootCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"config"})
	err := cmd.Execute()
	require.NoError(s.T(), err)
	_, err = os.Stat(Application)
	require.NoError(s.T(), err)
}
