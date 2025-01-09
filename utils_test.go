package queryBuilder_test

import (
	"testing"

	"github.com/linksports/queryBuilder"
	"github.com/stretchr/testify/assert"
)

func TestTrim(t *testing.T) {
	assert.Equal(t, queryBuilder.Trim(`
	a
	b	c
	`), "abc")
}
