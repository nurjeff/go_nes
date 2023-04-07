package mappers

type Mapper2 struct {
	Mapper
	PRGBanks   uint8
	CHRBanks   uint8
	MappedAddr uint32

	PRGBankSelectLO uint8
	PRGBankSelectHI uint8
}

func (m *Mapper2) CPUMapRead(addr uint16, data *uint8) (bool, uint32) {
	var mappedAddr uint32 = uint32(addr)

	if addr >= 0x8000 && addr <= 0xBFFF {
		mappedAddr = uint32(m.PRGBankSelectLO)*0x4000 + (uint32(addr) & 0x3FFF)
		return true, mappedAddr
	}

	if addr >= 0xC000 && addr <= 0xFFFF {
		mappedAddr = uint32(m.PRGBankSelectHI)*0x4000 + (uint32(addr) & 0x3FFF)
		return true, mappedAddr
	}

	return false, mappedAddr
}

func (m *Mapper2) CPUMapWrite(addr uint16, data *uint8) (bool, uint32) {
	mappedAddr := uint32(addr)

	if addr >= 0x8000 && addr <= 0xFFFF {
		m.PRGBankSelectLO = *data & 0x0F
	}

	return false, mappedAddr
}

func (m *Mapper2) PPUMapRead(addr uint16) (bool, uint32) {
	mappedAddr := uint32(addr)
	if addr < 0x2000 {
		mappedAddr = uint32(addr)
		return true, mappedAddr
	}
	return false, mappedAddr
}

func (m *Mapper2) PPUMapWrite(addr uint16) (bool, uint32) {
	mappedAddr := uint32(addr)
	if addr < 0x2000 {
		mappedAddr = uint32(addr)
		return true, mappedAddr
	}
	return false, mappedAddr
}

func (m *Mapper2) Initialize() {

}

func (m *Mapper2) Reset() {
	m.PRGBankSelectLO = 0
	m.PRGBankSelectHI = m.PRGBanks - 1
}

func (m *Mapper2) Mirror() uint8 {
	return VERTICAL
}
