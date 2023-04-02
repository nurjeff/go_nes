package cpu6502

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/sc-js/go_nes/emutools"
)

func (c *CPU6502) ADC() uint8 {
	c.Fetch()

	c.tmp = uint16(c.a) + uint16(c.fetched) + uint16(c.GetFlag(c.Flags.C))

	c.SetFlag(c.Flags.C, c.tmp > 255)

	c.SetFlag(c.Flags.Z, (c.tmp&0x00FF) == 0)

	c.SetFlag(c.Flags.N, ((c.tmp&0x80)>>7) >= 1)

	v := (((^(uint16(c.a) ^ uint16(c.fetched)) & (uint16(c.a) ^ uint16(c.tmp))) & 0x0080) >> 7) >= 1

	c.SetFlag(c.Flags.V, v)
	c.a = uint8(c.tmp) & 0x00FF
	return 1
}

func (c *CPU6502) AND() uint8 {
	c.Fetch()
	c.a = c.a & c.fetched
	c.SetFlag(c.Flags.Z, c.a == 0x00)
	c.SetFlag(c.Flags.N, ((c.a&0x80)>>7) >= 1)

	return 1
}

func (c *CPU6502) ASL() uint8 {
	c.Fetch()
	c.tmp = uint16(c.fetched) << 1
	c.SetFlag(c.Flags.C, (c.tmp&0xFF00) > 0)
	c.SetFlag(c.Flags.Z, (c.tmp&0x00FF) == 0)
	c.SetFlag(c.Flags.N, (c.tmp&0x80)>>7 >= 1)
	if runtime.FuncForPC(reflect.ValueOf(c.lookup[c.opcode].AddrMode).Pointer()).Name() == runtime.FuncForPC(reflect.ValueOf(c.IMP).Pointer()).Name() {
		c.a = uint8(c.tmp & 0x00FF)
	} else {
		c.Write(c.addr_abs, uint8(c.tmp&0x00FF))
	}
	return 0
}

func (c *CPU6502) BCC() uint8 {
	if c.GetFlag(c.Flags.C) == 0 {
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
	if c.GetFlag(c.Flags.C) == 1 {
		c.cycles++
		c.addr_abs = c.pc + c.addr_rel

		if (c.addr_abs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}

		c.pc = c.addr_abs
	}
	return 0x00
}

func (c *CPU6502) BEQ() uint8 {
	if c.GetFlag(c.Flags.Z) == 1 {
		c.cycles++
		c.addr_abs = c.pc + c.addr_rel

		if (c.addr_abs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}

		c.pc = c.addr_abs
	}
	return 0x00
}

func (c *CPU6502) BIT() uint8 {
	c.Fetch()
	c.tmp = uint16(c.a) & uint16(c.fetched)
	c.SetFlag(c.Flags.Z, (c.tmp&0x00FF) == 0x00)
	c.SetFlag(c.Flags.N, ((c.fetched)&(1<<7)) >= 1)
	c.SetFlag(c.Flags.V, ((c.fetched)&(1<<6)) >= 1)
	return 0x00
}

func (c *CPU6502) BMI() uint8 {
	if c.GetFlag(c.Flags.N) == 1 {
		c.cycles++
		c.addr_abs = c.pc + c.addr_rel

		if (c.addr_abs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}

		c.pc = c.addr_abs
	}
	return 0x00
}

func (c *CPU6502) BNE() uint8 {
	if c.GetFlag(c.Flags.Z) == 0 {
		c.cycles++
		c.addr_abs = c.pc + c.addr_rel

		if (c.addr_abs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}

		c.pc = c.addr_abs
	}
	return 0x00
}

func (c *CPU6502) BPL() uint8 {
	if c.GetFlag(c.Flags.N) == 0 {
		c.cycles++
		c.addr_abs = c.pc + c.addr_rel

		if (c.addr_abs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}

		c.pc = c.addr_abs
	}
	return 0x00
}

func (c *CPU6502) BRK() uint8 {
	c.pc++

	c.SetFlag(c.Flags.I, true)
	c.Write(0x0100+uint16(c.stkp), uint8((c.pc>>8)&0x00FF))
	c.stkp--
	c.Write(0x0100+uint16(c.stkp), uint8(c.pc&0x00FF))
	c.stkp--

	c.SetFlag(c.Flags.B, true)
	c.Write(0x0100+uint16(c.stkp), c.status)
	c.stkp--
	c.SetFlag(c.Flags.B, false)

	c.pc = uint16(c.Read(0xFFFE)) | (uint16(c.Read(0xFFFF)) << 8)
	return 0x00
}

func (c *CPU6502) BVC() uint8 {
	if c.GetFlag(c.Flags.V) == 0 {
		c.cycles++
		c.addr_abs = c.pc + c.addr_rel

		if (c.addr_abs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}

		c.pc = c.addr_abs
	}
	return 0x00
}

func (c *CPU6502) BVS() uint8 {
	if c.GetFlag(c.Flags.V) == 1 { // 1
		c.cycles++
		c.addr_abs = c.pc + c.addr_rel

		if (c.addr_abs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}

		c.pc = c.addr_abs
	}
	return 0x00
}

func (c *CPU6502) CLC() uint8 {
	c.SetFlag(c.Flags.C, false)
	return 0x00
}

func (c *CPU6502) CLD() uint8 {
	c.SetFlag(c.Flags.D, false)
	return 0x00
}

func (c *CPU6502) CLI() uint8 {
	c.SetFlag(c.Flags.I, false)
	return 0x00
}

func (c *CPU6502) CLV() uint8 {
	c.SetFlag(c.Flags.V, false)
	return 0x00
}

func (c *CPU6502) CMP() uint8 {
	c.Fetch()
	c.tmp = uint16(c.a) - uint16(c.fetched)
	c.SetFlag(c.Flags.C, c.a >= c.fetched)
	c.SetFlag(c.Flags.Z, (c.tmp&0x00FF) == 0x0000)
	c.SetFlag(c.Flags.N, (c.tmp&0x0080)>>7 >= 1)
	return 1
}

func (c *CPU6502) CPX() uint8 {
	c.Fetch()
	c.tmp = uint16(c.x) - uint16(c.fetched)
	c.SetFlag(c.Flags.C, c.x >= c.fetched)
	c.SetFlag(c.Flags.Z, (c.tmp&0x00FF) == 0x0000)
	c.SetFlag(c.Flags.N, (c.tmp&0x0080)>>7 >= 1)
	return 0x00
}

func (c *CPU6502) CPY() uint8 {
	c.Fetch()
	c.tmp = uint16(c.y) - uint16(c.fetched)
	c.SetFlag(c.Flags.C, c.y >= c.fetched)
	c.SetFlag(c.Flags.Z, (c.tmp&0x00FF) == 0x0000)
	c.SetFlag(c.Flags.N, (c.tmp&0x0080)>>7 >= 1)
	return 0x00
}

func (c *CPU6502) ISC() uint8 {
	c.INC()
	c.SBC()
	return 0
}

func (c *CPU6502) ALR() uint8 {
	c.AND()
	c.LSR()
	return 0
}

func (c *CPU6502) SLO() uint8 {
	c.ASL()
	c.ORA()
	return 0
}

func (c *CPU6502) RLA() uint8 {
	c.ROL()
	c.AND()
	return 0
}

func (c *CPU6502) SRE() uint8 {
	c.LSR()
	c.EOR()
	return 0
}

func (c *CPU6502) RRA() uint8 {
	c.ROR()
	c.ADC()
	return 0
}

func (c *CPU6502) DEC() uint8 {
	c.Fetch()
	c.tmp = uint16(c.fetched) - 1
	c.Write(c.addr_abs, uint8(c.tmp&0x00FF))
	c.SetFlag(c.Flags.Z, (c.tmp&0x00FF) == 0x0000)
	c.SetFlag(c.Flags.N, (c.tmp&0x0080)>>7 >= 1)
	return 0x00
}

func (c *CPU6502) DEX() uint8 {
	c.x--
	c.SetFlag(c.Flags.Z, c.x == 0x00)
	c.SetFlag(c.Flags.N, ((c.x&0x80)>>7) >= 1)
	return 0x00
}

func (c *CPU6502) DEY() uint8 {
	c.y--
	c.SetFlag(c.Flags.Z, c.y == 0x00)
	c.SetFlag(c.Flags.N, (c.y&0x80)>>7 >= 1)
	return 0x00
}

func (c *CPU6502) EOR() uint8 {
	c.Fetch()
	c.a = c.a ^ c.fetched
	c.SetFlag(c.Flags.Z, c.a == 0x00)
	c.SetFlag(c.Flags.N, (c.a&0x80)>>7 >= 1)
	return 0x00
}

func (c *CPU6502) INC() uint8 {
	c.Fetch()
	c.tmp = uint16(c.fetched) + 1
	c.Write(c.addr_abs, uint8(c.tmp&0x00FF))
	c.SetFlag(c.Flags.Z, (c.tmp&0x00FF) == 0x0000)
	c.SetFlag(c.Flags.N, (c.tmp&0x0080)>>7 >= 1)
	return 0x00
}

func (c *CPU6502) INX() uint8 {
	c.x++
	c.SetFlag(c.Flags.Z, c.x == 0x00)
	c.SetFlag(c.Flags.N, ((c.x&0x80)>>7) >= 1)
	return 0x00
}

func (c *CPU6502) INY() uint8 {
	c.y++
	c.SetFlag(c.Flags.Z, c.y == 0x00)
	c.SetFlag(c.Flags.N, (c.y&0x80)>>7 >= 1)
	return 0x00
}

func (c *CPU6502) JMP() uint8 {
	c.pc = c.addr_abs
	return 0x00
}

func (c *CPU6502) JSR() uint8 {
	c.pc--
	c.Write(0x0100+uint16(c.stkp), (uint8((c.pc >> 8) & 0x00FF)))
	c.stkp--
	c.Write(0x0100+uint16(c.stkp), uint8(c.pc&0x00FF))
	c.stkp--

	c.pc = c.addr_abs
	return 0x00
}

func (c *CPU6502) LDA() uint8 {
	c.Fetch()
	c.a = c.fetched
	c.SetFlag(c.Flags.Z, c.a == 0x00)
	c.SetFlag(c.Flags.N, (c.a&0x80)>>7 >= 1)
	return 1
}

func (c *CPU6502) LDX() uint8 {
	c.Fetch()
	c.x = c.fetched
	c.SetFlag(c.Flags.Z, c.x == 0x00)
	c.SetFlag(c.Flags.N, (c.x&0x80)>>7 >= 1)
	return 1
}

func (c *CPU6502) LAX() uint8 {
	c.Fetch()
	c.x = c.fetched
	c.a = c.fetched
	c.SetFlag(c.Flags.Z, c.x == 0x00)
	c.SetFlag(c.Flags.N, (c.x&0x80)>>7 >= 1)
	return 1
}

func (c *CPU6502) DCP() uint8 {
	c.DEC()
	c.CMP()
	return 0
}

func (c *CPU6502) LDY() uint8 {
	c.Fetch()
	c.y = c.fetched
	c.SetFlag(c.Flags.Z, c.y == 0x00)
	c.SetFlag(c.Flags.N, (c.y&0x80)>>7 >= 1)
	return 1
}

func (c *CPU6502) LSR() uint8 {
	c.Fetch()
	c.SetFlag(c.Flags.C, c.fetched&0x0001 >= 1)
	c.tmp = uint16(c.fetched) >> 1
	c.SetFlag(c.Flags.Z, (c.tmp&0x00FF) == 0x0000)
	c.SetFlag(c.Flags.N, (c.tmp&0x0080)>>7 >= 1)
	if runtime.FuncForPC(reflect.ValueOf(c.lookup[c.opcode].AddrMode).Pointer()).Name() == runtime.FuncForPC(reflect.ValueOf(c.IMP).Pointer()).Name() {
		c.a = uint8(c.tmp & 0x00FF)
	} else {
		c.Write(c.addr_abs, uint8(c.tmp&0x00FF))
	}

	return 0x00
}

func (c *CPU6502) NOP() uint8 {
	switch c.opcode {
	case 0x1C:
		return 1
	case 0x3C:
		return 1
	case 0x5C:
		return 1
	case 0x7C:
		return 1
	case 0xDC:
		return 1
	case 0xFC:
		return 1
	}
	return 0
}

func (c *CPU6502) ORA() uint8 {
	c.Fetch()
	c.a = c.a | c.fetched
	c.SetFlag(c.Flags.Z, c.a == 0x00)
	c.SetFlag(c.Flags.N, (c.a&0x80)>>7 >= 1)
	return 0x00
}

func (c *CPU6502) PHA() uint8 {
	c.Write(0x0100+uint16(c.stkp), c.a)
	c.stkp--
	return 0
}

func (c *CPU6502) PHP() uint8 {
	c.Write(0x100+uint16(c.stkp), c.status|(1<<4)|(1<<5))
	c.SetFlag(c.Flags.B, false)
	c.SetFlag(c.Flags.U, false)
	c.stkp--
	return 0x00
}

func (c *CPU6502) PLA() uint8 {
	c.stkp++
	c.a = c.Read(0x0100 + uint16(c.stkp))
	c.SetFlag(c.Flags.Z, c.a == 0x00)
	c.SetFlag(c.Flags.N, ((c.a&0x80)>>7) >= 1)
	return 0
}

func (c *CPU6502) PLP() uint8 {
	c.stkp++
	b := c.GetFlag(c.Flags.B)
	c.status = c.Read(0x0100 + uint16(c.stkp))
	c.SetFlag(c.Flags.B, b >= 1)
	return 0x00
}

func (c *CPU6502) ROL() uint8 {
	c.Fetch()
	c.tmp = (uint16(c.fetched) << 1) | uint16(c.GetFlag(c.Flags.C))
	c.SetFlag(c.Flags.C, (c.tmp&0xFF00) >= 1)
	c.SetFlag(c.Flags.Z, (c.tmp&0x00FF) == 0x0000)
	c.SetFlag(c.Flags.N, ((c.tmp&0x0080)>>7) >= 1)
	if emutools.GetFunNameAddr(c.lookup[c.opcode].AddrMode) == "IMP" {
		c.a = uint8(c.tmp & 0x00FF)
	} else {
		c.Write(c.addr_abs, uint8(c.tmp&0x00FF))
	}

	return 0x00
}

func (c *CPU6502) ROR() uint8 {
	c.Fetch()
	c.tmp = (uint16(c.GetFlag(c.Flags.C)) << 7) | (uint16(c.fetched) >> 1)
	c.SetFlag(c.Flags.C, (c.fetched&0x01) > 0)
	c.SetFlag(c.Flags.Z, (c.tmp&0x00FF) == 0x00)
	c.SetFlag(c.Flags.N, ((c.tmp&0x0080)>>7) >= 1)
	if emutools.GetFunNameAddr(c.lookup[c.opcode].AddrMode) == "IMP" {
		c.a = uint8(c.tmp & 0x00FF)
	} else {
		c.Write(c.addr_abs, uint8(c.tmp&0x00FF))
	}

	return 0x00
}

func (c *CPU6502) RTI() uint8 {
	c.stkp++
	c.status = c.Read(0x0100 + uint16(c.stkp))
	b := ^(1 << 4)
	u := ^(1 << 5)
	c.status &= uint8(b)
	c.status &= uint8(u)

	c.stkp++
	c.pc = uint16(c.Read(0x0100 + uint16(c.stkp)))
	c.stkp++
	c.pc |= uint16(c.Read(0x0100+uint16(c.stkp))) << 8
	return 0x00
}

func (c *CPU6502) RTS() uint8 {
	c.stkp++
	c.pc = uint16(c.Read(0x0100 + uint16(c.stkp)))
	c.stkp++
	c.pc |= (uint16(c.Read(0x100+uint16(c.stkp))) << 8)
	c.pc++
	return 0x00
}

func (c *CPU6502) SBC() uint8 {
	c.Fetch()

	var value uint16 = uint16(c.fetched) ^ 0x00FF

	c.tmp = uint16(c.a) + value + uint16(c.GetFlag(c.Flags.C))

	c.SetFlag(c.Flags.C, c.tmp > 255)

	c.SetFlag(c.Flags.Z, (c.tmp&0x00FF) == 0)

	c.SetFlag(c.Flags.N, ((c.tmp&0x80)>>7) >= 1)

	v := (c.tmp^uint16(c.a))&((c.tmp^value)&0x0080) >= 1

	c.SetFlag(c.Flags.V, v)
	c.a = uint8(c.tmp) & 0x00FF
	return 1

	/*c.tmp = uint16(c.a) + uint16(value) + uint16(c.GetFlag(c.Flags.C))
	c.SetFlag(c.Flags.C, c.tmp > 255)
	c.SetFlag(c.Flags.Z, (c.tmp&0x00FF) == 0)
	c.SetFlag(c.Flags.N, ((c.tmp&0x80)>>7) >= 1)
	v := (((^(uint16(c.a) ^ uint16(c.fetched)) & (uint16(c.a) ^ uint16(c.tmp))) & 0x0080) >> 7) >= 1
	c.SetFlag(c.Flags.V, v)
	c.a = uint8(c.tmp) & 0x00FF

	return 1*/
}

func (c *CPU6502) SEC() uint8 {
	c.SetFlag(c.Flags.C, true)
	return 0x00
}

//need 10100101
//got  11101101
//wrong D, V

func (c *CPU6502) SED() uint8 {
	c.SetFlag(c.Flags.D, true)
	return 0x00
}

func (c *CPU6502) SEI() uint8 {
	c.SetFlag(c.Flags.I, true)
	return 0x00
}

func (c *CPU6502) SAX() uint8 {
	c.Write(c.addr_abs, c.a&c.x)
	return 0x00
}

func (c *CPU6502) STA() uint8 {
	c.Write(c.addr_abs, c.a)
	return 0x00
}

func (c *CPU6502) STX() uint8 {
	c.Write(c.addr_abs, c.x)
	return 0x00
}

func (c *CPU6502) STY() uint8 {
	c.Write(c.addr_abs, c.y)
	return 0x00
}

func (c *CPU6502) TAX() uint8 {
	c.x = c.a
	c.SetFlag(c.Flags.Z, c.x == 0x00)
	c.SetFlag(c.Flags.N, (c.x&0x80)>>7 >= 1)
	return 0x00
}

func (c *CPU6502) TAY() uint8 {
	c.y = c.a
	c.SetFlag(c.Flags.Z, c.y == 0x00)
	c.SetFlag(c.Flags.N, (c.y&0x80)>>7 >= 1)
	return 0x00
}

func (c *CPU6502) TSX() uint8 {
	c.x = c.stkp
	c.SetFlag(c.Flags.Z, c.x == 0x00)
	c.SetFlag(c.Flags.N, (c.x&0x80)>>7 >= 1)
	return 0x00
}

func (c *CPU6502) TXA() uint8 {
	c.a = c.x
	c.SetFlag(c.Flags.Z, c.a == 0x00)
	c.SetFlag(c.Flags.N, (c.a&0x80)>>7 >= 1)
	return 0x00
}

func (c *CPU6502) TXS() uint8 {
	c.stkp = c.x
	return 0x00
}

func (c *CPU6502) TYA() uint8 {
	c.a = c.y
	c.SetFlag(c.Flags.Z, c.a == 0x00)
	c.SetFlag(c.Flags.N, (c.a&0x80)>>7 >= 1)
	return 0x00
}

// Not implemented opcodes
func (c *CPU6502) XXX() uint8 {
	fmt.Println("OpCode", emutools.Hex(c.opcode, 2), "not implemented")
	return 0x00
}

func (c *CPU6502) IRQ() {
	if c.GetFlag(c.Flags.I) == 0 {
		c.Write(0x0100+uint16(c.stkp), uint8((c.pc>>8)&0x00FF))
		c.stkp--
		c.Write(0x0100+uint16(c.stkp), uint8(c.pc&0x00FF))
		c.stkp--

		c.SetFlag(c.Flags.B, false)
		c.SetFlag(c.Flags.U, true)
		c.SetFlag(c.Flags.I, true)
		c.Write(0x0100+uint16(c.stkp), c.status)
		c.stkp--

		c.addr_abs = 0xFFFE
		var lo uint16 = uint16(c.Read(c.addr_abs + 0))
		var hi uint16 = uint16(c.Read(c.addr_abs + 1))
		c.pc = (hi >> 8) | lo

		c.cycles = 7
	}
}

func (c *CPU6502) NMI() {
	c.Write(0x0100+uint16(c.stkp), uint8((c.pc>>8)&0x00FF))
	c.stkp--
	c.Write(0x0100+uint16(c.stkp), uint8(c.pc&0x00FF))
	c.stkp--

	c.SetFlag(c.Flags.B, false)
	c.SetFlag(c.Flags.U, true)
	c.SetFlag(c.Flags.I, true)
	c.Write(0x0100+uint16(c.stkp), c.status)
	c.stkp--

	c.addr_abs = 0xFFFE
	var lo uint16 = uint16(c.Read(c.addr_abs + 0))
	var hi uint16 = uint16(c.Read(c.addr_abs + 1))
	c.pc = (hi >> 8) | lo

	c.cycles = 8
}
