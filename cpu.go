package main

import (
	"io"
	"math/rand"
	"os"
	"time"
)

const (
	height = byte(0x20)
	width  = byte(0x40)
)

type cpu struct {
	pc         uint16              // program counter
	memory     [4096]byte          // 4k memory
	stack      [16]uint16          // 16 level stack
	sp         uint16              // stack pointer
	V          [16]byte            // 16 registers
	I          uint16              // The address register
	delayTimer byte                // The delay timer counts down at 60hz
	soundTimer byte                //sound timer counts down at 60hz
	display    [height][width]byte // 2d array representing 64x32 grid
	keys       [16]byte            // state of the keys
}

var fontset = [...]byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

func newCpu() cpu {
	c := cpu{pc: 0x200}
	c.LoadFontSet()
	return c
}

func (c *cpu) LoadFontSet() {
	for i := 0x00; i < 0x50; i++ {
		c.memory[i] = fontset[i]
	}
}

func (c *cpu) ClearDisplay() {
	for x := 0x00; x < 0x20; x++ {
		for y := 0x00; y < 0x40; y++ {
			c.display[x][y] = 0x00
		}
	}
}

func (c *cpu) LoadProgram(rom string) int {
	f, err := os.Open(rom)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	memory := make([]byte, 3584)
	n, err := f.Read(memory)
	if err != nil {
		if err != io.EOF {
			panic(err)
		}
	}
	for index, b := range memory {
		c.memory[index+0x200] = b
	}
	return n
}

func (c *cpu) Reset() {
	c.pc = 0x200
	c.delayTimer = 0
	c.soundTimer = 0
	c.I = 0
	c.sp = 0
	for i := 0; i < len(c.memory); i++ {
		c.memory[i] = 0
	}
	for i := 0; i < len(c.stack); i++ {
		c.stack[i] = 0
	}
	for i := 0; i < len(c.V); i++ {
		c.V[i] = 0
	}
	for i := 0; i < len(c.keys); i++ {
		c.keys[i] = 0
	}
	c.LoadFontSet()
	c.ClearDisplay()
}

func (c *cpu) Run() {
	c.RunCpuCycle()

	if c.delayTimer > 0 {
		c.delayTimer = c.delayTimer - 1
	}

	if c.soundTimer > 0 {
		c.soundTimer = c.soundTimer - 1
	}
}

func (c *cpu) RunCpuCycle() {
	opcode := uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
	c.pc = c.pc + 2
	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode & 0x000F {
		case 0x0000:
			c.ClearDisplay()
		case 0x000E:
			c.pc = c.stack[c.sp-1]
			c.sp = c.sp - 1
		}
	case 0x1000:
		c.pc = opcode & 0x0FFF
	case 0x2000:
		c.stack[c.sp] = c.pc
		c.sp = c.sp + 1
		c.pc = opcode & 0x0FFF
	case 0x3000:
		compareTo := byte(opcode & 0x00FF)
		register := (opcode & 0x0F00) >> 8
		if c.V[register] == compareTo {
			c.pc = c.pc + 2
		}
	case 0x4000:
		compareTo := byte(opcode & 0x00FF)
		register := (opcode & 0x0F00) >> 8
		if c.V[register] != compareTo {
			c.pc = c.pc + 2
		}
	case 0x5000:
		registerX := (opcode & 0x0F00) >> 8
		registerY := (opcode & 0x00F0) >> 4
		if c.V[registerX] == c.V[registerY] {
			c.pc = c.pc + 2
		}
	case 0x6000:
		register := byte((opcode & 0x0F00) >> 8)
		c.V[register] = byte(opcode & 0x00FF)
	case 0x7000:
		register := byte((opcode & 0x0F00) >> 8)
		value := byte(opcode & 0x00FF)
		c.V[register] = c.V[register] + value
	case 0x8000:
		switch opcode & 0x000F {
		case 0x0000:
			registerX := (opcode & 0x0F00) >> 8
			registerY := (opcode & 0x00F0) >> 4
			c.V[registerX] = c.V[registerY]
		case 0x0001:
			registerX := (opcode & 0x0F00) >> 8
			registerY := (opcode & 0x00F0) >> 4
			c.V[registerX] = c.V[registerX] | c.V[registerY]
		case 0x0002:
			registerX := (opcode & 0x0F00) >> 8
			registerY := (opcode & 0x00F0) >> 4
			c.V[registerX] = c.V[registerX] & c.V[registerY]
		case 0x0003:
			registerX := (opcode & 0x0F00) >> 8
			registerY := (opcode & 0x00F0) >> 4
			c.V[registerX] = c.V[registerX] ^ c.V[registerY]
		case 0x0004:
			registerX := byte((opcode & 0x0F00) >> 8)
			registerY := byte((opcode & 0x00F0) >> 4)
			c.V[registerX] = c.V[registerX] + c.V[registerY]
			if uint16(c.V[registerX])+uint16(c.V[registerY]) > 0xFF {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
		case 0x0005:
			registerX := (opcode & 0x0F00) >> 8
			registerY := (opcode & 0x00F0) >> 4
			if c.V[registerX] > c.V[registerY] {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
			c.V[registerX] = c.V[registerX] - c.V[registerY]
		case 0x0006:
			registerX := (opcode & 0x0F00) >> 8
			if c.V[registerX]&0x1 == 1 {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
			c.V[registerX] = c.V[registerX] >> 1
		case 0x0007:
			registerX := (opcode & 0x0F00) >> 8
			registerY := (opcode & 0x00F0) >> 4
			if c.V[registerY] > c.V[registerX] {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
			c.V[registerX] = c.V[registerY] - c.V[registerX]
		case 0x000E:
			registerX := (opcode & 0x0F00) >> 8
			if c.V[registerX]&0x80 == 0x80 {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
			c.V[registerX] = c.V[registerX] << 1
		}
	case 0x9000:
		registerX := (opcode & 0x0F00) >> 8
		registerY := (opcode & 0x00F0) >> 4
		if c.V[registerX] != c.V[registerY] {
			c.pc = c.pc + 2
		}
	case 0xA000:
		c.I = (opcode & 0x0FFF)
	case 0xB000:
		c.pc = (opcode & 0x0FFF) + uint16(c.V[0x0])
	case 0xC000:
		registerX := (opcode & 0x0F00) >> 8
		value := byte(opcode & 0x00FF)
		rand.Seed(time.Now().Unix())
		c.V[registerX] = byte(rand.Intn(256)) & value
	case 0xD000:
		registerX := (opcode & 0x0F00) >> 8
		registerY := (opcode & 0x00F0) >> 4
		nibble := byte(opcode & 0x000F)
		x := c.V[registerX]
		y := c.V[registerY]
		c.V[0xF] = 0x00
		for i := y; i < y+nibble; i++ {
			for j := x; j < x+8; j++ {
				bit := (c.memory[c.I+uint16(i-y)] >> (7 - j + x)) & 0x01
				xIndex, yIndex := j, i
				if j >= width {
					xIndex = j - width
				}
				if i >= height {
					yIndex = i - height
				}
				if bit == 0x01 && c.display[yIndex][xIndex] == 0x01 {
					c.V[0xF] = 0x01
				}
				c.display[yIndex][xIndex] = c.display[yIndex][xIndex] ^ bit
			}
		}
	case 0xE000:
		switch opcode & 0x00FF {
		case 0x009E:
			register := (opcode & 0x0F00) >> 8
			if c.keys[c.V[register]] == 0x01 {
				c.pc = c.pc + 2
			}
		case 0x00A1:
			register := (opcode & 0x0F00) >> 8
			if c.keys[c.V[register]] == 0x00 {
				c.pc = c.pc + 2
			}
		}
	case 0xF000:
		switch opcode & 0x00FF {
		case 0x007:
			register := (opcode & 0x0F00) >> 8
			c.V[register] = c.delayTimer
		case 0x0015:
			register := (opcode & 0x0F00) >> 8
			c.delayTimer = c.V[register]
		case 0x0018:
			register := (opcode & 0x0F00) >> 8
			c.soundTimer = c.V[register]
		case 0x001E:
			register := (opcode & 0x0F00) >> 8
			c.I = c.I + uint16(c.V[register])
		case 0x0029:
			register := (opcode & 0x0F00) >> 8
			c.I = uint16(c.V[register] * 0x5)
		case 0x0033:
			register := (opcode & 0x0F00) >> 8
			number := c.V[register]
			c.memory[c.I] = (number / 100) % 10
			c.memory[c.I+1] = (number / 10) % 10
			c.memory[c.I+2] = number % 10
		case 0x0055:
			register := (opcode & 0x0F00) >> 8
			for i := uint16(0x00); i <= register; i++ {
				c.memory[c.I+i] = c.V[i]
			}
		case 0x0065:
			register := (opcode & 0x0F00) >> 8
			for i := uint16(0x00); i <= register; i++ {
				c.V[i] = c.memory[c.I+i]
			}
		}
	}
}
