package main

import (
	"fmt"
	"reflect"
	"runtime"
)

func (c *CPU6502) disassemble(start uint16, stop uint16) {
	var addr uint32 = uint32(start)
	var value uint8 = 0x00
	var lo uint8 = 0x00
	var hi uint8 = 0x00

	var mapLines map[uint16]string = make(map[uint16]string)

	var line_addr uint16 = 0

	for addr <= uint32(stop) {
		line_addr = uint16(addr)
		var sInst string = "$" + hex(addr, 4) + ": "

		var opcode uint8 = c.Read(uint16(addr))
		addr++
		sInst += c.lookup[opcode].name + " "

		if getFunNameAddr(c.lookup[opcode].AddrMode) == "IMP" {
			sInst += " {IMP}"
		} else if getFunNameAddr(c.lookup[opcode].AddrMode) == "IMM" {

			value = c.Read(uint16(addr))
			sInst += "#$ " + hex(value, 2) + " {IMM}"

		} else if getFunNameAddr(c.lookup[opcode].AddrMode) == "ZP0" {
			lo = c.Read(uint16(addr))
			addr++
			hi = 0x00
			sInst += "$ " + hex(lo, 2) + " {ZP0}"
		} else if getFunNameAddr(c.lookup[opcode].AddrMode) == "ZPX" {
			lo = c.Read(uint16(addr))
			addr++
			hi = 0x00
			sInst += "$ " + hex(lo, 2) + ", X {ZPX}"
		} else if getFunNameAddr(c.lookup[opcode].AddrMode) == "ZPY" {
			lo = c.Read(uint16(addr))
			addr++
			hi = 0x00
			sInst += "$ " + hex(lo, 2) + ", Y {ZPY}"
		} else if getFunNameAddr(c.lookup[opcode].AddrMode) == "IZX" {
			lo = c.Read(uint16(addr))
			addr++
			hi = 0x00
			sInst += "($ " + hex(lo, 2) + "), X {IZX}"
		} else if getFunNameAddr(c.lookup[opcode].AddrMode) == "IZY" {
			lo = c.Read(uint16(addr))
			addr++
			hi = 0x00
			sInst += "($ " + hex(lo, 2) + "), Y {IZY}"
		} else if getFunNameAddr(c.lookup[opcode].AddrMode) == "ABS" {
			lo = c.Read(uint16(addr))
			addr++
			hi = c.Read(uint16(addr))
			addr++
			sInst += "$ " + hex((uint16(hi)<<8)|uint16(lo), 4) + " {ABS}"
		} else if getFunNameAddr(c.lookup[opcode].AddrMode) == "ABX" {
			lo = c.Read(uint16(addr))
			addr++
			hi = c.Read(uint16(addr))
			addr++
			sInst += "$ " + hex((uint16(hi)<<8)|uint16(lo), 4) + ", X {ABX}"
		} else if getFunNameAddr(c.lookup[opcode].AddrMode) == "ABY" {
			lo = c.Read(uint16(addr))
			addr++
			hi = c.Read(uint16(addr))
			addr++
			sInst += "$ " + hex((uint16(hi)<<8)|uint16(lo), 4) + ", Y {ABY}"
		} else if getFunNameAddr(c.lookup[opcode].AddrMode) == "IND" {
			lo = c.Read(uint16(addr))
			addr++
			hi = c.Read(uint16(addr))
			addr++
			sInst += "($ " + hex((uint16(hi)<<8)|uint16(lo), 4) + ") {IND}"
		} else if getFunNameAddr(c.lookup[opcode].AddrMode) == "REL" {
			value = c.Read(uint16(addr))
			addr++
			sInst += "$ " + hex(value, 2) + " [$" + hex(addr+uint32(value), 4) + "] {REL}"
		} else {
			fmt.Println("Unknown Addressing Mode?", getFunNameAddr(c.lookup[opcode].AddrMode))
		}

		mapLines[line_addr] = sInst
	}

	c.Disassembly = mapLines
}

func hex(variable interface{}, n int) string {
	h := fmt.Sprintf("%x", variable)
	if len(h) > n {
		h = h[:n]
	}
	return h
}

func getFunName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func getFunNameAddr(f interface{}) string {
	n := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	n = n[len(n)-6 : len(n)-3]
	return n
}
