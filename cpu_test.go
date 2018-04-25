package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCpu(t *testing.T) {
	c := newCpu()
	assert.NotNil(t, c)
}

func TestLoadProgram(t *testing.T) {
	c := newCpu()
	n := c.LoadProgram("roms/PONG")
	assert.Equal(t, 246, n, "246 bytes should be read as the game is 246 bytes long")
}
