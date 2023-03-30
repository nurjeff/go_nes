package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"runtime"
	"strconv"
)

func main() {
	bus := Bus{}
	bus.Initialize()

	runAssembly()
	return
	runAllOpTests("./tests/")
}

func runAssembly() {
	bus := Bus{}
	bus.Initialize()

	hex := []string{"A2", "0A", "8E", "00", "00", "A2", "03", "8E", "01", "00", "AC", "00", "00", "A9", "00", "18", "6D", "01", "00", "88", "D0", "FA", "8D", "02", "00", "EA", "EA", "EA"}
	offset := 0
	for _, h := range hex {
		value, _ := strconv.ParseInt(h, 16, 64)
		bus.RAM[0x8000+offset] = uint8(value)
		offset++
	}

	fmt.Println(bus.RAM[0x8000 : len(hex)+0x8000])
	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x80

	bus.CPU.Reset()

	fmt.Println("Start PC, RAM[PC]")
	fmt.Println(bus.CPU.pc, bus.RAM[bus.CPU.pc])
	fmt.Println()

	for n := 0; n < 50; n++ {
		for {
			bus.CPU.Clock()
			if bus.CPU.cycles == 0 {
				bus.CPU.PrintRegisters()
				fmt.Println("----")
				fmt.Println()
				break
			}
		}

	}

	bus.CPU.PrintRegisters()
	for i := 0; i < 10; i++ {
		fmt.Println(bus.RAM[i])
	}

}

var offset = 16

func runAllOpTests(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	resmap := make(map[string][]string)
	for _, element := range files[offset : offset+1] {
		bus := Bus{}
		bus.Initialize()
		n := new(big.Int)
		n.SetString(element.Name()[:2], 16)
		opname := bus.CPU.lookup[int(n.Int64())].name
		if opname == "???" {
			continue
		}

		addrname := runtime.FuncForPC(reflect.ValueOf(bus.CPU.lookup[int(n.Int64())].AddrMode).Pointer()).Name()
		addrname = addrname[len(addrname)-6 : len(addrname)-3]

		failed := runOpTest(path+element.Name(), &bus)
		fstr := "PASS"
		if failed {
			fstr = "FAILED"
		}
		resmap[opname] = append(resmap[opname], addrname+" - "+fstr+"\n")

		fmt.Println("running test:", element.Name())
	}
	for k, v := range resmap {
		fmt.Println(k)
		fmt.Println(v)
		fmt.Println("---")
	}
}

func runOpTest(path string, bus *Bus) bool {
	tests := readOpTest(path)
	failedTests := []int{}
	for ind, element := range tests {
		corr := true
		for _, c := range element.CyclesRaw {
			tests[ind].Cycles = append(tests[ind].Cycles, Cycle{Address: uint16(c[0].(float64)), Value: uint8(c[1].(float64)), Op: fmt.Sprint(c[2])})
		}
		bus.CPU.Reset()
		bus.CPU.pc = element.Initial.PC
		bus.CPU.status = element.Initial.P
		bus.CPU.a = element.Initial.A
		bus.CPU.x = element.Initial.X
		bus.CPU.y = element.Initial.Y
		bus.CPU.stkp = element.Initial.S
		for _, r := range element.Initial.RAM {
			bus.RAM[r[0]] = uint8(r[1])
		}
		for _, _ = range tests[ind].Cycles {
			for {
				bus.CPU.Clock()
				if bus.CPU.cycles == 0 {
					break
				}
			}
		}

		for _, m := range element.Final.RAM {
			if bus.RAM[m[0]] != uint8(m[1]) {
				corr = false
			}
		}

		if !corr {
			failedTests = append(failedTests, ind)
		}
	}

	return len(failedTests) > 0
}

func readOpTest(path string) OpTests {
	optests := OpTests{}
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &optests); err != nil {
		panic(err)
	}
	return optests
}

type OpTests []OpTest

type OpTest struct {
	Name      string   `json:"name"`
	Initial   State    `json:"initial"`
	Final     State    `json:"final"`
	CyclesRaw CycleRaw `json:"cycles"`
	Cycles    []Cycle
}

type Cycle struct {
	Address uint16
	Value   uint8
	Op      string
}

type State struct {
	PC  uint16     `json:"pc"`
	S   uint8      `json:"s"`
	A   uint8      `json:"a"`
	X   uint8      `json:"x"`
	Y   uint8      `json:"y"`
	P   uint8      `json:"p"`
	RAM [][]uint16 `json:"ram"`
}

type CycleRaw [][]interface{}
