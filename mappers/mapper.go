package mappers

type Mapper interface {
	CPUMapRead(uint16) (bool, uint32)
	CPUMapWrite(uint16) (bool, uint32)
	PPUMapRead(uint16) (bool, uint32)
	PPUMapWrite(uint16) (bool, uint32)
}
