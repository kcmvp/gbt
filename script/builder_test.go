package script

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/thedevsaddam/gojsonq/v2"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type ScriptTestSuite struct {
	suite.Suite
	args []string
}

func (suit *ScriptTestSuite) SetupTest() {
	suit.args = []string{"-run", "TestNormal*"}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestScriptTestSuite(t *testing.T) {
	suite.Run(t, new(ScriptTestSuite))
}

func (suit *ScriptTestSuite) TestCleanTestFlow() {
	cqc := NewCQC()
	cqc.Clean()

	_, err := os.Stat("./target/line_coverage.data")
	require.Error(suit.T(), err, "should no coverage.data")
	_, err = os.Stat("./target/test.data")
	require.Error(suit.T(), err, "should no test.data")
	_, err = os.Stat("./target/test.json")
	require.Error(suit.T(), err, "should no test.json")
	cqc.Test(suit.args...)
	_, err = os.Stat("./target/line_coverage.data")
	require.NoError(suit.T(), err, "should have coverage.data")
	_, err = os.Stat("./target/test.data")
	require.NoError(suit.T(), err, "should have test.data")
	_, err = os.Stat("./target/test.json")
	require.NoError(suit.T(), err, "should have test.json")
}

func (suit *ScriptTestSuite) TestJsonDataIncludeDummy() {
	cqc := NewCQC()
	cqc.Clean().Test(suit.args...)
	data, err := os.ReadFile("./target/test.json")
	require.NoError(suit.T(), err, "test.json should be generated")
	assert.NotEmpty(suit.T(), data)
	jq := gojsonq.New(gojsonq.WithSeparator("#")).FromString(string(data))
	node := jq.Find("Packages#github.com/kcmvp/gbt/script/sandbox/dummy")
	assert.NotNil(suit.T(), node)
}

func (suit *ScriptTestSuite) TestBuildWithDefault() {
	cqc := NewCQC()
	cqc.Clean().Build()
	found := false
	filepath.Walk(cqc.root, func(path string, info fs.FileInfo, err error) error {
		found = strings.EqualFold("main", info.Name()) || strings.EqualFold("main.exe", info.Name())
		return nil
	})
	require.True(suit.T(), found)
}
func (suit *ScriptTestSuite) TestBuildWithSpecifiedFiles() {
	cqc := NewCQC()
	cqc.Clean().Build("nothing.go")
	found := false
	filepath.Walk(cqc.root, func(path string, info fs.FileInfo, err error) error {
		found = strings.EqualFold("nothing", info.Name()) || strings.EqualFold("nothing.exe", info.Name())
		return nil
	})
	require.True(suit.T(), found)
}

func (suit *ScriptTestSuite) TestBuildWithMultipleFiles() {
	cqc := NewCQC()
	cqc.Clean().Build("nothing.go", "main.go")
	nothing := false
	main := false
	filepath.Walk(cqc.root, func(path string, info fs.FileInfo, err error) error {
		if !nothing {
			nothing = strings.EqualFold("nothing", info.Name()) || strings.EqualFold("nothing.exe", info.Name())
		} else if !main {
			main = strings.EqualFold("main", info.Name()) || strings.EqualFold("main.exe", info.Name())
		}
		return nil
	})
	require.True(suit.T(), nothing)
	require.True(suit.T(), main)
}

func (suit *ScriptTestSuite) TestJsonDataUncovered() {
	//pkgName := "github.com/kcmvp/gbt/script"
	cqc := NewCQC()
	cqc.Clean()
	cqc.Test(suit.args...)
	data, err := os.ReadFile("./target/test.json")
	require.NoError(suit.T(), err, "test.json should be generated")
	assert.NotEmpty(suit.T(), data)
}

func (suit *ScriptTestSuite) TestSecScan() {
	cqc := NewCQC()
	cqc.Clean().SecScan()
	_, err := os.Stat(filepath.Join(cqc.target, "security.json"))
	require.NoError(suit.T(), err)
}

func (suit *ScriptTestSuite) TestStaticScan() {
	cqc := NewCQC()
	cqc.Clean().StaticScan()
	_, err := os.Stat(filepath.Join(cqc.target, "static.json"))
	require.NoError(suit.T(), err)
}
