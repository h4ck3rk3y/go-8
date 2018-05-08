package main

import (
	"fmt"
	"io"
	"os"
)

type cpu struct {
	pc         uint16     // program counter
	memory     [4096]byte // 4k memory
	stack      [16]uint16 // 16 level stack
	sp         uint16     // stack pointer
	V          [16]byte   // 16 registers
	I          uint16     // The address register
	delayTimer uint16     // The delay timer counts down at 60hz
	soundTimer uint16     //sound timer counts down at 60hz
}

func newCpu() cpu {
	c := cpu{pc: 0x200}
	return c
}

func (c *cpu) LoadProgram(rom string) int {
	f, err := os.Open(rom)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	memory := make([]byte, 3584)
	n, err := f.Read(memory)
	for index, b := range memory {
		c.memory[index+0x200] = b
	}
	if err != nil {
		if err != io.EOF {
			panic(err)
		}
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
}

func (c *cpu) RunCpuCycle() {
	opcode := uint16(c.memory[c.pc]<<8) | uint16(c.memory[c.pc+1])
	c.pc = c.pc + 2
	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode & 0x000F {
		case 0x0000:
			fmt.Println("Clear Screen not implemented")
		case 0x000E:
			c.pc = c.stack[c.sp-1]
			c.sp = c.sp - 1
		}
	case 0x1000:
		c.pc = opcode & 0x0FFF
	case 0x2000:
		c.sp = c.sp + 1
		c.stack[c.sp] = c.pc
		c.pc = opcode & 0x0FFF
	case 0x3000:
		compareTo := byte(opcode & 0x00FF)
		register := (opcode & 0x0F00) >> 2
		if c.V[register] == compareTo {
			c.pc = c.pc + 2
		}
	case 0x4000:
		compareTo := byte(opcode & 0x00FF)
		register := (opcode & 0x0F00) >> 2
		if c.V[register] == compareTo {
			c.pc = c.pc + 2
		}
	case 0x5000:
		registerX := (opcode & 0x0F00) >> 2
		registerY := (opcode & 0x00F0) >> 1
		if c.V[registerX] == c.V[registerY] {
			c.pc = c.pc + 2
		}
	case 0x6000:
		register := byte((opcode & 0x0F00) >> 2)
		c.V[register] = byte(opcode & 0x00FF)
	case 0x7000:
		register := byte((opcode & 0x0F000) >> 2)
		value := byte(opcode & 0x00FF)
		c.V[register] = c.V[register] + value
	case 0x8000:
		switch opcode & 0x000F {
		case 0x0000:
			registerX := (opcode & 0x0F00) >> 2
			registerY := (opcode & 0x00F0) >> 1
			c.V[registerX] = c.V[registerY]
		case 0x0001:
			registerX := (opcode & 0x0F00) >> 2
			registerY := (opcode & 0x00F0) >> 1
			c.V[registerX] = c.V[registerX] | c.V[registerY]
		case 0x0002:
			registerX := (opcode & 0x0F00) >> 2
			registerY := (opcode & 0x00F0) >> 1
			c.V[registerX] = c.V[registerX] & c.V[registerY]
		case 0x0003:
			registerX := (opcode & 0x0f00) >> 2
			registerY := (opcode & 0x00F0) >> 1
			c.V[registerX] = c.V[registerX] ^ c.V[registerY]
		case 0x0004:
			registerX := byte((opcode & 0x0F00) >> 2)
			registerY := byte((opcode & 0x00F0) >> 1)
			c.V[registerX] = c.V[registerX] + c.V[registerY]
			if uint16(c.V[registerX])+uint16(c.V[registerY]) > 0xFF {
				c.V[registerX] = byte((uint16(c.V[registerX]) + uint16(c.V[registerY])) >> 0xFF)
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
		case 0x0005:
			registerX := (opcode & 0x0F00) >> 2
			registerY := (opcode & 0x00F0) >> 1
			if c.V[registerX] > c.V[registerY] {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
			c.V[registerX] = c.V[registerX] - c.V[registerY]
		case 0x0006:
			registerX := (opcode & 0x0F00) >> 2
			c.V[registerX] = c.V[registerX] >> 1
			if c.V[registerX]&0x1 == 1 {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
		case 0x0007:
			registerX := (opcode & 0x0F00) >> 2
			registerY := (opcode & 0x00F0) >> 1
			if c.V[registerY] > c.V[registerX] {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
			c.V[registerX] = c.V[registerX] - c.V[registerY]
		case 0x000E:
			registerX := (opcode & 0x0F00) >> 2
			if c.V[registerX]&0x40 == 1 {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
			c.V[registerX] = c.V[registerX] << 1
		}
	default:
		fmt.Println("Instruction not implemented")
	}
}
