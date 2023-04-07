package mappers

type Mapper0 struct {
	Mapper
	PRGBanks   uint8
	CHRBanks   uint8
	MappedAddr uint32
}

func (m *Mapper0) CPUMapRead(addr uint16, data *uint8) (bool, uint32) {
	var mappedAddr uint32 = uint32(addr)

	if addr >= 0x8000 && addr <= 0xFFFF {
		var tmp uint16 = 0x3FFF
		if m.PRGBanks > 1 {
			tmp = 0x7FFF
		}
		mappedAddr = uint32(addr & tmp)
		return true, mappedAddr
	}

	return false, mappedAddr
}

func (m *Mapper0) CPUMapWrite(addr uint16, data *uint8) (bool, uint32) {
	mappedAddr := uint32(addr)
	if addr >= 0x8000 && addr <= 0xFFFF {
		var tmp uint16 = 0x3FFF
		if m.PRGBanks > 1 {
			tmp = 0x7FFF
		}
		mappedAddr = uint32(addr & tmp)
		return true, mappedAddr
	}
	return false, mappedAddr
}

func (m *Mapper0) PPUMapRead(addr uint16) (bool, uint32) {
	mappedAddr := uint32(addr)
	if addr <= 0x1FFF {
		return true, mappedAddr
	}
	return false, mappedAddr
}

func (m *Mapper0) PPUMapWrite(addr uint16) (bool, uint32) {
	mappedAddr := uint32(addr)
	if addr <= 0x1FFF {
		if m.CHRBanks == 0 {
			return true, mappedAddr
		}
	}
	return false, mappedAddr
}

func (m *Mapper0) Initialize() {

}

func (m *Mapper0) Reset() {

}

func (m *Mapper0) Mirror() uint8 {
	return VERTICAL
}
