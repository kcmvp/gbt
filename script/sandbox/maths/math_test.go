package maths

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormalNumberAdd(t *testing.T) {
	assert.Equal(t, 3, add(1, 2))
}
