package bus

import (
	"time"

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

	Controller      [2]uint8
	ControllerState [2]uint8

	DMAPage uint8
	DMAAddr uint8
	DMAData uint8

	DMATransfer bool
	DMADummy    bool
}

func (b *Bus) Initialize() {
	b.DMADummy = true
	if b.Cartridge == nil {
		panic("insert cartridge first")
	}
	// Reset RAM
	for i := range b.CPURAM {
		b.CPURAM[i] = 0x00
	}
	// Create CPU with reference to this bus
	b.PPU = ppu2c02.PPUC202{Cartridge: b.Cartridge}
	for index := range b.PPU.OAM {
		b.PPU.OAM[index] = ppu2c02.ObjectAttributeEntity{}
	}
	b.PPU.Initialize()
	b.CPU = cpu6502.CPU6502{ReadBus: b.cpuRead, WriteBus: b.cpuWrite}

	b.CPU.Initialize()
	b.CPU.Disassemble(0x0000, 0xFFFF)
	b.CPU.Reset()
	b.Cartridge.Mapper.Reset()

}

func (b *Bus) cpuWrite(addr uint16, data uint8) {

	if b.Cartridge.CPUWrite(addr, data) {
	} else if addr <= 0x1FFF {
		b.CPURAM[addr&0x07FF] = data
	} else if addr >= 0x2000 && addr <= 0x3FFF {
		b.PPU.CPUWrite(addr&0x0007, data)
	} else if addr == 0x4014 {
		b.DMAPage = data
		b.DMAAddr = 0x00
		b.DMATransfer = true
	} else if addr >= 0x4016 && addr <= 0x4017 {
		b.ControllerState[addr&0x0001] = b.Controller[addr&0x0001]
	}
}

func (b *Bus) cpuRead(addr uint16, readOnly bool) uint8 {
	var data uint8 = 0x00
	if b.Cartridge.CPURead(addr, &data) {
	} else if addr <= 0x1FFF {
		return b.CPURAM[addr&0x07FF]
	} else if addr >= 0x2000 && addr <= 0x3FFF {
		return b.PPU.CPURead(addr&0x0007, readOnly)
	} else if addr >= 0x4016 && addr <= 0x4017 {
		if (b.ControllerState[addr&0x0001] & 0x80) > 0 {
			data = 1
		} else {
			data = 0
		}
		b.ControllerState[addr&0x0001] <<= 1
	}

	return data
}

func (b *Bus) Boot() {
	var dur time.Duration
	for {
		st := time.Now()
		for {
			b.Clock()
			if b.PPU.FrameComplete {
				b.PPU.FrameComplete = false
				for b.CPU.GetCycles() > 0 {
					b.CPU.Clock()
				}

				break
			}
		}
		dur = time.Since(st)
		if dur.Milliseconds() < 16 {
			time.Sleep(time.Millisecond*16 - dur)
		}

	}
}
