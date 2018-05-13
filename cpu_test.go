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
	for i := 0x50; i < 0x200; i++ {
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
	c.memory[0x200] = 0x10
	c.memory[0x201] = 0xFF
	c.RunCpuCycle()
	assert.Equal(t, uint16(0xFF), c.pc)
}

func TestCallAddr(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x26
	c.memory[0x201] = 0x93
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x01), c.sp)
	assert.Equal(t, uint16(0x200), c.stack[c.sp-1])
	assert.Equal(t, uint16(0x693), c.pc)
}

func TestSkipIfVxIsKKIsTrue(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x3B
	c.memory[0x201] = 0x54
	c.V[0xB] = 0x54
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x204), c.pc)
}

func TestSkipIfVxIsKKIsFalse(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x31
	c.memory[0x201] = 0x54
	c.V[1] = 0x95
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x202), c.pc)
}

func TestSkipIfVxIsNotKKIsTrue(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x4B
	c.memory[0x201] = 0x54
	c.V[0xB] = 0x54
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x202), c.pc)
}

func TestSkipIfVxIsNotKKIsFalse(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x41
	c.memory[0x201] = 0x54
	c.V[0x1] = 0x95
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x204), c.pc)
}

func TestSkipIfVxIsVyIsTrue(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x53
	c.memory[0x201] = 0xB0
	c.V[0x3] = 0x96
	c.V[0xB] = 0x96
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x204), c.pc)
}

func TestSkipIfVxIsVyIsFalse(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x53
	c.memory[0x201] = 0xB0
	c.V[0x3] = 0x94
	c.V[0xB] = 0x96
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x202), c.pc)
}

func TestSetVxToKK(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x63
	c.memory[0x201] = 0x94
	c.RunCpuCycle()
	assert.Equal(t, byte(0x94), c.V[0x3])
}

func TestAddByteToVx(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x7C
	c.memory[0x201] = 0xFE
	c.V[0xC] = 0x1
	c.RunCpuCycle()
	assert.Equal(t, byte(0xFF), c.V[0xC])
}

func TestAddByteToVxOverflow(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x7C
	c.memory[0x201] = 0xFF
	c.V[0xC] = 0x90
	c.RunCpuCycle()
	assert.Equal(t, byte(0x8f), c.V[0xC])
}

func TestVxAssignVy(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x8A
	c.memory[0x201] = 0xB0
	c.V[0xB] = 0x90
	c.RunCpuCycle()
	assert.Equal(t, byte(0x90), c.V[0xA])
}

func TestVxOrVy(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x8A
	c.memory[0x201] = 0xC1
	c.V[0xA] = 0x11
	c.V[0xC] = 0x43
	c.RunCpuCycle()
	assert.Equal(t, byte(0x53), c.V[0xA])
}

func TestVxAndVy(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x8A
	c.memory[0x201] = 0xC2
	c.V[0xA] = 0x34
	c.V[0xC] = 0xD3
	c.RunCpuCycle()
	assert.Equal(t, byte(0x10), c.V[0xA])
}

func TestVxXorVy(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x8A
	c.memory[0x201] = 0xD3
	c.V[0xA] = 0xA3
	c.V[0xD] = 0x3A
	c.RunCpuCycle()
	assert.Equal(t, byte(0x99), c.V[0xA])
}

func TestAddVxVyNoOverflow(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x8B
	c.memory[0x201] = 0xE4
	c.V[0xB] = 0x11
	c.V[0xE] = 0x53
	c.RunCpuCycle()
	assert.Equal(t, byte(0x64), c.V[0xB])
	assert.Equal(t, byte(0x0), c.V[0xF])
}

func TestAddVxVyOverflow(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x8B
	c.memory[0x201] = 0xF4
	c.V[0xB] = 0xAA
	c.V[0xF] = 0xFF
	c.RunCpuCycle()
	assert.Equal(t, byte(0xA9), c.V[0xB])
	assert.Equal(t, byte(0x1), c.V[0xF])
}

func TestSubVxVyNoBorrow(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x87
	c.memory[0x201] = 0x65
	c.V[0x7] = 0x99
	c.V[0x6] = 0x33
	c.RunCpuCycle()
	assert.Equal(t, byte(0x66), c.V[0x7])
	assert.Equal(t, byte(0x1), c.V[0xF])
}

func TestSubVxVyBorrow(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x87
	c.memory[0x201] = 0x65
	c.V[0x7] = 0x98
	c.V[0x6] = 0xAA
	c.RunCpuCycle()
	assert.Equal(t, byte(0xEE), c.V[0x7])
	assert.Equal(t, byte(0x0), c.V[0xF])
}

func TestShrVxLsbIsOne(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x8B
	c.memory[0x201] = 0xC6
	c.V[0xB] = 0x99
	c.RunCpuCycle()
	assert.Equal(t, byte(0x1), c.V[0xF])
	assert.Equal(t, byte(0x4C), c.V[0xB])
}

func TestShrVxLsbIsNotOne(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x8B
	c.memory[0x201] = 0xC6
	c.V[0xB] = 0x98
	c.RunCpuCycle()
	assert.Equal(t, byte(0x0), c.V[0xF])
	assert.Equal(t, byte(0x4C), c.V[0xB])
}

func TestVySubVxNoBorrow(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x8B
	c.memory[0x201] = 0xC7
	c.V[0xB] = 0x89
	c.V[0xC] = 0x95
	c.RunCpuCycle()
	assert.Equal(t, byte(0x1), c.V[0xF])
	assert.Equal(t, byte(0xC), c.V[0xB])
}

func TestVySubVxBorrow(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x8B
	c.memory[0x201] = 0xC7
	c.V[0xB] = 0x01
	c.V[0xC] = 0x00
	c.RunCpuCycle()
	assert.Equal(t, byte(0x0), c.V[0xF])
	assert.Equal(t, byte(0xFF), c.V[0xB])
}

func TestShlVxMsbIsOne(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x8A
	c.memory[0x201] = 0xCE
	c.V[0xA] = 0xAB
	c.RunCpuCycle()
	assert.Equal(t, byte(0x01), c.V[0xF])
	assert.Equal(t, byte(0x56), c.V[0xA])
}

func TestShlVxMsbIsNotOne(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x8A
	c.memory[0x201] = 0xCE
	c.V[0xA] = 0x3B
	c.RunCpuCycle()
	assert.Equal(t, byte(0x00), c.V[0xF])
	assert.Equal(t, byte(0x76), c.V[0xA])
}

func TestSneVxVyNotEqual(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x9B
	c.memory[0x201] = 0xD0
	c.V[0xB] = 0xDD
	c.V[0xD] = 0xCC
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x204), c.pc)
}

func TestSneVxVyEqual(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x9B
	c.memory[0x201] = 0xD0
	c.V[0xB] = 0xDD
	c.V[0xD] = 0xDD
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x202), c.pc)
}

func TestLoadAddress(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0xAB
	c.memory[0x201] = 0x34
	c.RunCpuCycle()
	assert.Equal(t, uint16(0xB34), c.I)
}

func TestJumpToLocationPlusV0(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0xB1
	c.memory[0x201] = 0x94
	c.V[0x0] = 0x6
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x19A), c.pc)
}

func TestSetVxToRandomNumberAndKK(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0xCA
	c.memory[0x201] = 0xFF
	c.RunCpuCycle()
}

func TestLoadFontSet(t *testing.T) {
	c := newCpu()
	c.LoadFontSet()
	for i := 0x00; i < 0x50; i++ {
		assert.Equal(t, fontset[i], c.memory[i])
	}
}

func TestClearDisplay(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0x00
	c.memory[0x201] = 0xE0
	c.display[0][0] = 0x1
	c.display[9][23] = 0x1
	c.RunCpuCycle()
	for x := 0x00; x < 0x20; x++ {
		for y := 0x00; y < 0x40; y++ {
			assert.Equal(t, byte(0), c.display[x][y])
		}
	}
}

func TestDXYNNoWrapAroundNoCollision(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0xD3
	c.memory[0x201] = 0xD2
	c.I = 0x300
	c.V[0x3] = 0
	c.V[0xD] = 0
	c.memory[0x300] = 0x11
	c.memory[0x301] = 0x88
	c.ClearDisplay()
	c.RunCpuCycle()
	assert.Equal(t, byte(0x00), c.V[0xF])
	assert.Equal(t, byte(0x01), c.display[0][3])
	assert.Equal(t, byte(0x01), c.display[0][7])
	assert.Equal(t, byte(0x01), c.display[1][0])
	assert.Equal(t, byte(0x01), c.display[1][4])
}

func TestDXYNNoWrapAroundYesCollision(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0xD3
	c.memory[0x201] = 0xD2
	c.I = 0x300
	c.memory[0x300] = 0x11
	c.memory[0x301] = 0x88
	c.V[0x3] = 0
	c.V[0xD] = 0
	c.ClearDisplay()
	c.display[0][3] = 0x01
	c.RunCpuCycle()
	assert.Equal(t, byte(0x01), c.V[0xF])
	assert.Equal(t, byte(0x00), c.display[0][3])
	assert.Equal(t, byte(0x01), c.display[0][7])
	assert.Equal(t, byte(0x01), c.display[1][0])
	assert.Equal(t, byte(0x01), c.display[1][4])
}

func TestDXYNWithWrapAroundNoCollision(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0xD3
	c.memory[0x201] = 0xD2
	c.I = 0x300
	c.memory[0x300] = 0x11
	c.memory[0x301] = 0x88
	c.V[0x3] = 0x3F
	c.V[0xD] = 0x1F
	c.ClearDisplay()
	c.RunCpuCycle()
	assert.Equal(t, byte(0x00), c.V[0xF])
	assert.Equal(t, byte(0x01), c.display[31][2])
	assert.Equal(t, byte(0x01), c.display[31][6])
	assert.Equal(t, byte(0x01), c.display[0][3])
	assert.Equal(t, byte(0x01), c.display[0][63])
}

func TestDXYNWithWrapAroundYesCollision(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0xD3
	c.memory[0x201] = 0xD2
	c.I = 0x300
	c.memory[0x300] = 0x11
	c.memory[0x301] = 0x88
	c.V[0x3] = 0x3F
	c.V[0xD] = 0x1F
	c.ClearDisplay()
	c.display[31][2] = 0x01
	c.RunCpuCycle()
	assert.Equal(t, byte(0x01), c.V[0xF])
	assert.Equal(t, byte(0x00), c.display[31][2])
	assert.Equal(t, byte(0x01), c.display[31][6])
	assert.Equal(t, byte(0x01), c.display[0][3])
	assert.Equal(t, byte(0x01), c.display[0][63])
}

func TestSetDelayTimer(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0xFD
	c.memory[0x201] = 0x15
	c.V[0xD] = 0x33
	c.RunCpuCycle()
	assert.Equal(t, byte(0x33), c.delayTimer)
}

func TestSetVxToDelayTimer(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0xFD
	c.memory[0x201] = 0x07
	c.delayTimer = 0x44
	c.RunCpuCycle()
	assert.Equal(t, byte(0x44), c.V[0xD])
}

func TestSetSoundTimer(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0xFD
	c.memory[0x201] = 0x18
	c.V[0xD] = 0x99
	c.RunCpuCycle()
	assert.Equal(t, byte(0x99), c.soundTimer)
}

func TestSetIToIPlusVx(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0xFD
	c.memory[0x201] = 0x1E
	c.I = 0x32
	c.V[0xD] = 0x33
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x65), c.I)
}

func TestSetIToLocationOfDigit(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0xFD
	c.memory[0x201] = 0x29
	c.V[0xD] = 0x7
	c.RunCpuCycle()
	assert.Equal(t, uint16(0x23), c.I)
}

func TestSetBCDRepresentation(t *testing.T) {
	c := newCpu()
	c.memory[0x200] = 0xFB
	c.memory[0x201] = 0x33
	c.I = 0x90
	c.V[0xB] = 0x7B
	c.RunCpuCycle()
	assert.Equal(t, byte(0x1), c.memory[c.I])
	assert.Equal(t, byte(0x2), c.memory[c.I+1])
	assert.Equal(t, byte(0x3), c.memory[c.I+2])
}
