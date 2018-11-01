package main

type SavVersion uint8

const (
	SavInvalid SavVersion = iota
	SavGen1
)

const (
	/*SizeG7USUM     = 0x6CC00
	SizeG7SM       = 0x6BE00
	SizeG6XY       = 0x65600
	SizeG6ORAS     = 0x76000
	SizeG6ORASDEMO = 0x5A00
	SizeG5RAW      = 0x80000
	SizeG5BW       = 0x24000
	SizeG5B2W2     = 0x26000
	SizeG4BR       = 0x380000
	SizeG4RAW      = 0x80000
	SizeG3BOX      = 0x76000
	SizeG3BOXGCI   = 0x76040
	SizeG3COLO     = 0x60000
	SizeG3COLOGCI  = 0x60040
	SizeG3XD       = 0x56000
	SizeG3XDGCI    = 0x56040
	SizeG3RAW      = 0x20000
	SizeG3RAWHALF  = 0x10000
	SizeG2RAW_U    = 0x8000
	SizeG2VC_U     = 0x8010
	SizeG2BAT_U    = 0x802C
	SizeG2EMU_U    = 0x8030
	SizeG2RAW_J    = 0x10000
	SizeG2VC_J     = 0x10010
	SizeG2BAT_J    = 0x1002C
	SizeG2EMU_J    = 0x10030*/
	SizeG1RAW = 0x8000
	SizeG1BAT = 0x802C
)

func IsValidSavSize(len int64) bool {
	return len == SizeG1BAT ||
		len == SizeG1RAW
}

func DetectSavVersion(data []byte) SavVersion {
	if IsSavGen1(data) {
		return SavGen1
	}
	return SavInvalid
}

func IsG12ListValid(data []byte, offset int, listCount int) bool {
	numEntries := int(data[offset])
	return numEntries <= listCount && data[offset+1+numEntries] == 0xFF
}

func IsSavGen1U(data []byte) bool {
	return IsG12ListValid(data, 0x2F2C, 20) && IsG12ListValid(data, 0x30C0, 20)
}

func IsSavGen1J(data []byte) bool {
	return IsG12ListValid(data, 0x2ED5, 30) && IsG12ListValid(data, 0x302D, 30)
}

func IsSavGen1(data []byte) bool {
	if len(data) != SizeG1RAW && len(data) != SizeG1BAT {
		return false
	}
	if !IsSavGen1J(data) && !IsSavGen1U(data) {
		return false
	}
	return true
}

func CreateComputerGen1(data []byte) *Computer {
	cd := CreateComputer(uint8(len(BankData.Computers)), 1)
	BankData.Computers = append(BankData.Computers, *cd)
	c := &BankData.Computers[len(BankData.Computers)-1]

	//TODO read save data
	return c
}
