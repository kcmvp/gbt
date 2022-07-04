package sandbox

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	sum := Add(1, 2)
	assert.Equal(t, 3, sum)
}
func TestMin(t *testing.T) {
	sum := Minus(2, 1)
	assert.Equal(t, 2, sum)
}
