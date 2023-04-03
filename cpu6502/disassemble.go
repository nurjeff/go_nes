package cpu6502

import (
	t "github.com/sc-js/go_nes/emutools"
)

func (c *CPU6502) Disassemble(start uint16, stop uint16) {
	var addr uint32 = uint32(start)
	var value uint8 = 0x00
	var lo uint8 = 0x00
	var hi uint8 = 0x00

	var mapLines map[uint16]string = make(map[uint16]string)

	var line_addr uint16 = 0

	for addr <= uint32(stop) {
		line_addr = uint16(addr)
		var sInst string = "$" + t.Hex(addr, 4) + ": "

		var opcode uint8 = c.Read(uint16(addr))
		addr++
		sInst += c.lookup[opcode].name + " "

		fname := t.GetFunNameAddr(c.lookup[opcode].AddrMode)

		switch fname {
		case "IMP":
			sInst += " {IMP}"
		case "IMM":
			value = c.Read(uint16(addr))
			addr++
			sInst += "#$ " + t.Hex(value, 2) + " {IMM}"
		case "ZP0":
			lo = c.Read(uint16(addr))
			addr++
			hi = 0x00
			sInst += "$ " + t.Hex(lo, 2) + " {ZP0}"
		case "ZPX":
			lo = c.Read(uint16(addr))
			addr++
			hi = 0x00
			sInst += "$ " + t.Hex(lo, 2) + ", X {ZPX}"
		case "ZPY":
			lo = c.Read(uint16(addr))
			addr++
			hi = 0x00
			sInst += "$ " + t.Hex(lo, 2) + ", Y {ZPY}"
		case "IZX":
			lo = c.Read(uint16(addr))
			addr++
			hi = 0x00
			sInst += "($ " + t.Hex(lo, 2) + "), X {IZX}"
		case "IZY":
			lo = c.Read(uint16(addr))
			addr++
			hi = 0x00
			sInst += "($ " + t.Hex(lo, 2) + "), Y {IZY}"
		case "ABS":
			lo = c.Read(uint16(addr))
			addr++
			hi = c.Read(uint16(addr))
			addr++
			sInst += "$ " + t.Hex((uint16(hi)<<8)|uint16(lo), 4) + " {ABS}"
		case "ABX":
			lo = c.Read(uint16(addr))
			addr++
			hi = c.Read(uint16(addr))
			addr++
			sInst += "$ " + t.Hex((uint16(hi)<<8)|uint16(lo), 4) + ", X {ABX}"
		case "ABY":
			lo = c.Read(uint16(addr))
			addr++
			hi = c.Read(uint16(addr))
			addr++
			sInst += "$ " + t.Hex((uint16(hi)<<8)|uint16(lo), 4) + ", Y {ABY}"
		case "IND":
			lo = c.Read(uint16(addr))
			addr++
			hi = c.Read(uint16(addr))
			addr++
			sInst += "($ " + t.Hex((uint16(hi)<<8)|uint16(lo), 4) + ") {IND}"
		case "REL":
			value = c.Read(uint16(addr))
			addr++
			sInst += "$ " + t.Hex(int(value), 2) + " [$" + t.Hex(addr+uint32(value), 4) + "] {REL}"
		default:
			panic("unknown addressing mode: " + fname)
		}
		mapLines[line_addr] = sInst
	}

	c.Disassembly = mapLines
}
