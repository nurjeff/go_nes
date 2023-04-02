package cpu6502

// Implement Adressing Modes
func (c *CPU6502) IMP() uint8 {
	c.fetched = c.a
	return 0
}

func (c *CPU6502) IMM() uint8 {
	c.addr_abs = c.pc
	c.pc++
	return 0
}

func (c *CPU6502) ZP0() uint8 {
	c.addr_abs = uint16(c.Read(c.pc))
	c.pc++
	c.addr_abs &= 0x00FF
	return 0
}

func (c *CPU6502) ZPX() uint8 {
	c.addr_abs = (uint16(c.Read(c.pc) + c.x))
	c.pc++
	c.addr_abs &= 0x00FF
	return 0
}

func (c *CPU6502) ZPY() uint8 {
	c.addr_abs = (uint16(c.Read(c.pc) + c.y))
	c.pc++
	c.addr_abs &= 0x00FF
	return 0
}

func (c *CPU6502) REL() uint8 {
	c.addr_rel = uint16(c.Read(c.pc))
	c.pc++
	// might bug out

	if (c.addr_rel & 0x80 >> 7) == 1 {
		c.addr_rel |= 0xFF00
	}
	return 0x00
}

func (c *CPU6502) ABS() uint8 {
	var lo uint16 = uint16(c.Read(c.pc))
	c.pc++
	var hi uint16 = uint16(c.Read(c.pc))
	c.pc++
	c.addr_abs = (hi << 8) | lo
	return 0
}

func (c *CPU6502) ABX() uint8 {
	var lo uint16 = uint16(c.Read(c.pc))
	c.pc++
	var hi uint16 = uint16(c.Read(c.pc))
	c.pc++

	c.addr_abs = (hi << 8) | lo
	c.addr_abs += uint16(c.x)
	if c.addr_abs&0xFF00 != (hi << 8) {
		return 1
	}
	return 0
}

func (c *CPU6502) ABY() uint8 {
	var lo uint16 = uint16(c.Read(c.pc))
	c.pc++
	var hi uint16 = uint16(c.Read(c.pc))
	c.pc++

	c.addr_abs = (hi << 8) | lo
	c.addr_abs += uint16(c.y)
	if c.addr_abs&0xFF00 != (hi << 8) {
		return 1
	}
	return 0
}

func (c *CPU6502) IND() uint8 {
	var ptr_lo uint16 = uint16(c.Read(c.pc))
	c.pc++
	var ptr_hi uint16 = uint16(c.Read(c.pc))
	c.pc++

	var ptr uint16 = (ptr_hi << 8) | ptr_lo
	if ptr_lo == 0x00FF {
		c.addr_abs = (uint16(c.Read(ptr&0xFF00)) << 8) | uint16(c.Read(ptr+0))
	} else {
		c.addr_abs = (uint16(c.Read(ptr+1)) << 8) | uint16(c.Read(ptr+0))
	}
	return 0
}

func (c *CPU6502) IZX() uint8 {
	var t uint16 = uint16(c.Read(c.pc))
	c.pc++

	var lo uint16 = uint16(c.Read((t + uint16(c.x)) & 0x00FF))
	var hi uint16 = uint16(c.Read((t + uint16(c.x) + 1) & 0x00FF))

	c.addr_abs = (hi << 8) | lo

	return 0
}

func (c *CPU6502) IZY() uint8 {
	var t uint16 = uint16(c.Read(c.pc))
	c.pc++

	var lo uint16 = uint16(c.Read(t & 0x00FF))
	var hi uint16 = uint16(c.Read((t + 1) & 0x00FF))

	c.addr_abs = (hi << 8) | lo
	c.addr_abs += uint16(c.y)

	if (c.addr_abs & 0xFF00) != (hi << 8) {
		return 1
	}
	return 0
}

func (c *CPU6502) NOPX() uint8 {
	var ptr_lo uint16 = uint16(c.Read(c.pc))
	c.pc++
	var ptr_hi uint16 = uint16(c.Read(c.pc))
	c.pc++

	_ = ptr_lo
	_ = ptr_hi
	return 1
}
