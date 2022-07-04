package script

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestName(t *testing.T) {
	setCaller("script/builder_test.go")
	NewCQC()
	assert.Equal(t, 1, 2)
}
