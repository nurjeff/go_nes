package cpu6502

import (
	"reflect"
	"runtime"
)

type CPU6502 struct {
	//Bus   *Bus      // Bus connection pointer
	Flags FLAGS6502 // 6502 Flags

	ReadBus  func(uint16, bool) uint8
	WriteBus func(uint16, uint8)

	// Registers
	status uint8  // Status
	a      uint8  // Accumulator
	x      uint8  // X
	y      uint8  // Y
	stkp   uint8  // Stack pointer
	pc     uint16 // Program counter

	// Data
	fetched uint8

	// Addressing
	addr_abs uint16
	addr_rel uint16
	opcode   uint8
	cycles   uint8

	// Temporary variable used by ops to prevent constant allocation/deallocation
	tmp uint16

	// Total instructions
	totalInstructions uint64

	// Total cycles taken
	totalCycles uint64

	// Instruction lookup table
	lookup Instructions

	// Map to hold disassembled instructions - for display
	Disassembly map[uint16]string
}

func (c *CPU6502) GetCycles() uint8 {
	return c.cycles
}

// Setup the CPU
func (c *CPU6502) Initialize() {
	c.Flags = FLAGS6502{}
	c.Flags.Initialize()
	c.InitInternals()
}

func (c *CPU6502) InitInternals() {
	// Init Registers
	c.status = 0x00
	c.a = 0x00
	c.x = 0x00
	c.y = 0x00
	c.stkp = 0x00
	c.pc = 0x0000

	// Init Data
	c.fetched = 0x00
	c.addr_abs = 0x0000
	c.addr_rel = 0x0000
	c.opcode = 0x00
	c.cycles = 0

	c.lookup = Instructions{}
	c.lookup = *c.lookup.Fill(c)
}

// Read from the bus
func (c *CPU6502) Read(addr uint16) uint8 {
	return c.ReadBus(addr, false)
}

// Write to the bus
func (c *CPU6502) Write(addr uint16, data uint8) {
	c.WriteBus(addr, data)
}

func (c *CPU6502) Clock() {
	if c.cycles == 0 {
		c.opcode = c.ReadBus(c.pc, false)
		c.pc++

		//fmt.Println(c.lookup[c.opcode].name)

		c.cycles = c.lookup[c.opcode].cycles
		addrCycle := c.lookup[c.opcode].AddrMode()
		opCycle := c.lookup[c.opcode].OpCode()

		c.cycles += (addrCycle & opCycle)

		c.SetFlag(c.Flags.U, true)
		c.totalInstructions++
	}
	c.cycles--
	c.totalCycles++
}

func (c *CPU6502) Fetch() uint8 {
	if runtime.FuncForPC(reflect.ValueOf(c.lookup[c.opcode].AddrMode).Pointer()).Name() != runtime.FuncForPC(reflect.ValueOf(c.IMP).Pointer()).Name() {
		//fmt.Println("ADDR HERE:", emutools.Hex(c.addr_abs, 4))
		c.fetched = c.ReadBus(c.addr_abs, false)
	}
	return c.fetched
}

func (c *CPU6502) GetFlag(f uint8) uint8 {
	if (c.status & f) > 0 {
		return 1
	} else {
		return 0
	}
}

func (c *CPU6502) SetFlag(f uint8, v bool) {
	if v {
		c.status |= f
	} else {
		c.status &= ^f
	}
}

func (c *CPU6502) Reset() {
	c.a = 0
	c.x = 0
	c.y = 0
	c.stkp = 0xFD
	c.status = 0x00 | (1 << 5)

	c.addr_abs = 0xFFFC
	var lo uint16 = uint16(c.Read(c.addr_abs + 0))
	var hi uint16 = uint16(c.Read(c.addr_abs + 1))

	c.pc = (hi << 8) | lo
	c.addr_rel = 0x0000
	c.addr_abs = 0x0000
	c.fetched = 0x00

	c.cycles = 8
}

func (c *CPU6502) A() uint8 {
	return c.a
}

func (c *CPU6502) X() uint8 {
	return c.x
}

func (c *CPU6502) Y() uint8 {
	return c.y
}

func (c *CPU6502) P() uint8 {
	return c.status
}

func (c *CPU6502) SP() uint8 {
	return c.stkp
}

func (c *CPU6502) PC() uint16 {
	return c.pc
}

func (c *CPU6502) TC() uint64 {
	return c.totalCycles
}

func (c *CPU6502) SetPC(pc uint16) {
	c.pc = pc
}

func (c *CPU6502) GetOP() Instruction {
	return c.lookup[c.pc]
}
