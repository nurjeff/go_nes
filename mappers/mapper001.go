package mappers

type Mapper1 struct {
	Mapper
	PRGBanks   uint8
	CHRBanks   uint8
	MappedAddr uint32

	CHRBankSelect4LO uint8
	CHRBankSelect4HI uint8
	CHRBankSelect8   uint8

	PRGBankSelect16LO uint8
	PRGBankSelect16HI uint8
	PRGBankSelect32   uint8

	LoadRegister      uint8
	LoadRegisterCount uint8
	ControlRegister   uint8

	MirrorMode uint8

	RAMStatic [32768]uint8
}

func (m *Mapper1) Initialize() {
	for index := range m.RAMStatic {
		m.RAMStatic[index] = 0x00
	}
}

func (m *Mapper1) Reset() {
	m.ControlRegister = 0x1C
	m.LoadRegister = 0x00
	m.LoadRegisterCount = 0x00
	m.CHRBankSelect4LO = 0x00
	m.CHRBankSelect4HI = 0x00
	m.CHRBankSelect8 = 0x00
	m.PRGBankSelect32 = 0x00
	m.PRGBankSelect16LO = 0x00
	m.PRGBankSelect16HI = m.PRGBanks - 1
}

func (m *Mapper1) Mirror() uint8 {
	return m.MirrorMode
}

func (m *Mapper1) CPUMapRead(addr uint16, data *uint8) (bool, uint32) {
	m.MappedAddr = uint32(addr)
	if addr >= 0x6000 && addr <= 0x7FFF {
		m.MappedAddr = 0xFFFFFFFF
		*data = m.RAMStatic[addr&0x1FFF]
		return true, m.MappedAddr
	}

	if addr >= 0x8000 {
		if (m.ControlRegister & 0b01000) >= 1 {
			if addr >= 0x8000 && addr <= 0xBFFF {
				m.MappedAddr = uint32(m.PRGBankSelect16LO)*0x4000 + (uint32(addr) & 0x3FFF)
				return true, m.MappedAddr
			}
			if addr >= 0xC000 && addr <= 0xFFFF {
				m.MappedAddr = uint32(m.PRGBankSelect16HI)*0x4000 + (uint32(addr) & 0x3FFF)
				return true, m.MappedAddr
			}
		}
	}
	return false, uint32(addr)
}

func (m *Mapper1) CPUMapWrite(addr uint16, data *uint8) (bool, uint32) {
	if addr >= 0x6000 && addr <= 0x7FFF {
		m.MappedAddr = 0xFFFFFFFF
		m.RAMStatic[addr&0x1FFF] = *data
		return true, m.MappedAddr
	}

	if addr >= 0x8000 {
		if (*data & 0x80) >= 1 {
			m.LoadRegister = 0x00
			m.LoadRegisterCount = 0
			m.ControlRegister = m.ControlRegister | 0x0C
		} else {
			m.LoadRegister >>= 1
			m.LoadRegister |= (*data & 0x01) << 4
			m.LoadRegisterCount++

			if m.LoadRegisterCount == 5 {
				var targetRegister uint8 = uint8((addr >> 13) & 0x03)
				if targetRegister == 0 {
					m.ControlRegister = m.LoadRegister & 0x1F
					switch m.ControlRegister & 0x03 {
					case 0:
						m.MirrorMode = ONESCREEN_LO
					case 1:
						m.MirrorMode = ONESCREEN_HI
					case 2:
						m.MirrorMode = VERTICAL
					case 3:
						m.MirrorMode = HORIZONTAL
					}
				} else if targetRegister == 1 {
					if (m.ControlRegister & 0b10000) >= 1 {
						m.CHRBankSelect4LO = m.LoadRegister & 0x1F
					} else {
						m.CHRBankSelect8 = m.LoadRegister & 0x1E
					}
				} else if targetRegister == 2 {
					if (m.ControlRegister & 0b10000) >= 1 {
						m.CHRBankSelect4HI = m.LoadRegister & 0x1F
					}
				} else if targetRegister == 3 {
					var PRGMode uint8 = (m.ControlRegister >> 2) & 0x03
					if PRGMode == 0 || PRGMode == 1 {
						m.PRGBankSelect32 = (m.LoadRegister & 0x0E) >> 1
					} else if PRGMode == 2 {
						m.PRGBankSelect16LO = 0
						m.PRGBankSelect16HI = m.LoadRegister & 0x0F
					} else if PRGMode == 3 {
						m.PRGBankSelect16LO = m.LoadRegister & 0x0F
						m.PRGBankSelect16HI = m.PRGBanks - 1
					}
				}
				m.LoadRegister = 0x00
				m.LoadRegisterCount = 0
			}
		}
	}

	return false, uint32(addr)
}

func (m *Mapper1) PPUMapRead(addr uint16) (bool, uint32) {
	if addr < 0x2000 {
		if m.CHRBanks == 0 {
			m.MappedAddr = uint32(addr)
			return true, m.MappedAddr
		} else {
			if (m.ControlRegister & 0b10000) >= 1 {
				if addr <= 0x0FFF {
					m.MappedAddr = uint32(m.CHRBankSelect4LO)*0x1000 + (uint32(addr) & 0x0FFF)
					return true, m.MappedAddr
				}

				if addr >= 0x1000 && addr <= 0x1FFF {
					m.MappedAddr = uint32(m.CHRBankSelect4HI)*0x1000 + (uint32(addr) & 0x0FFF)
					return true, m.MappedAddr
				}
			} else {
				m.MappedAddr = uint32(m.CHRBankSelect8)*0x1000 + (uint32(addr) & 0x1FFF)
				return true, m.MappedAddr
			}
		}
	}
	return false, uint32(addr)
}

func (m *Mapper1) PPUMapWrite(addr uint16) (bool, uint32) {
	if addr < 0x2000 {
		if m.CHRBanks == 0 {
			m.MappedAddr = uint32(addr)
			return true, m.MappedAddr
		}
		return true, m.MappedAddr
	} else {
		return false, uint32(addr)
	}
}
