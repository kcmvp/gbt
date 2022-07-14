package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type CmdTestSuit struct {
	suite.Suite
	root string
}

func (s *CmdTestSuit) SetupSuite() {
	os.Remove(Application)
	os.RemoveAll(builderDir)
	pwd, _ := os.Getwd()
	s.root = filepath.Dir(pwd)
}

func (s *CmdTestSuit) SetupTest() {
	os.Chdir(s.root)
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

func (s *CmdTestSuit) TestBuilderCmd() {
	cmd := NewRootCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"builder", "-u"})
	err := cmd.Execute()
	assert.True(s.T(), bFlag.update)
	require.NoError(s.T(), err)
}
