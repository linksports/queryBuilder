package queryBuilder_test

import (
	"testing"

	queryBuilder "github.com/linksports/esQueryBuilder"
	"github.com/stretchr/testify/assert"
)

func TestTrim(t *testing.T) {
	assert.Equal(t, queryBuilder.Trim(`
	a
	b	c
	`), "abc")
}
