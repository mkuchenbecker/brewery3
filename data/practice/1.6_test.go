package practice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test1_6(t *testing.T) {
	t.Run("book example", func(t *testing.T) {
		t.Parallel()
		s, err := Compress("aabcccccaaa")
		assert.NoError(t, err)
		assert.Equal(t, "a2b1c5a3", s)
	})
	t.Run("no compression", func(t *testing.T) {
		t.Parallel()
		s, err := Compress("AaaBbbCcc")
		assert.NoError(t, err)
		assert.Equal(t, "AaaBbbCcc", s)
	})
}
