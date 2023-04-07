package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/sc-js/go_nes/bus"
	"github.com/sc-js/go_nes/cartridge"
	"github.com/sc-js/go_nes/window"
)

func main() {
	cartridge := cartridge.Cartridge{}
	cartridge.Initialize("./assets/ducktales.nes")

	bus := bus.Bus{}
	bus.InsertCartidge(&cartridge)

	sdlController := window.SDLController{Bus: &bus}
	bus.Initialize(&sdlController.ScreenTransfer)

	go bus.Boot()
	//go loadROM("./nt.nes", &bus)

	if err := sdlController.Initialize(256*window.UPRES, 240*window.UPRES, "cmgc"); err != nil {
		panic(err)
	}

	sdlController.Run()
}

func loadROM(path string, b *bus.Bus) {
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
		b.CPURAM[0x8000+i] = buf[0x0010+i]
		b.CPURAM[0xC000+i] = buf[0x0010+i]
	}
	b.CPU.Disassemble(0x8000, 0xFFFF)
	b.CPU.Reset()
	b.CPU.SetPC(0xC000)

	for n := 0; n < 0; n++ {
		for {
			b.CPU.Clock()
			if b.CPU.GetCycles() == 0 {
				break
			}
		}
	}
}

/*
func runAssembly(con *window.SDLController, bus *bus.Bus) {
	time.Sleep(time.Second)
	hex := []string{"A2", "0A", "8E", "00", "00", "A2", "03", "8E", "01", "00", "AC", "00", "00", "A9", "00", "18", "6D", "01", "00", "88", "D0", "FA", "8D", "02", "00", "EA", "EA", "EA"}
	offset := 0
	for _, h := range hex {
		value, _ := strconv.ParseInt(h, 16, 64)
		bus.RAM[0x8000+offset] = uint8(value)
		offset++
	}

	bus.CPU.Disassemble(0x0000, 0xFFFF)
	bus.CPU.Reset()

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x80

	fmt.Println("Start PC, RAM[PC]")
	fmt.Println(bus.CPU.PC(), bus.RAM[bus.CPU.PC()])
	fmt.Println()

	for n := 0; n < 120; n++ {
		for {
			bus.CPU.Clock()
			if bus.CPU.GetCycles() == 0 {
				break
			}
			time.Sleep(time.Millisecond * 100)

		}
	}
}
*/
