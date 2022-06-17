package xast

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var dir = []string{
	"/home/path/pkg",
	"/home/path/pkg1",
	"/home/path/pkg2",
	"/home/path1/go",
}

func TestExcludeKeys(t *testing.T) {
	l := len(dir)
	assert.Equal(t, len(excludeKeys(dir, "")), l)
	assert.Equal(t, len(excludeKeys(dir, "path/")), l-3)
	assert.Equal(t, len(excludeKeys(dir, "path1/")), l-1)
	assert.Equal(t, len(excludeKeys(dir, " path1/ ")), l-1)
	assert.Equal(t, len(excludeKeys(dir, "pkg1,pkg2")), l-2)
}
