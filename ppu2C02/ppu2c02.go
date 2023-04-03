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
	FrameComplete bool

	ControlReg uint8
	MaskReg    uint8
	StatusReg  uint8

	AddressLatch uint8
	DataBuffer   uint8
	Address      uint16

	NMI bool

	palScreen [0x40]Pixel

	sprNameTable    [2]Display
	sprPatternTable [2]Display
	sprScreen       Display
}

func (p *PPUC202) GetColourFromPaletteRam(palette uint8, pixel uint8) Pixel {
	temp := p.PPURead(0x3F00+(uint16(palette)<<2)+uint16(pixel), false)
	return p.palScreen[temp&0x3F]
}

func (p *PPUC202) GetNameTable(i uint8) *Display {
	return &p.sprNameTable[i]
}

func (p *PPUC202) GetPatternTable(i uint8, palette uint8) *Display {
	for nTileY := 0; nTileY < 16; nTileY++ {
		for nTileX := 0; nTileX < 16; nTileX++ {
			nOffset := nTileY*256 + nTileX*16
			for row := 0; row < 8; row++ {
				var tileLSB uint8 = p.PPURead(uint16(i)*0x1000+uint16(nOffset)+uint16(row)+0x0000, false)
				var tileMSB uint8 = p.PPURead(uint16(i)*0x1000+uint16(nOffset)+uint16(row)+0x0008, false)
				for col := 0; col < 8; col++ {
					var pixel uint8 = (tileLSB & 0x01) + ((tileMSB & 0x01) << 1)
					tileLSB >>= 1
					tileMSB >>= 1
					p.sprPatternTable[i].SetPixel(
						int32(nTileX)*8+(7-int32(col)),
						int32(nTileY)*8+int32(row),
						p.GetColourFromPaletteRam(palette, pixel))
				}
			}
		}
	}
	return &p.sprPatternTable[i]
}

const (
	S_VBLANK_FLAG     = (1 << 0)
	S_SPRITEZERO_HIT  = (1 << 1)
	S_SPRITE_OVERFLOW = (1 << 2)

	C_ENABLE_MMI     = (1 << 0)
	C_SLAVE_MODE     = (1 << 1)
	C_SPRITE_SIZE    = (1 << 2)
	C_PATTERN_BKG    = (1 << 3)
	C_PATTERN_SPR    = (1 << 4)
	C_INCREMENT_MODE = (1 << 5)
	C_NAMETABLE_X    = (1 << 6)
	C_NAMETABLE_Y    = (1 << 7)
)

func (p *PPUC202) Initialize() {
	p.palScreen = GetPal()
	p.sprNameTable[0] = Display{Width: 256, Height: 240}
	p.sprNameTable[1] = Display{Width: 256, Height: 240}
	p.sprPatternTable[0] = Display{Width: 128, Height: 128}
	p.sprPatternTable[1] = Display{Width: 128, Height: 128}

	p.ControlReg = 0xFF
	p.MaskReg = 0xFF
	p.StatusReg = 0xFF
}

func (p *PPUC202) ConnectCartridge(cartridge *cartridge.Cartridge) {
	p.Cartridge = cartridge
}

func (p *PPUC202) Clock() {
	if p.scanline == -1 && p.cycle == 1 {
		p.SetStatus(S_VBLANK_FLAG, false)
	}

	if p.scanline == 241 && p.cycle == 1 {
		p.SetStatus(S_VBLANK_FLAG, true)
		if p.GetControl(C_ENABLE_MMI) {
			p.NMI = true
		}
	}

	p.cycle++

	if p.cycle >= 341 {
		p.cycle = 0
		p.scanline++
		if p.scanline >= 261 {
			p.scanline = -1
			p.FrameComplete = true
			fmt.Println("Frame completed")
		}
	}
}

func (p *PPUC202) CPUWrite(addr uint16, data uint8) {
	switch addr {
	case 0x0000: // Control
		p.ControlReg = data
	case 0x0001: //Mask
		p.MaskReg = data
	case 0x0002: //Status
	case 0x0003: // OAM Address
		break
	case 0x0004: //OAM Data
		break
	case 0x0005: // Scroll
		break
	case 0x0006: // PPU Address
		if p.AddressLatch == 0 {
			p.Address = (p.Address & 0x00FF) | (uint16(data) << 8)
			p.AddressLatch = 1
		} else {
			p.Address = (p.Address & 0xFF00) | uint16(data)
			p.AddressLatch = 0
		}
	case 0x0007: // PPU Data
		p.PPUWrite(p.Address, data)
		p.Address++
	default:
		panic("cpu tried accessing forbidden data [WRITE]:" + fmt.Sprint(addr))
	}
}

func (p *PPUC202) SetStatus(FLAG uint8, v bool) {
	if v {
		p.StatusReg |= FLAG
	} else {
		p.StatusReg &= ^FLAG
	}
}

func (p *PPUC202) GetStatus(FLAG uint8) bool {
	if (p.StatusReg & FLAG) > 0 {
		return true
	} else {
		return false
	}
}

func (p *PPUC202) GetControl(FLAG uint8) bool {
	if (p.ControlReg & FLAG) > 0 {
		return true
	} else {
		return false
	}
}

func (p *PPUC202) SetControl(FLAG uint8, v bool) {
	if v {
		p.ControlReg |= FLAG
	} else {
		p.ControlReg &= ^FLAG
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
		data = (p.StatusReg & 0xE0) | (p.DataBuffer & 0x1F)
		p.SetStatus(S_VBLANK_FLAG, false)
		p.AddressLatch = 0
	case 0x0003: // OAM Address
		break
	case 0x0004: //OAM Data
		break
	case 0x0005: // Scroll
		break
	case 0x0006: // PPU Address
		break
	case 0x0007: // PPU Data
		data = p.DataBuffer
		p.DataBuffer = p.PPURead(p.Address, false)
		if p.Address > 0x3F00 {
			data = p.DataBuffer
		}
		p.Address++
	default:
		panic("cpu tried accessing forbidden data [READ]:" + fmt.Sprint(addr))
	}
	return data
}

func (p *PPUC202) PPURead(addr uint16, readOnly bool) uint8 {
	var data uint8 = 0x00
	addr &= 0x3FFF

	if p.Cartridge.PPURead(addr, &data) {
	} else if addr <= 0x1FFF {
		data = p.Pattern[(addr&0x1000)>>12][addr&0x0FFF]

	} else if addr >= 0x2000 && addr <= 0x3EFF {
	} else if addr >= 0x3F00 && addr <= 0x3FFF {
		addr &= 0x001F
		if addr == 0x0010 {
			addr = 0x0000
		} else if addr == 0x0014 {
			addr = 0x0004
		} else if addr == 0x0018 {
			addr = 0x0008
		} else if addr == 0x001C {
			addr = 0x000C
		}
		data = p.Palette[addr]
	}

	return data
}

func (p *PPUC202) PPUWrite(addr uint16, data uint8) {
	addr &= 0x3FFF

	if p.Cartridge.PPUWrite(addr, data) {
	} else if addr <= 0x1FFF {
		p.Pattern[(addr&0x1000)>>12][addr&0x0FFF] = data
	} else if addr >= 0x2000 && addr <= 0x3EFF {

	} else if addr >= 0x3F00 && addr <= 0x3FFF {
		addr &= 0x001F
		if addr == 0x0010 {
			addr = 0x0000
		} else if addr == 0x0014 {
			addr = 0x0004
		} else if addr == 0x0018 {
			addr = 0x0008
		} else if addr == 0x001C {
			addr = 0x000C
		}
		p.Palette[addr] = data
	}
}
