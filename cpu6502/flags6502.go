package cpu6502

type FLAGS6502 struct {
	C uint8 // Carry bit
	Z uint8 // Zero
	I uint8 // Disable Interrupts
	D uint8 // Decimal Mode
	B uint8 // Break
	U uint8 // Unused
	V uint8 // Overflow
	N uint8 // Negative
}

func (f *FLAGS6502) Initialize() {
	f.C = (1 << 0)
	f.Z = (1 << 1)
	f.I = (1 << 2)
	f.D = (1 << 3)
	f.B = (1 << 4)
	f.U = (1 << 5)
	f.V = (1 << 6)
	f.N = (1 << 7)
}
