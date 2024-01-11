package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	src := `
class Main {
	static boolean test;
	field int x, y;
}
`
	engine := NewCompilationEngine(strings.NewReader(src))
	cls := engine.compileClass()

	assert.Equal(t, cls.Name, "Main")
	assert.Equal(t, cls.Vars[0].StorageClass, "static")
	assert.Equal(t, cls.Vars[0].Type, "boolean")
	assert.Equal(t, cls.Vars[0].Names[0], "test")
}
