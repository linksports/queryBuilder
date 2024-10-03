package queryBuilder

import (
	"strings"
)

func Trim(s string) string {
	return strings.Replace(strings.Replace(s, "\n", "", -1), "\t", "", -1)
}
