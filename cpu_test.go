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
	for i := 0; i < 512; i++ {
		assert.Equal(t, uint8(0), c.memory[i], "Should be 0 as first 512 is where emulator resides")
	}
}

func TestLoadProgramFailsWithWrongFile(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	c := newCpu()
	c.LoadProgram("roms/FOO")
}

func TestReset(t *testing.T) {
	c := newCpu()
	c.LoadProgram("roms/PONG")
	c.I = 42
	c.Reset()
	f := newCpu()
	assert.Equal(t, f, c, "After reset it should be same as new")
}

func TestRunCpuCycle(t *testing.T) {
	c := newCpu()
	c.LoadProgram("roms/PONG")
	c.RunCpuCycle()
}

func TestReturnFromSubRoutine(t *testing.T) {
	c := newCpu()
	c.Reset()
	c.stack[c.sp] = 0x30
	c.pc = 0x200
	c.sp = c.sp + 1
	c.memory[0x200] = 0x00
	c.memory[0x201] = 0xEE
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x30), c.pc)
	assert.Equal(t, uint16(0x00), c.sp)
}
