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
