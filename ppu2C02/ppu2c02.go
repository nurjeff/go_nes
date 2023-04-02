package ppu2c02

import (
	"fmt"

	"github.com/sc-js/go_nes/cartridge"
)

type PPUC202 struct {
	Cartridge *cartridge.Cartridge

	Nametable [2][1024]uint8 // VRAM
	Palette   [32]uint8
	Pattern   [2][4096]uint8 // This is really on the cartridge

	cycle         int16
	scanline      int16
	frameComplete bool
}

func (p *PPUC202) ConnectCartridge(cartridge *cartridge.Cartridge) {
	p.Cartridge = cartridge
}

func (p *PPUC202) Clock() {
	p.cycle++

	if p.cycle >= 341 {
		p.cycle = 0
		p.scanline++
		if p.scanline >= 261 {
			p.scanline = -1
			p.frameComplete = true
		}
	}
}

func (p *PPUC202) CPUWrite(addr uint16, data uint8) {
	switch addr {
	case 0x0000: // Control
		break
	case 0x0001: //Mask
		break
	case 0x0002: //Status
		break
	case 0x0003: // OAM Address
		break
	case 0x0004: //OAM Data
		break
	case 0x0005: // Scroll
		break
	case 0x0006: // PPU Address
		break
	case 0x0007: // PPU Data
		break
	default:
		panic("cpu tried accessing forbidden data:" + fmt.Sprint(addr))
	}
}

func (p *PPUC202) CPURead(addr uint16, readOnly bool) uint8 {
	var data uint8 = 0x00
	switch addr {
	case 0x0000: // Control
		break
	case 0x0001: //Mask
		break
	case 0x0002: //Status
		break
	case 0x0003: // OAM Address
		break
	case 0x0004: //OAM Data
		break
	case 0x0005: // Scroll
		break
	case 0x0006: // PPU Address
		break
	case 0x0007: // PPU Data
		break
	default:
		panic("cpu tried accessing forbidden data:" + fmt.Sprint(addr))
	}
	return data
}

func (p *PPUC202) PPURead(addr uint16, readOnly bool) uint8 {
	var data uint8 = 0x00
	addr &= 0x3FFF

	if p.Cartridge.PPURead(addr, &data) {
		fmt.Println("ppu read")
	}

	return data
}

func (p *PPUC202) PPUWrite(addr uint16, data uint8) {
	addr &= 0x3FFF

	if p.Cartridge.PPUWrite(addr, data) {
		fmt.Println("ppu write")
	}
}
