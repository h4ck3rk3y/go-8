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
	c := cpu{}
	return c
}

func (c cpu) LoadProgram(rom string) int {
	f, err := os.Open(rom)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	memory := make([]byte, 4096)
	n, err := f.Read(memory)
	for index, b := range memory {
		c.memory[index] = b
	}
	if err != nil {
		if err != io.EOF {
			panic(err)
		}
	}
	return n
}
