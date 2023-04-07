package cartridge

import (
	"fmt"
	"os"

	"github.com/sc-js/go_nes/mappers"
)

type Cartridge struct {
	VPRGMemory []uint8
	VCHRMemory []uint8

	MapperID uint8
	PRGBanks uint8
	CHRBanks uint8

	Mirror uint8

	Header cartridgeHeader
	Mapper mappers.Mapper
}

type cartridgeHeader struct {
	Name         string
	PRGRomChunks uint8
	CHRRomChunks uint8
	Mapper1      uint8
	Mapper2      uint8
	PRGRAMSize   uint8
	TVSystem1    uint8
	TVSystem2    uint8
	Unused       string
}

const (
	HORIZONTAL   = 1
	VERTICAL     = 2
	ONESCREEN_LO = 3
	ONESCREEN_HI = 4
)

const (
	HEADER_LENGTH   = 16
	TRAINING_LENGTH = 512
	PRG_BANK_SIZE   = 16384
	CHR_BANK_SIZE   = 8192
)

func (c *Cartridge) readCartridgeData(data []byte) {
	c.Header = cartridgeHeader{
		Name:         string(data[:3]),
		PRGRomChunks: data[4],
		CHRRomChunks: data[5],
		Mapper1:      data[6],
		Mapper2:      data[7],
		PRGRAMSize:   data[8],
		TVSystem1:    data[9],
		TVSystem2:    data[10],
		Unused:       string(data[10:HEADER_LENGTH]),
	}
	c.MapperID = ((c.Header.Mapper2 >> 4) << 4) | (c.Header.Mapper1 >> 4)
	if (c.Header.Mapper1 & 0x01) >= 1 {
		c.Mirror = VERTICAL
	} else {
		c.Mirror = HORIZONTAL
	}

	trLenth := 0

	// Check if we need to skip trainer data
	if (c.Header.Mapper1 & 0x04) >= 1 {
		trLenth = TRAINING_LENGTH

	}

	var fileType uint8 = 1
	if (c.Header.Mapper2 & 0x0C) == 0x08 {
		fileType = 2
	}

	if fileType == 0 {

	} else if fileType == 1 {
		c.PRGBanks = c.Header.PRGRomChunks
		c.VPRGMemory = make([]uint8, PRG_BANK_SIZE*int(c.PRGBanks))
		for ind, i := range data[HEADER_LENGTH+trLenth : HEADER_LENGTH+trLenth+len(c.VPRGMemory)] {
			c.VPRGMemory[ind] = i
		}

		c.CHRBanks = c.Header.CHRRomChunks
		if c.CHRBanks > 0 {
			c.VCHRMemory = make([]uint8, CHR_BANK_SIZE*int(c.CHRBanks))
		} else {
			c.VCHRMemory = make([]uint8, CHR_BANK_SIZE)
		}

		for ind, i := range data[HEADER_LENGTH+trLenth+len(c.VPRGMemory):] { //HEADER_LENGTH+trLenth+len(c.VPRGMemory)+len(c.VCHRMemory)] {
			c.VCHRMemory[ind] = i
		}
	} else if fileType == 2 {
	} else {
		panic("unsupported file type: " + fmt.Sprint(fileType))
	}

	switch c.MapperID {
	case 0:
		c.Mapper = mappers.Mapper0{PRGBanks: c.PRGBanks, CHRBanks: c.CHRBanks}
	case 1:
		c.Mapper = &mappers.Mapper1{PRGBanks: c.PRGBanks, CHRBanks: c.CHRBanks}

	default:
		panic("unimplemented mapper:" + fmt.Sprint(c.MapperID))
	}

	c.Mapper.Initialize()
	c.Mapper.Reset()
	fmt.Println("Using Mapper:", c.MapperID)
}

func (c *Cartridge) Initialize(filepath string) {
	f, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fileinfo, err := f.Stat()
	if err != nil {
		panic(err)
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = f.Read(buffer)
	if err != nil {
		panic(err)
	}

	c.readCartridgeData(buffer)
}

func (p *Cartridge) CPUWrite(addr uint16, data uint8) bool {
	cng, add := p.Mapper.CPUMapWrite(addr, &data)
	if cng {
		if add == 0xFFFFFFFF {
			return true
		} else {
			p.VPRGMemory[add] = data
		}
		return true
	}
	return false
}

func (p *Cartridge) CPURead(addr uint16, data *uint8) bool {
	cng, add := p.Mapper.CPUMapRead(addr, data)
	if cng {
		if add == 0xFFFFFFFF {
			return true
		} else {
			*data = p.VPRGMemory[add]
		}

		return true
	}
	return false
}

func (p *Cartridge) PPURead(addr uint16, data *uint8) bool {
	cng, add := p.Mapper.PPUMapRead(addr)
	if cng {
		*data = p.VCHRMemory[add]
		return true
	}
	return false
}

func (p *Cartridge) PPUWrite(addr uint16, data uint8) bool {
	cng, add := p.Mapper.PPUMapWrite(addr)
	if cng {
		p.VCHRMemory[add] = data
		return true
	}
	return false
}

func (p *Cartridge) GetMirror() uint8 {
	m := p.Mapper.Mirror()
	if m == 0 {
		return p.Mirror
	} else {
		return m
	}
}
