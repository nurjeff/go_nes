package bus

import (
	"github.com/sc-js/go_nes/cartridge"
	"github.com/sc-js/go_nes/cpu6502"
	ppu2c02 "github.com/sc-js/go_nes/ppu2C02"
)

type Bus struct {
	CPU       cpu6502.CPU6502
	PPU       ppu2c02.PPUC202
	Cartridge *cartridge.Cartridge
	CPURAM    [2048]uint8

	SystemClockCounter uint64
}

func (b *Bus) Initialize() {
	if b.Cartridge == nil {
		panic("insert cartridge first")
	}
	// Reset RAM
	for i := range b.CPURAM {
		b.CPURAM[i] = 0x00
	}

	b.CPURAM[0xFFFC&0x07FF] = 0x00
	b.CPURAM[0xFFFD&0x07FF] = 0x80

	// Create CPU with reference to this bus
	b.CPU = cpu6502.CPU6502{ReadBus: b.cpuRead, WriteBus: b.cpuWrite}

	b.CPU.Initialize()
	b.CPU.Disassemble(0x0000, 0xFFFF)
	b.CPU.Reset()

	b.PPU = ppu2c02.PPUC202{Cartridge: b.Cartridge}
}

func (b *Bus) cpuWrite(addr uint16, data uint8) {
	if b.Cartridge.CPUWrite(addr, data) {
	} else if addr <= 0x1FFF {
		b.CPURAM[addr&0x07FF] = data
	} else if addr >= 0x2000 && addr <= 0x3FFF {
		b.PPU.CPUWrite(addr&0x0007, data)
	}
}

func (b *Bus) cpuRead(addr uint16, readOnly bool) uint8 {
	var data uint8 = 0x00
	if b.Cartridge.CPURead(addr, &data) {
	} else if addr <= 0x1FFF {
		return b.CPURAM[addr&0x07FF]
	} else if addr >= 0x2000 && addr <= 0x3FFF {
		return b.PPU.CPURead(addr&0x0007, readOnly)
	}
	return data
}
