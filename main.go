package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	bus := Bus{}
	bus.Initialize()

	//runAllOpTests("./tests/")

	//fmt.Println(getFunNameAddr(bus.CPU.lookup[100].AddrMode))
	//return

	sdlController := SDLController{Bus: &bus}
	go loadROM("./nt.nes", &bus)
	//go runAssembly(&sdlController, &bus)

	if err := sdlController.Initialize(1100, 750, "cmgc"); err != nil {
		panic(err)
	}
}

func loadROM(path string, b *Bus) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	buf := make([]byte, 1000000)
	for {
		_, err := reader.Read(buf)

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
	}

	for i := 0; i < 0x4000; i++ {
		b.RAM[0x8000+i] = buf[0x0010+i]
		b.RAM[0xC000+i] = buf[0x0010+i]
	}
	b.CPU.disassemble(0x0000, 0xFFFF)
	b.CPU.Reset()
	b.CPU.pc = 0xC000

	//n5002 14572
	for n := 0; n < 6000; n++ {
		for {
			b.CPU.Clock()
			if b.CPU.cycles == 0 {
				break
			}
		}
	}
}

func run() (err error) {
	var window *sdl.Window
	var font *ttf.Font
	var surface *sdl.Surface
	var text *sdl.Surface

	if err = ttf.Init(); err != nil {
		return
	}
	defer ttf.Quit()

	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		return
	}
	defer sdl.Quit()

	// Create a window for us to draw the text on
	if window, err = sdl.CreateWindow("Drawing text", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN); err != nil {
		return
	}
	defer window.Destroy()

	if surface, err = window.GetSurface(); err != nil {
		return
	}

	// Load the font for our text
	if font, err = ttf.OpenFont("./pixel.ttf", 48); err != nil {
		return
	}
	defer font.Close()

	// Create a red text with the font
	if text, err = font.RenderUTF8Blended("Hello, World!", sdl.Color{R: 255, G: 0, B: 0, A: 255}); err != nil {
		return
	}
	defer text.Free()

	// Draw the text around the center of the window
	if err = text.Blit(nil, surface, &sdl.Rect{X: 400 - (text.W / 2), Y: 300 - (text.H / 2), W: 0, H: 0}); err != nil {
		return
	}

	// Update the window surface with what we have drawn
	window.UpdateSurface()

	// Run infinite loop until user closes the window
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		sdl.Delay(16)
	}

	return
}

func initSDL() {
	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	rect := sdl.Rect{X: 0, Y: 0, W: 200, H: 200}
	surface.FillRect(&rect, 0xffff0000)
	window.UpdateSurface()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false

			}
		}
	}
}

func runAssembly(con *SDLController, bus *Bus) {
	time.Sleep(time.Second)
	hex := []string{"A2", "0A", "8E", "00", "00", "A2", "03", "8E", "01", "00", "AC", "00", "00", "A9", "00", "18", "6D", "01", "00", "88", "D0", "FA", "8D", "02", "00", "EA", "EA", "EA"}
	offset := 0
	for _, h := range hex {
		value, _ := strconv.ParseInt(h, 16, 64)
		bus.RAM[0x8000+offset] = uint8(value)
		offset++
	}

	bus.CPU.disassemble(0x0000, 0xFFFF)
	bus.CPU.Reset()

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x80

	fmt.Println("Start PC, RAM[PC]")
	fmt.Println(bus.CPU.pc, bus.RAM[bus.CPU.pc])
	fmt.Println()

	for n := 0; n < 120; n++ {
		for {
			bus.CPU.Clock()
			if bus.CPU.cycles == 0 {
				break
			}
			time.Sleep(time.Millisecond * 100)

		}
	}

	/*bus.CPU.PrintRegisters()
	for i := 0; i < 10; i++ {
		fmt.Println(bus.RAM[i])
	}*/

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
