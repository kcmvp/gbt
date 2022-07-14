package cmd_test

import (
	"fmt"
	"github.com/kcmvp/gbt/gbtc/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

type ScriptTestSuit struct {
	suite.Suite
	root string
}

var builderMsg = "Generate build script for the project"
var configMsg = "Generate system configuration files application.yml"
var notInRootMsg = "Error: please run the command from project root"

func (s *ScriptTestSuit) SetupSuite() {
	_, filename, _, _ := runtime.Caller(0)
	cmd := filepath.Dir(filename)
	s.root = filepath.Dir(cmd)
}

func (s *ScriptTestSuit) SetupTest() {
	os.Chdir(s.root)
	os.Remove("application.yml")
	os.RemoveAll(filepath.Join(s.root, "scripts"))
}

func (w *ScriptTestSuit) TearDownTest() {
	os.Remove("application.yml")
	os.RemoveAll("scripts")
}

func (s *ScriptTestSuit) TestRootUsage() {
	out, err := exec.Command("go", "run", filepath.Join(s.root, "main.go")).CombinedOutput()
	require.NoError(s.T(), err)
	assert.Contains(s.T(), string(out), "Usage:")
	assert.Contains(s.T(), string(out), builderMsg)
	assert.Contains(s.T(), string(out), configMsg)
	assert.NotContains(s.T(), string(out), notInRootMsg)
}

func (s *ScriptTestSuit) TestConfig() {
	out, err := exec.Command("go", "run", filepath.Join(s.root, "main.go"), "config").CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprintf("%+v", err))
	}
	require.NoError(s.T(), err)
	pwd, _ := os.Getwd()
	assert.Contains(s.T(), string(out), fmt.Sprintf("create files: %s/application.yml successfully", pwd))
}

func (s *ScriptTestSuit) TestWithExistingConfig() {
	_, err := os.Create(cmd.Application)
	require.NoError(s.T(), err)
	out, err := exec.Command("go", "run", filepath.Join(s.root, "main.go"), "config").CombinedOutput()
	fmt.Println(string(out))
	require.NoError(s.T(), err)
	pwd, _ := os.Getwd()
	assert.Contains(s.T(), string(out), fmt.Sprintf("%v/application.yml exists", pwd))
}

func TestScriptTestSuit(t *testing.T) {
	suite.Run(t, new(ScriptTestSuit))
}
