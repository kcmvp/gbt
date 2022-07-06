package script

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type ScriptTestSuite struct {
	suite.Suite
	args []string
}

func (suit *ScriptTestSuite) SetupTest() {
	suit.args = []string{"-run", "TestNormal.*"}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestScriptTestSuite(t *testing.T) {
	suite.Run(t, new(ScriptTestSuite))
}

func (suit *ScriptTestSuite) TestCleanTestFlow() {
	cqc, _ := NewCQC()
	cqc.Clean()

	_, err := os.Stat("./target/coverage.data")
	require.Error(suit.T(), err, "should no coverage.data")
	_, err = os.Stat("./target/test.data")
	require.Error(suit.T(), err, "should no test.data")
	_, err = os.Stat("./target/test.json")
	require.Error(suit.T(), err, "should no test.json")
	cqc.Test(suit.args...)
	_, err = os.Stat("./target/coverage.data")
	require.NoError(suit.T(), err, "should have coverage.data")
	_, err = os.Stat("./target/test.data")
	require.NoError(suit.T(), err, "should have test.data")
	_, err = os.Stat("./target/test.json")
	require.NoError(suit.T(), err, "should have test.json")
}

func (suit *ScriptTestSuite) TestJsonDataIncludeDummy() {
	cqc, _ := NewCQC()
	cqc.Clean()
	cqc.Test(suit.args...)

}
