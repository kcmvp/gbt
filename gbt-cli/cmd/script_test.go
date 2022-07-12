package cmd_test

import (
	"fmt"
	"github.com/kcmvp/gbt/gbt-cli/cmd"
	"github.com/kcmvp/gbt/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

type ScriptTestSuit struct {
	suite.Suite
}

var builderMsg = "Generate build script for the project"
var configMsg = "Generate system configuration files application.yml"
var notInRootMsg = "Error: please run the command from project root"

func (s *ScriptTestSuit) SetupTest() {
	err := os.Chdir(script.ProjectRoot())
	require.NoError(s.T(), err)
}

func (s *ScriptTestSuit) TearDownSuite() {
	err := os.Remove("application.yml")
	require.NoError(s.T(), err)
}

func (s *ScriptTestSuit) TestRootUsage() {
	out, err := exec.Command("go", "run", filepath.Join(script.ProjectRoot(), "main.go")).CombinedOutput()
	require.NoError(s.T(), err)
	assert.Contains(s.T(), string(out), "Usage:")
	assert.Contains(s.T(), string(out), builderMsg)
	assert.Contains(s.T(), string(out), configMsg)
	assert.NotContains(s.T(), string(out), notInRootMsg)
}

func (s *ScriptTestSuit) TestNotInRootDir() {
	os.Chdir(filepath.Dir(script.ProjectRoot()))
	out, err := exec.Command("go", "run", filepath.Join(script.ProjectRoot(), "main.go")).CombinedOutput()
	require.Error(s.T(), err)
	fmt.Println(string(out))
	assert.Contains(s.T(), string(out), notInRootMsg)
}

func (s *ScriptTestSuit) TestConfig() {
	out, err := exec.Command("go", "run", filepath.Join(script.ProjectRoot(), "main.go"), "config").CombinedOutput()
	require.NoError(s.T(), err)
	pwd, _ := os.Getwd()
	assert.Contains(s.T(), string(out), fmt.Sprintf("create files: %s/application.yml successfully", pwd))
}

func (s *ScriptTestSuit) TestWithExistingConfig() {
	_, err := os.Create(cmd.Application)
	require.NoError(s.T(), err)
	out, err := exec.Command("go", "run", filepath.Join(script.ProjectRoot(), "main.go"), "config").CombinedOutput()
	fmt.Println(string(out))
	require.NoError(s.T(), err)
	pwd, _ := os.Getwd()
	assert.Contains(s.T(), string(out), fmt.Sprintf("%v/application.yml exists", pwd))
}

func (s ScriptTestSuit) TestBuilder() {
	out, err := exec.Command("go", "run", filepath.Join(script.ProjectRoot(), "main.go"), "defaultBuilder").CombinedOutput()
	fmt.Println(string(out))
	require.NoError(s.T(), err)
	pwd, _ := os.Getwd()
	assert.Contains(s.T(), string(out), fmt.Sprintf("%v/application.yml exists", pwd))
}

func TestScriptTestSuit(t *testing.T) {
	suite.Run(t, new(ScriptTestSuit))
}
