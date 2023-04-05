package ppu2c02

type LoopyReg struct {
	reg uint16
}

const (
	_ uint16 = 1 << iota
	_
	_
	_
	_

	_
	_
	_
	_
	_

	LOOPY_REG_NAMETABLE_X
	LOOPY_REG_NAMETABLE_Y
)

func (l *LoopyReg) CoarseX() uint16 {
	var mask uint16 = 0x1F
	var bitoff uint16 = 0
	return (l.reg >> bitoff) & mask
}
func (l *LoopyReg) CoarseY() uint16 {
	var mask uint16 = 0x1F
	var bitoff uint16 = 5
	return (l.reg >> bitoff) & mask
}

func (l *LoopyReg) NametableX() uint16 {
	var mask uint16 = 0x1
	var bitoff uint16 = 10
	return (l.reg >> bitoff) & mask
}
func (l *LoopyReg) NametableY() uint16 {
	var mask uint16 = 0x1
	var bitoff uint16 = 11
	return (l.reg >> bitoff) & mask
}

func (l *LoopyReg) FineY() uint16 {
	var mask uint16 = 0x7
	var bitoff uint16 = 12
	return (l.reg >> bitoff) & mask
}

func (l *LoopyReg) SetCoarseX(val uint16) {
	var mask uint16 = 0x1F
	var bitoff uint16 = 0
	l.reg &^= mask << bitoff
	l.reg |= (val & mask) << bitoff
}
func (l *LoopyReg) SetCoarseY(val uint16) {
	var mask uint16 = 0x1F
	var bitoff uint16 = 5
	l.reg &^= mask << bitoff
	l.reg |= (val & mask) << bitoff
}

func (l *LoopyReg) SetNametableX(val uint16) {
	var mask uint16 = 0x1
	var bitoff uint16 = 10
	l.reg &^= mask << bitoff
	l.reg |= (val & mask) << bitoff
}
func (l *LoopyReg) SetNametableY(val uint16) {
	var mask uint16 = 0x1
	var bitoff uint16 = 11
	l.reg &^= mask << bitoff
	l.reg |= (val & mask) << bitoff
}

func (l *LoopyReg) SetFineY(val uint16) {
	var mask uint16 = 0x7
	var bitoff uint16 = 12
	l.reg &^= mask << bitoff
	l.reg |= (val & mask) << bitoff
}

/*
const (
	L_COARSE_X_MASK = 0b111110000000000
	L_COARSE_Y_MASK = 0b0000011111000000
	L_NT_X_MASK     = 0b0000000000100000
	L_NT_Y_MASK     = 0b0000000000010000
	L_FINE_Y_MASK   = 0b0000000000001110
	L_UNUSED_MASK   = 0b0000000000000001
)

func (p *PPUC202) GetLVRAM(flag int) uint16 {
	switch flag {
	case L_COARSE_X:
		return (p.LVRAM & L_COARSE_X_MASK) >> 11
	case L_COARSE_Y:
		return (p.LVRAM & L_COARSE_Y_MASK) >> 6
	case L_NT_X:
		return (p.LVRAM & L_NT_X_MASK) >> 5
	case L_NT_Y:
		return (p.LVRAM & L_NT_Y_MASK) >> 4
	case L_FINE_Y:
		return (p.LVRAM & L_FINE_Y_MASK) >> 1
	case L_UNUSED:
		return (p.LVRAM & L_UNUSED_MASK) >> 0
	}
	return 0
}

func (p *PPUC202) GetLTRAM(flag int) uint16 {
	switch flag {
	case L_COARSE_X:
		return (p.LTRAM & L_COARSE_X_MASK) >> 11
	case L_COARSE_Y:
		return (p.LTRAM & L_COARSE_Y_MASK) >> 6
	case L_NT_X:
		return (p.LTRAM & L_NT_X_MASK) >> 5
	case L_NT_Y:
		return (p.LTRAM & L_NT_Y_MASK) >> 4
	case L_FINE_Y:
		return (p.LTRAM & L_FINE_Y_MASK) >> 1
	case L_UNUSED:
		return (p.LTRAM & L_UNUSED_MASK) >> 0
	}
	return 0
}

func (p *PPUC202) SetLVRAM(flag int, val uint16) {
	var bitsToSet uint16 = 0x0000
	switch flag {
	case L_COARSE_X:
		bitsToSet = (p.GetLVRAM(L_COARSE_X) ^ val) << 11
	case L_COARSE_Y:
		bitsToSet = (p.GetLVRAM(L_COARSE_Y) ^ val) << 6
	case L_NT_X:
		bitsToSet = (p.GetLVRAM(L_NT_X) ^ val) << 5
	case L_NT_Y:
		bitsToSet = (p.GetLVRAM(L_NT_Y) ^ val) << 4
	case L_FINE_Y:
		bitsToSet = (p.GetLVRAM(L_FINE_Y) ^ val) << 1
	case L_UNUSED:
		bitsToSet = (p.GetLVRAM(L_UNUSED) ^ val) << 0
	}

	p.LVRAM = bitsToSet ^ p.LVRAM
}

func (p *PPUC202) SETLTRAM(flag int, val uint16) {
	var bitsToSet uint16 = 0x0000
	switch flag {
	case L_COARSE_X:
		bitsToSet = (p.GetLTRAM(L_COARSE_X) ^ val) << 11
	case L_COARSE_Y:
		bitsToSet = (p.GetLTRAM(L_COARSE_Y) ^ val) << 6
	case L_NT_X:
		bitsToSet = (p.GetLTRAM(L_NT_X) ^ val) << 5
	case L_NT_Y:
		bitsToSet = (p.GetLTRAM(L_NT_Y) ^ val) << 4
	case L_FINE_Y:
		bitsToSet = (p.GetLTRAM(L_FINE_Y) ^ val) << 1
	case L_UNUSED:
		bitsToSet = (p.GetLTRAM(L_UNUSED) ^ val) << 0
	}

	p.LTRAM = bitsToSet ^ p.LTRAM
}
*/
