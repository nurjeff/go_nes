package main

func (c *CPU6502) ADC() uint8 {
	c.Fetch()
	tmp := uint16(c.a) + uint16(c.fetched) + uint16(c.GetFlag(c.flags.C))
	c.SetFlag(c.flags.C, tmp > 255)
	c.SetFlag(c.flags.Z, (tmp&0x00FF) == 0)
	c.SetFlag(c.flags.N, ((tmp&0x80)>>7) == 1)
	v := ((^(uint16(c.a)^uint16(c.fetched))&(uint16(c.a)^uint16(tmp)))&0x0080)>>7 == 1
	c.SetFlag(c.flags.V, v)
	c.a = uint8(tmp) & 0x00FF
	return 1
}

func (c *CPU6502) AND() uint8 {
	c.Fetch()
	c.a = c.a & c.fetched
	c.SetFlag(c.flags.Z, c.a == 0x00)
	c.SetFlag(c.flags.N, ((c.a&0x80)>>7) == 1)

	return 1
}

func (c *CPU6502) ASL() uint8 {
	return 0x00
}

func (c *CPU6502) BCC() uint8 {
	if c.GetFlag(c.flags.C) == 1 {
		c.cycles++
		c.addr_abs = c.pc + c.addr_rel

		if (c.addr_abs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}

		c.pc = c.addr_abs
	}
	return 0
}

func (c *CPU6502) BCS() uint8 {
	return 0x00
}

func (c *CPU6502) BEQ() uint8 {
	return 0x00
}

func (c *CPU6502) BIT() uint8 {
	return 0x00
}

func (c *CPU6502) BMI() uint8 {
	return 0x00
}

func (c *CPU6502) BNE() uint8 {
	return 0x00
}

func (c *CPU6502) BPL() uint8 {
	return 0x00
}

func (c *CPU6502) BRK() uint8 {
	return 0x00
}

func (c *CPU6502) BVC() uint8 {
	return 0x00
}

func (c *CPU6502) BVS() uint8 {
	return 0x00
}

func (c *CPU6502) CLC() uint8 {
	return 0x00
}

func (c *CPU6502) CLD() uint8 {
	return 0x00
}

func (c *CPU6502) CLI() uint8 {
	return 0x00
}

func (c *CPU6502) CLV() uint8 {
	return 0x00
}

func (c *CPU6502) CMP() uint8 {
	return 0x00
}

func (c *CPU6502) CPX() uint8 {
	return 0x00
}

func (c *CPU6502) CPY() uint8 {
	return 0x00
}

func (c *CPU6502) DEC() uint8 {
	return 0x00
}

func (c *CPU6502) DEX() uint8 {
	return 0x00
}

func (c *CPU6502) DEY() uint8 {
	return 0x00
}

func (c *CPU6502) EOR() uint8 {
	return 0x00
}

func (c *CPU6502) INC() uint8 {
	return 0x00
}

func (c *CPU6502) INX() uint8 {
	return 0x00
}

func (c *CPU6502) INY() uint8 {
	return 0x00
}

func (c *CPU6502) JMP() uint8 {
	return 0x00
}

func (c *CPU6502) JSR() uint8 {
	return 0x00
}

func (c *CPU6502) LDA() uint8 {
	return 0x00
}

func (c *CPU6502) LDX() uint8 {
	return 0x00
}

func (c *CPU6502) LDY() uint8 {
	return 0x00
}

func (c *CPU6502) LSR() uint8 {
	return 0x00
}

func (c *CPU6502) NOP() uint8 {
	return 0x00
}

func (c *CPU6502) ORA() uint8 {
	return 0x00
}

func (c *CPU6502) PHA() uint8 {
	c.Write(0x0100+uint16(c.stkp), c.a)
	c.stkp--
	return 0
}

func (c *CPU6502) PHP() uint8 {
	return 0x00
}

func (c *CPU6502) PLA() uint8 {
	c.stkp++
	c.a = c.Read(0x0100 + uint16(c.stkp))
	c.SetFlag(c.flags.Z, c.a == 0x00)
	c.SetFlag(c.flags.N, ((c.a&0x80)>>7) == 1)
	return 0
}

func (c *CPU6502) PLP() uint8 {
	return 0x00
}

func (c *CPU6502) ROL() uint8 {
	return 0x00
}

func (c *CPU6502) ROR() uint8 {
	return 0x00
}

func (c *CPU6502) RTI() uint8 {
	return 0x00
}

func (c *CPU6502) RTS() uint8 {
	return 0x00
}

func (c *CPU6502) SBC() uint8 {
	c.Fetch()

	tmp := uint16(c.a) + uint16(c.fetched) + uint16(c.GetFlag(c.flags.C))
	c.SetFlag(c.flags.C, tmp > 255)
	c.SetFlag(c.flags.Z, (tmp&0x00FF) == 0)
	c.SetFlag(c.flags.N, ((tmp&0x80)>>7) == 1)
	v := ((^(uint16(c.a)^uint16(c.fetched))&(uint16(c.a)^uint16(tmp)))&0x0080)>>7 == 1
	c.SetFlag(c.flags.V, v)
	c.a = uint8(tmp) & 0x00FF
	return 1
}

func (c *CPU6502) SEC() uint8 {
	return 0x00
}

func (c *CPU6502) SED() uint8 {
	return 0x00
}

func (c *CPU6502) SEI() uint8 {
	return 0x00
}

func (c *CPU6502) STA() uint8 {
	return 0x00
}

func (c *CPU6502) STX() uint8 {
	return 0x00
}

func (c *CPU6502) STY() uint8 {
	return 0x00
}

func (c *CPU6502) TAX() uint8 {
	return 0x00
}

func (c *CPU6502) TAY() uint8 {
	return 0x00
}

func (c *CPU6502) TSX() uint8 {
	return 0x00
}

func (c *CPU6502) TXA() uint8 {
	return 0x00
}

func (c *CPU6502) TXS() uint8 {
	return 0x00
}

func (c *CPU6502) TYA() uint8 {
	return 0x00
}

// Illegal opcode
func (c *CPU6502) XXX() uint8 {
	return 0x00
}
