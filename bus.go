package main

import "fmt"

type Bus struct {
	CPU CPU6502
	RAM [64 * 1024]uint8
}

func (b *Bus) Initialize() {
	// Reset RAM
	for i := range b.RAM {
		b.RAM[i] = 0x00
	}

	b.RAM[0xFFFC] = 0x00
	b.RAM[0xFFFD] = 0x80

	// Create CPU with reference to this bus
	b.CPU = CPU6502{Bus: b}
	b.CPU.Initialize()
	b.CPU.Reset()
}

func (b *Bus) Write(addr uint16, data uint8) {
	if addr <= 0xFFFF {
		b.RAM[addr] = data
	} else {
		fmt.Println("Write bus outside address range:", addr)
	}
}

func (b *Bus) Read(addr uint16, readOnly bool) uint8 {
	if addr <= 0xFFFF {
		return b.RAM[addr]
	} else {
		fmt.Println("Read bus outside address range:", addr)
	}
	return 0x00
}
