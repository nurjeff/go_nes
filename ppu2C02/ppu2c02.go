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

	LVRAM LoopyReg
	LTRAM LoopyReg

	FineX uint8

	NMI bool

	palScreen [0x40]Pixel

	sprNameTable    [2]Display
	sprPatternTable [2]Display
	SprScreen       Display

	BGNextTileID     uint8
	BGNextTileAttrib uint8
	BGNextTileLSB    uint8
	BGNextTileMSB    uint8

	BGShifterPatternLO uint16
	BGShifterPatternHI uint16
	BGShifterAttribLO  uint16
	BGShifterAttribHI  uint16
}

func (p *PPUC202) GetColourFromPaletteRam(palette uint8, pixel uint8) Pixel {
	temp := p.PPURead(0x3F00+(uint16(palette)<<2)+uint16(pixel), false)
	val := p.palScreen[temp&0x3F]
	//fmt.Println(val)
	return val
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
	/*
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
	*/

	S_VBLANK_FLAG     = (1 << 7)
	S_SPRITEZERO_HIT  = (1 << 6)
	S_SPRITE_OVERFLOW = (1 << 5)

	C_ENABLE_MMI     = (1 << 7)
	C_SLAVE_MODE     = (1 << 6)
	C_SPRITE_SIZE    = (1 << 5)
	C_PATTERN_BKG    = (1 << 4)
	C_PATTERN_SPR    = (1 << 3)
	C_INCREMENT_MODE = (1 << 2)
	C_NAMETABLE_Y    = (1 << 1)
	C_NAMETABLE_X    = (1 << 0)

	M_ENHANCE_BLUE    = (1 << 7)
	M_ENHANGE_GREEN   = (1 << 6)
	M_ENHANCE_RED     = (1 << 5)
	M_RENDER_SPR      = (1 << 4)
	M_RENDER_BKG      = (1 << 3)
	M_RENDER_SPR_LEFT = (1 << 2)
	M_RENDER_BKG_LEFT = (1 << 1)
	M_GRAYSCALE       = (1 << 0)

	L_COARSE_X = (1 << 11)
	L_COARSE_Y = (1 << 6)
	L_NT_X     = (1 << 5)
	L_NT_Y     = (1 << 4)
	L_FINE_Y   = (1 << 1)
	L_UNUSED   = (1 << 0)
)

func (p *PPUC202) Initialize() {
	p.palScreen = GetPal()
	p.sprNameTable[0] = Display{Width: 256, Height: 240}
	p.sprNameTable[1] = Display{Width: 256, Height: 240}
	p.sprPatternTable[0] = Display{Width: 128, Height: 128}
	p.sprPatternTable[1] = Display{Width: 128, Height: 128}
	p.SprScreen = Display{Width: 256, Height: 240}
	p.SprScreen.Initialize()

	p.ControlReg = 0x00
	p.MaskReg = 0x00
	p.StatusReg = 0x00
	p.SetStatus(S_VBLANK_FLAG, true)
	p.SetStatus(S_SPRITE_OVERFLOW, true)
}

func (p *PPUC202) ConnectCartridge(cartridge *cartridge.Cartridge) {
	p.Cartridge = cartridge
}

func (p *PPUC202) IncrementScrollX() {
	if p.GetMask(M_RENDER_BKG) || p.GetMask(M_RENDER_SPR) {
		if p.LVRAM.CoarseX() == 31 {
			p.LVRAM.SetCoarseX(0)
			p.LVRAM.SetNametableX(^p.LVRAM.NametableX())
		} else {
			p.LVRAM.SetCoarseX(p.LVRAM.CoarseX() + 1)
		}
	}
}

func (p *PPUC202) IncrementScrollY() {
	if p.GetMask(M_RENDER_BKG) || p.GetMask(M_RENDER_SPR) {
		if p.LVRAM.FineY() < 7 {
			p.LVRAM.SetFineY(p.LVRAM.FineY() + 1)
		} else {
			p.LVRAM.SetFineY(0)
			if p.LVRAM.CoarseY() == 29 {
				p.LVRAM.SetCoarseY(0)
				p.LVRAM.SetNametableY(^p.LVRAM.NametableY())
			} else if p.LVRAM.CoarseY() == 31 {
				p.LVRAM.SetCoarseY(0)
			} else {
				p.LVRAM.SetCoarseY(p.LVRAM.CoarseY() + 1)
			}
		}
	}
}

func (p *PPUC202) TransferAddressX() {
	if p.GetMask(M_RENDER_BKG) || p.GetMask(M_RENDER_SPR) {
		p.LVRAM.SetNametableX(p.LTRAM.NametableX())
		p.LVRAM.SetCoarseX(p.LTRAM.CoarseX())
	}
}

func (p *PPUC202) TransferAddressY() {
	if p.GetMask(M_RENDER_BKG) || p.GetMask(M_RENDER_SPR) {
		p.LVRAM.SetNametableY(p.LTRAM.NametableY())
		p.LVRAM.SetCoarseY(p.LTRAM.CoarseY())
		p.LVRAM.SetFineY(p.LTRAM.FineY())
	}
}

func (p *PPUC202) LoadBackgroundShifters() {
	p.BGShifterPatternLO = (p.BGShifterPatternLO & 0xFF00) | uint16(p.BGNextTileLSB)
	p.BGShifterPatternHI = (p.BGShifterPatternHI & 0xFF00) | uint16(p.BGNextTileMSB)

	var nextTileAttribModLO uint16 = 0x00
	var nextTileAttribModHI uint16 = 0x00

	if (p.BGNextTileAttrib & 0b01) >= 1 {
		nextTileAttribModLO = 0xFF
	}
	if (p.BGNextTileAttrib & 0b10) >= 1 {
		nextTileAttribModHI = 0xFF
	}

	p.BGShifterAttribLO = (p.BGShifterAttribLO & 0xFF00) | nextTileAttribModLO
	p.BGShifterAttribHI = (p.BGShifterAttribHI & 0xFF00) | nextTileAttribModHI
}

func (p *PPUC202) UpdateShifters() {
	if p.GetMask(M_RENDER_BKG) {
		p.BGShifterPatternLO <<= 1
		p.BGShifterPatternHI <<= 1
		p.BGShifterAttribLO <<= 1
		p.BGShifterAttribHI <<= 1
	}
}

func (p *PPUC202) Clock() {
	if p.scanline >= -1 && p.scanline < 240 {

		if p.scanline == 0 && p.cycle == 0 {
			p.cycle = 1
		}

		if p.scanline == -1 && p.cycle == 1 {
			p.SetStatus(S_VBLANK_FLAG, false)
		}

		if (p.cycle >= 2 && p.cycle < 258) || (p.cycle >= 321 && p.cycle < 338) {

			p.UpdateShifters()

			switch (p.cycle - 1) % 8 {
			case 0:
				p.LoadBackgroundShifters()
				p.BGNextTileID = p.PPURead(0x2000|(p.LVRAM.reg&0x0FFF), false)
			case 2:

				p.BGNextTileAttrib = p.PPURead(0x23C0|((p.LVRAM.NametableY())<<11)|(p.LVRAM.NametableX()<<10)|((p.LVRAM.CoarseY()>>2)<<3)|(p.LVRAM.CoarseX()>>2), false)

				if (p.LVRAM.CoarseY() & 0x02) >= 1 {
					p.BGNextTileAttrib >>= 4
				}
				if (p.LVRAM.CoarseX() & 0x02) >= 1 {
					p.BGNextTileAttrib >>= 2
				}
				p.BGNextTileAttrib &= 0x03
			case 4:
				var f uint16 = 0
				if p.GetControl(C_PATTERN_BKG) {
					f = 1
				}
				f = f << 12
				p.BGNextTileLSB = p.PPURead(f+(uint16(p.BGNextTileID)<<4)+(p.LVRAM.FineY()+0), false)
			case 6:
				var f uint16 = 0
				if p.GetControl(C_PATTERN_BKG) {
					f = 1
				}
				f = f << 12
				p.BGNextTileMSB = p.PPURead(f+(uint16(p.BGNextTileID)<<4)+(p.LVRAM.FineY()+8), false)
			case 7:
				p.IncrementScrollX()
			}
		}

		if p.cycle == 256 {
			p.IncrementScrollY()
		}

		if p.cycle == 257 {
			p.LoadBackgroundShifters()
			p.TransferAddressX()
		}

		if p.cycle == 338 || p.cycle == 340 {
			p.BGNextTileID = p.PPURead(0x2000|(p.LVRAM.reg&0x0FFF), false)
		}

		if p.scanline == -1 && p.cycle >= 280 && p.cycle < 305 {
			p.TransferAddressY()
		}
	}

	/*if p.scanline == 240 {
		// Nothing happens
	}*/

	if p.scanline >= 241 && p.scanline < 261 {
		if p.scanline == 241 && p.cycle == 1 {
			p.SetStatus(S_VBLANK_FLAG, true)

			if p.GetControl(C_ENABLE_MMI) {
				//fmt.Println("nmi")
				p.NMI = true
			}
		}
	}

	var BGPixel uint8 = 0x00
	var BGPalette uint8 = 0x00

	if p.GetMask(M_RENDER_BKG) {
		var bit_mux uint16 = 0x8000 >> p.FineX
		var p0_pixel uint8
		var p1_pixel uint8
		if (p.BGShifterPatternLO & bit_mux) > 0 {
			p0_pixel = 1
		}
		if (p.BGShifterPatternHI & bit_mux) > 0 {
			p1_pixel = 1
		}
		BGPixel = (p1_pixel << 1) | p0_pixel

		var bg_pal0 uint8
		var bg_pal1 uint8
		if (p.BGShifterAttribLO & bit_mux) > 0 {
			bg_pal0 = 1
		}
		if (p.BGShifterAttribHI & bit_mux) > 0 {
			bg_pal1 = 1
		}
		BGPalette = (bg_pal1 << 1) | bg_pal0
	}

	p.SprScreen.SetPixel(int32(p.cycle)-1, int32(p.scanline), p.GetColourFromPaletteRam(BGPalette, BGPixel))

	p.cycle++

	if p.cycle >= 341 {
		p.cycle = 0
		p.scanline++
		if p.scanline >= 261 {
			p.scanline = -1
			p.FrameComplete = true
		}
	}
}

func (p *PPUC202) CPUWrite(addr uint16, data uint8) {
	switch addr {
	case 0x0000: // Control
		p.ControlReg = data
		if p.GetControl(C_NAMETABLE_X) {
			//p.SETLTRAM(L_NT_X, 1)
			p.LTRAM.SetNametableX(1)
		} else {
			p.LTRAM.SetNametableX(0)
		}
		if p.GetControl(C_NAMETABLE_Y) {
			p.LTRAM.SetNametableY(1)
		} else {
			p.LTRAM.SetNametableY(0)
		}
	case 0x0001: //Mask
		p.MaskReg = data
	case 0x0002: //Status
	case 0x0003: // OAM Address
		break
	case 0x0004: //OAM Data
		break
	case 0x0005: // Scroll
		if p.AddressLatch == 0 {
			p.FineX = data & 0x07
			p.LTRAM.SetCoarseX(uint16(data) >> 3)
			p.AddressLatch = 1
		} else {
			p.LTRAM.SetFineY(uint16(data) & 0x07)
			p.LTRAM.SetCoarseY(uint16(data) >> 3)
			p.AddressLatch = 0
		}
	case 0x0006: // PPU Address
		if p.AddressLatch == 0 {
			p.LTRAM.reg = (p.LTRAM.reg & 0x00FF) | (uint16(data) << 8)
			p.AddressLatch = 1
		} else {
			p.LTRAM.reg = (p.LTRAM.reg & 0xFF00) | uint16(data)
			p.LVRAM.reg = p.LTRAM.reg
			p.AddressLatch = 0
		}
	case 0x0007: // PPU Data
		p.PPUWrite(p.LVRAM.reg, data)
		if p.GetControl(C_INCREMENT_MODE) {
			p.LVRAM.reg += 32
		} else {
			p.LVRAM.reg++
		}

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

func (p *PPUC202) GetMask(FLAG uint8) bool {
	if (p.MaskReg & FLAG) > 0 {
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

func (p *PPUC202) SetMask(FLAG uint8, v bool) {
	if v {
		p.MaskReg |= FLAG
	} else {
		p.MaskReg &= ^FLAG
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
		p.DataBuffer = p.PPURead(p.LVRAM.reg, false)
		if p.LVRAM.reg >= 0x3F00 {
			data = p.DataBuffer
		}
		if p.GetControl(C_INCREMENT_MODE) {
			p.LVRAM.reg += 32
		} else {
			p.LVRAM.reg++
		}

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
		addr &= 0x0FFF
		if p.Cartridge.Mirror == cartridge.VERTICAL {
			if addr <= 0x03FF {
				data = p.Nametable[0][addr&0x3FF]
			}
			if addr >= 0x0400 && addr <= 0x07FF {
				data = p.Nametable[1][addr&0x03FF]
			}
			if addr >= 0x0800 && addr <= 0x0BFF {
				data = p.Nametable[0][addr&0x03FF]
			}
			if addr >= 0x0C00 && addr <= 0x0FFF {
				data = p.Nametable[1][addr&0x03FF]
			}
		} else if p.Cartridge.Mirror == cartridge.HORIZONTAL {
			if addr <= 0x03FF {
				data = p.Nametable[0][addr&0x3FF]
			}
			if addr >= 0x0400 && addr <= 0x07FF {
				data = p.Nametable[0][addr&0x3FF]
			}
			if addr >= 0x0800 && addr <= 0x0BFF {
				data = p.Nametable[1][addr&0x3FF]
			}
			if addr >= 0x0C00 && addr <= 0x0FFF {
				data = p.Nametable[1][addr&0x3FF]
			}
		}

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
		addr &= 0x0FFF
		if p.Cartridge.Mirror == cartridge.VERTICAL {

			if addr <= 0x03FF {
				p.Nametable[0][addr&0x3FF] = data
			}
			if addr >= 0x0400 && addr <= 0x07FF {
				p.Nametable[1][addr&0x03FF] = data
			}
			if addr >= 0x0800 && addr <= 0x0BFF {
				p.Nametable[0][addr&0x03FF] = data
			}
			if addr >= 0x0C00 && addr <= 0x0FFF {
				p.Nametable[1][addr&0x03FF] = data
			}
		} else if p.Cartridge.Mirror == cartridge.HORIZONTAL {
			//fmt.Println(emutools.Hex(addr, 4))
			if addr <= 0x03FF {
				p.Nametable[0][addr&0x3FF] = data
			}
			if addr >= 0x0400 && addr <= 0x07FF {
				p.Nametable[0][addr&0x3FF] = data
			}
			if addr >= 0x0800 && addr <= 0x0BFF {
				p.Nametable[1][addr&0x3FF] = data
			}
			if addr >= 0x0C00 && addr <= 0x0FFF {
				p.Nametable[1][addr&0x3FF] = data
			}
		}
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
	} else {
		fmt.Println("invalid ppu write")
	}
}
