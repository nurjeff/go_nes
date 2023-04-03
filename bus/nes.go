package bus

import "github.com/sc-js/go_nes/cartridge"

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
		b.CPU.Clock()
	}

	if b.PPU.NMI {
		b.PPU.NMI = false
		b.CPU.NMI()
	}
	b.SystemClockCounter++
}
