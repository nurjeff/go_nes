package mappers

type Mapper interface {
	CPUMapRead(uint16, *uint8) (bool, uint32)
	CPUMapWrite(uint16, *uint8) (bool, uint32)
	PPUMapRead(uint16) (bool, uint32)
	PPUMapWrite(uint16) (bool, uint32)
	Initialize()
	Reset()
	Mirror() uint8
}

const (
	HORIZONTAL   = 1
	VERTICAL     = 2
	ONESCREEN_LO = 3
	ONESCREEN_HI = 4
)
