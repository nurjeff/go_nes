package main

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type SDLController struct {
	Window  *sdl.Window
	Surface *sdl.Surface
	ResX    uint
	ResY    uint
	Fonts   map[int]*ttf.Font
	Running bool
	Bus     *Bus
	Rand    int
}

var fonts []int = []int{FONT_12, FONT_13, FONT_14, FONT_15, FONT_18, FONT_20, FONT_24, FONT_32, FONT_38, FONT_42}

const (
	FONT_12 = 12
	FONT_13 = 13
	FONT_14 = 14
	FONT_15 = 15
	FONT_18 = 18
	FONT_20 = 20
	FONT_24 = 24
	FONT_32 = 32
	FONT_38 = 38
	FONT_42 = 42

	WHITE  = 80
	RED    = 81
	YELLOW = 82
	GREEN  = 83
)

func getColor(font int) sdl.Color {
	switch font {
	case WHITE:
		return sdl.Color{R: 255, G: 255, B: 255, A: 255}
	case RED:
		return sdl.Color{R: 255, G: 0, B: 0, A: 255}
	case YELLOW:
		return sdl.Color{R: 245, G: 215, B: 40, A: 255}
	case GREEN:
		return sdl.Color{R: 0, G: 255, B: 0, A: 255}
	}
	return sdl.Color{R: 255, G: 255, B: 255, A: 255}
}

func (c *SDLController) Initialize(resx uint, resy uint, fontname string) error {
	c.Rand = rand.Int()
	c.Fonts = make(map[int]*ttf.Font)
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return err
	}
	if resx == 0 || resy == 0 {
		return errors.New("resolution axis can't be 0")
	}

	c.ResX = resx
	c.ResY = resy

	if window, err := sdl.CreateWindow("GO NES", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, int32(resx), int32(resy), sdl.WINDOW_SHOWN|sdl.WINDOW_ALLOW_HIGHDPI); err != nil {
		return err
	} else {
		c.Window = window
		defer c.Window.Destroy()
	}

	if surface, err := c.Window.GetSurface(); err != nil {
		return err
	} else {
		c.Surface = surface
	}

	if err := ttf.Init(); err != nil {
		return err
	} else {
		defer ttf.Quit()
	}

	for _, element := range fonts {
		if font, err := ttf.OpenFont(fmt.Sprintf("./assets/%s.ttf", fontname), element); err != nil {
			return err
		} else {
			c.Fonts[element] = font
			defer c.Fonts[element].Close()
		}
	}

	c.Refresh()
	c.Start()

	return nil
}

func (c *SDLController) Refresh() {
	c.Surface.Free()
	rect := sdl.Rect{X: 0, Y: 0, W: c.Surface.W, H: c.Surface.H}
	c.Surface.FillRect(&rect, 0x1e2124)

	c.DrawCPUFlags()
	c.DrawRAMPage0()
	c.DrawRAMPage8000()
	c.DrawDisassembly()
	c.Window.UpdateSurface()
}

func (c *SDLController) Start() {
	c.Running = true

	for c.Running {
		c.Refresh()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			//rect := sdl.Rect{X: 0, Y: 0, W: c.Surface.W, H: c.Surface.H}
			//c.Surface.FillRect(&rect, 0x0)
			case *sdl.QuitEvent:
				c.Running = false
			case *sdl.KeyboardEvent:
				if t.State == sdl.RELEASED {
					if t.Keysym.Sym == sdl.K_SPACE {
						for {
							c.Bus.CPU.Clock()
							if c.Bus.CPU.cycles == 0 {
								break
							}
						}
					}
				}
			}
		}

		sdl.Delay(16)
	}
}

func (c *SDLController) DrawRAMPage0() {
	var offset uint = 34
	c.DrawText(30, 17, "RAM - 0x0000 : 0x0100", FONT_20, WHITE)
	c.DrawText(30, offset, "------", FONT_20, WHITE)
	brk := 0
	for index, _ := range c.Bus.RAM[:256] {
		col := WHITE
		if c.Bus.RAM[index] > 0 {
			col = YELLOW
		}
		c.DrawTextCentered(34+27*uint(index-brk*27), offset+17+uint(brk)*24, fmt.Sprint(c.Bus.RAM[index]), FONT_15, col)
		if ((index + 1) % 27) == 0 {
			brk++
		}
	}
}

func (c *SDLController) DrawRAMPage8000() {
	var offset uint = 400
	c.DrawText(30, offset-17, "RAM - 0x8000 : 0x8100", FONT_20, WHITE)
	c.DrawText(30, offset, "------", FONT_20, WHITE)
	brk := 0
	for index, _ := range c.Bus.RAM[32768:33024] {
		col := WHITE
		if c.Bus.RAM[index+32768] > 0 {
			col = YELLOW
		}
		c.DrawTextCentered(34+27*uint(index-brk*27), offset+17+uint(brk)*24, fmt.Sprint(c.Bus.RAM[index+32768]), FONT_15, col)
		if ((index + 1) % 27) == 0 {
			brk++
		}
	}
}

func (c *SDLController) DrawCPUFlags() {
	var offset uint = 34
	var xoff uint = 280
	var flagoff uint = 180

	c.DrawText(uint(c.Surface.W)-xoff, 17, "CPU", FONT_20, WHITE)
	c.DrawText(uint(c.Surface.W)-xoff, offset, "------", FONT_20, WHITE)
	c.DrawText(uint(c.Surface.W)-xoff, 17+offset, fmt.Sprint("A:    ", hex(c.Bus.CPU.a, 2)), FONT_20, WHITE)
	c.DrawText(uint(c.Surface.W)-xoff, 34+offset, fmt.Sprint("X:    ", hex(c.Bus.CPU.x, 2)), FONT_20, WHITE)
	c.DrawText(uint(c.Surface.W)-xoff, 51+offset, fmt.Sprint("Y:    ", hex(c.Bus.CPU.y, 2)), FONT_20, WHITE)
	c.DrawText(uint(c.Surface.W)-xoff, 68+offset, fmt.Sprint("S:    ", hex(c.Bus.CPU.stkp, 2)), FONT_20, WHITE)
	c.DrawText(uint(c.Surface.W)-xoff, 85+offset, fmt.Sprint("P:    ", hex(c.Bus.CPU.status, 2)), FONT_20, WHITE)
	c.DrawText(uint(c.Surface.W)-xoff, 102+offset, fmt.Sprintf("PC:  0x%x", c.Bus.CPU.pc), FONT_20, WHITE)
	c.DrawText(uint(c.Surface.W)-xoff, 119+offset, fmt.Sprint("TC:  ", c.Bus.CPU.totalCycles), FONT_20, WHITE)

	if c.Bus.CPU.GetFlag(c.Bus.CPU.flags.B) == 0 {
		c.DrawText((uint(c.Surface.W)-xoff)+flagoff, 17, "B", FONT_20, RED)
	} else {
		c.DrawText((uint(c.Surface.W)-xoff)+flagoff, 17, "B", FONT_20, GREEN)
	}

	if c.Bus.CPU.GetFlag(c.Bus.CPU.flags.C) == 0 {
		c.DrawText((uint(c.Surface.W)-xoff)+flagoff, 34, "C", FONT_20, RED)
	} else {
		c.DrawText((uint(c.Surface.W)-xoff)+flagoff, 34, "C", FONT_20, GREEN)
	}

	if c.Bus.CPU.GetFlag(c.Bus.CPU.flags.D) == 0 {
		c.DrawText((uint(c.Surface.W)-xoff)+flagoff, 51, "D", FONT_20, RED)
	} else {
		c.DrawText((uint(c.Surface.W)-xoff)+flagoff, 51, "D", FONT_20, GREEN)
	}

	if c.Bus.CPU.GetFlag(c.Bus.CPU.flags.I) == 0 {
		c.DrawText((uint(c.Surface.W)-xoff)+flagoff, 68, "I", FONT_20, RED)
	} else {
		c.DrawText((uint(c.Surface.W)-xoff)+flagoff, 68, "I", FONT_20, GREEN)
	}

	if c.Bus.CPU.GetFlag(c.Bus.CPU.flags.N) == 0 {
		c.DrawText((uint(c.Surface.W)-xoff)+flagoff, 85, "N", FONT_20, RED)
	} else {
		c.DrawText((uint(c.Surface.W)-xoff)+flagoff, 85, "N", FONT_20, GREEN)
	}

	if c.Bus.CPU.GetFlag(c.Bus.CPU.flags.U) == 0 {
		c.DrawText((uint(c.Surface.W)-xoff)+flagoff, 102, "U", FONT_20, RED)
	} else {
		c.DrawText((uint(c.Surface.W)-xoff)+flagoff, 102, "U", FONT_20, GREEN)
	}

	if c.Bus.CPU.GetFlag(c.Bus.CPU.flags.V) == 0 {
		c.DrawText((uint(c.Surface.W)-xoff)+flagoff, 119, "V", FONT_20, RED)
	} else {
		c.DrawText((uint(c.Surface.W)-xoff)+flagoff, 119, "V", FONT_20, GREEN)
	}

	if c.Bus.CPU.GetFlag(c.Bus.CPU.flags.Z) == 0 {
		c.DrawText((uint(c.Surface.W)-xoff)+flagoff, 136, "Z", FONT_20, RED)
	} else {
		c.DrawText((uint(c.Surface.W)-xoff)+flagoff, 136, "Z", FONT_20, GREEN)
	}
}

func (c *SDLController) DrawDisassembly() {
	var offset uint = 240
	var xoff uint = 280
	for i := -10; i < 13; i++ {
		nextInst := ""
		o := 0
		for len(nextInst) == 0 {
			nextInst = c.Bus.CPU.Disassembly[c.Bus.CPU.pc+uint16(i+o)]
			o++
		}
		co := WHITE
		if uint16(i+0)+c.Bus.CPU.pc == c.Bus.CPU.pc {
			co = GREEN
		}
		c.DrawText(uint(c.Surface.W-int32(xoff)), (uint(c.Surface.H)-350)+(uint(i)*20), c.Bus.CPU.Disassembly[c.Bus.CPU.pc+uint16(i)], FONT_15, co)
	}

	_ = offset

}

func (c *SDLController) DrawText(x uint, y uint, text string, size int, color int) error {
	if t, err := c.Fonts[size].RenderUTF8Blended(text, getColor(color)); err != nil {
		return err
	} else {
		defer t.Free()
		if err := t.Blit(nil, c.Surface, &sdl.Rect{X: int32(x), Y: int32(y), W: 0, H: 0}); err != nil {
			return err
		}
	}

	return nil
}

func (c *SDLController) DrawTextCentered(x uint, y uint, text string, size int, color int) error {
	if t, err := c.Fonts[size].RenderUTF8Blended(text, getColor(color)); err != nil {
		fmt.Println(err)
		return err
	} else {
		defer t.Free()
		if err := t.Blit(nil, c.Surface, &sdl.Rect{X: int32(x) - t.W/2, Y: int32(y), W: 0, H: 0}); err != nil {
			return err
		}
	}

	return nil
}
