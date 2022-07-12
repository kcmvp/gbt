package string

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormalStringAdd(t *testing.T) {
	assert.Equal(t, "abc", add("ab", "c"))
}
