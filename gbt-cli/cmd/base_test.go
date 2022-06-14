package cmd

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type BasicTestSuite struct {
	suite.Suite
}

func (s RootTestSuit) SetupTest() {
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "../..")
	err := os.Chdir(filepath.Join(root, "gbt-cli"))
	require.NoError(s.T(), err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(BasicTestSuite))
}
