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

func TestReturnFromSubRoutine(t *testing.T) {
	c := newCpu()
	c.Reset()
	c.stack[c.sp] = 0x30
	c.sp = c.sp + 1
	c.memory[0x200] = 0x00
	c.memory[0x201] = 0xEE
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x30), c.pc)
	assert.Equal(t, uint16(0x00), c.sp)
}

func TestJumpToNNN(t *testing.T) {
	c := newCpu()
	c.Reset()
	c.memory[0x200] = 0x10
	c.memory[0x201] = 0xFF
	c.RunCpuCycle()
	assert.Equal(t, uint16(0xFF), c.pc)
}

func TestCallAddr(t *testing.T) {
	c := newCpu()
	c.Reset()
	c.memory[0x200] = 0x26
	c.memory[0x201] = 0x93
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x01), c.sp)
	assert.Equal(t, uint16(0x200), c.stack[c.sp-1])
	assert.Equal(t, uint16(0x693), c.pc)
}

func TestSkipIfVxIsKKIsTrue(t *testing.T) {
	c := newCpu()
	c.Reset()
	c.memory[0x200] = 0x3B
	c.memory[0x201] = 0x54
	c.V[0xB] = 0x54
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x204), c.pc)
}

func TestSkipIfVxIsKKIsFalse(t *testing.T) {
	c := newCpu()
	c.Reset()
	c.memory[0x200] = 0x31
	c.memory[0x201] = 0x54
	c.V[1] = 0x95
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x202), c.pc)
}

func TestSkipIfVxIsNotKKIsTrue(t *testing.T) {
	c := newCpu()
	c.Reset()
	c.memory[0x200] = 0x4B
	c.memory[0x201] = 0x54
	c.V[0xB] = 0x54
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x202), c.pc)
}

func TestSkipIfVxIsNotKKIsFalse(t *testing.T) {
	c := newCpu()
	c.Reset()
	c.memory[0x200] = 0x41
	c.memory[0x201] = 0x54
	c.V[0x1] = 0x95
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x204), c.pc)
}

func TestSkipIfVxIsVyIsTrue(t *testing.T) {
	c := newCpu()
	c.Reset()
	c.memory[0x200] = 0x53
	c.memory[0x201] = 0xB0
	c.V[0x3] = 0x96
	c.V[0xB] = 0x96
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x204), c.pc)
}

func TestSkipIfVxIsVyIsFalse(t *testing.T) {
	c := newCpu()
	c.Reset()
	c.memory[0x200] = 0x53
	c.memory[0x201] = 0xB0
	c.V[0x3] = 0x94
	c.V[0xB] = 0x96
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x202), c.pc)
}

func TestSetVxToKK(t *testing.T) {
	c := newCpu()
	c.Reset()
	c.memory[0x200] = 0x63
	c.memory[0x201] = 0x94
	c.RunCpuCycle()
	assert.Equal(t, byte(0x94), c.V[0x3])
}
