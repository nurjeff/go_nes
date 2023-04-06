package bus

import (
	"github.com/sc-js/go_nes/cartridge"
)

func (b *Bus) InsertCartidge(c *cartridge.Cartridge) {
	b.Cartridge = c
	b.PPU.ConnectCartridge(b.Cartridge)
}

func (b *Bus) Reset() {
	b.CPU.Reset()
	b.SystemClockCounter = 0
}

func (b *Bus) Clock() {
	b.PPU.Clock()
	if b.SystemClockCounter%3 == 0 {
		if b.DMATransfer {
			if b.DMADummy {
				if b.SystemClockCounter%2 == 1 {
					b.DMADummy = false
				}
			} else {
				if b.SystemClockCounter%2 == 0 {
					b.DMAData = b.cpuRead((uint16(b.DMAPage)<<8)|uint16(b.DMAAddr), false)
				} else {

					b.PPU.POAM[b.DMAAddr] = b.DMAData
					b.DMAAddr++

					if b.DMAAddr == 0x00 {
						b.DMATransfer = false
						b.DMADummy = true
					}
				}
			}
		} else {
			b.CPU.Clock()
		}
	}

	if b.PPU.NMI {
		b.CPU.NMI()
		b.PPU.NMI = false
	}
	b.SystemClockCounter++
}
