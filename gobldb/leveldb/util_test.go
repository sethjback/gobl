package leveldb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntersectSlice(t *testing.T) {
	assert := assert.New(t)

	s1 := []string{"a", "b", "c", "d", "e", "f"}
	s2 := []string{"d", "e", "f", "g", "h", "i", "k"}

	intersect := intersectSlice(s1, s2)
	if assert.Len(intersect, 3) {
		assert.Contains(intersect, "d")
		assert.Contains(intersect, "e")
		assert.Contains(intersect, "f")
	}
}

func TestStringInSlice(t *testing.T) {
	assert := assert.New(t)

	s1 := []string{"a", "b", "c", "d", "e", "f"}
	assert.True(stringInSlice(s1, "c"))
	assert.False(stringInSlice(s1, "j"))
}
