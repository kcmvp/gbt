package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type MageTestSuit struct {
	BasicTestSuite
}

func TestMageSuit(t *testing.T) {
	suite.Run(t, new(MageTestSuit))
}

func (s *MageTestSuit) TestNotInRoot() {
	os.Chdir("cmd")
	cmd := exec.Command("go", "run", "github.com/kcmvp/gbt/gbt-cli", "mage")
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
			fmt.Println("msg", msg)
			assert.True(s.T(), strings.Contains(msg, not_in_root.Error()))
			verified = true
		}
	}
	assert.True(s.T(), verified)
}

//func (s *MageTestSuit) TestMageCommands() {
//	tests := []command{
//		{
//			"mage",
//			"go",
//			[]string{"run", "github.com/kcmvp/gbt/gbt-cli", "mage"},
//		},
//	}
//	for _, tt := range tests {
//		s.T().Run(tt.name, func(t *testing.T) {
//			cmd := exec.Command(tt.cmd, tt.args...)
//			stderr := bytes.NewBuffer(nil)
//			stdout := bytes.NewBuffer(nil)
//			cmd.Stderr = stderr
//			cmd.Stdout = stdout
//			require.NoError(s.T(), cmd.Run())
//		})
//	}
//}
