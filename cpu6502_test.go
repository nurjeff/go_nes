package main

import "testing"

// check correct accumulator
func TestADCAdd(t *testing.T) {
	bus := Bus{}
	bus.Initialize()
	bus.CPU.a = 10
	bus.Write(0, 100)
	bus.CPU.ADC()

	if bus.CPU.a != 110 || bus.CPU.GetFlag(bus.CPU.flags.C) > 0 {
		t.Errorf("ADC err")
	}
}

// check carry flag C
func TestADCCarry(t *testing.T) {
	bus := Bus{}
	bus.Initialize()
	bus.CPU.a = 126
	bus.Write(0, 200)
	bus.CPU.ADC()

	if bus.CPU.a != 70 || bus.CPU.GetFlag(bus.CPU.flags.C) != 1 {
		t.Errorf("ADC err")
	}
}

// check overflow flag V
func TestADCOverflow(t *testing.T) {
	bus := Bus{}
	bus.Initialize()
	bus.CPU.a = 0b10011111
	bus.Write(0, 200)
	bus.CPU.ADC()

	if bus.CPU.a != 103 || bus.CPU.GetFlag(bus.CPU.flags.V) != 1 || bus.CPU.GetFlag(bus.CPU.flags.C) != 1 {
		t.Errorf("ADC err")
	}
}

// Check negative flag N
func TestADCNegative(t *testing.T) {
	bus := Bus{}
	bus.Initialize()
	bus.CPU.a = 0b10011111
	bus.Write(0, 0x0001)
	bus.CPU.ADC()

	if bus.CPU.a != 160 || bus.CPU.GetFlag(bus.CPU.flags.N) != 1 {
		t.Errorf("ADC err")
	}
}

func TestSBCSubtract(t *testing.T) {
	bus := Bus{}
	bus.Initialize()

	bus.CPU.a = 100
	bus.Write(0, 104)
	bus.CPU.SBC()

	if bus.CPU.a != 251 || bus.CPU.GetFlag(bus.CPU.flags.N) != 1 || bus.CPU.GetFlag(bus.CPU.flags.V) != 1 {
		t.Errorf("ADC err")
	}
}
