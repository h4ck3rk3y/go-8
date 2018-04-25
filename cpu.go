package main

import (
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
