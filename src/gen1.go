package main

import (
	"math"
	"sort"
)

func InitGen1() {
	InitGen1IdMap()
	InitGen1HeldItem()
}

type PK1G struct {
	Generation        uint8
	Species           uint8
	CurrentHp         uint16
	Lvl               uint8
	Status            uint8
	Type1             uint8
	Type2             uint8
	HeldItem          uint8
	Move1             uint8
	Move2             uint8
	Move3             uint8
	Move4             uint8
	OriginalTrainerID uint16
	Experience        uint32
	HpEV              uint16
	AttackEV          uint16
	DefenseEV         uint16
	SpeedEV           uint16
	SpecialEV         uint16
	IVData            uint16
	Move1PP           uint8
	Move2PP           uint8
	Move3PP           uint8
	Move4PP           uint8
	HpIV              uint8
	AttackIV          uint8
	DefenseIV         uint8
	SpeedIV           uint8
	SpecialIV         uint8
	Hp                uint8
	Attack            uint8
	Defense           uint8
	Speed             uint8
	Special           uint8
}

const (
	PK1GStatusAsleep    = 0x04
	PK1GStatusPoisoned  = 0x08
	PK1GStatusBurned    = 0x10
	PK1GStatusFrozen    = 0x20
	PK1GStatusParalyzed = 0x40
)

func newPK1G(data []byte) *PK1G {
	pk := PK1G{}
	pk.Generation = 1

	pk.Species = data[0]
	pk.CurrentHp = UInt16(data, 1)
	pk.Lvl = data[3]
	pk.Status = data[4]
	pk.Type1 = data[5]
	pk.Type2 = data[6]
	pk.HeldItem = data[7]
	pk.Move1 = data[8]
	pk.Move2 = data[9]
	pk.Move3 = data[10]
	pk.Move4 = data[11]
	pk.OriginalTrainerID = UInt16(data, 12)
	pk.Experience = UInt24(data, 14) //TODO repair this
	pk.HpEV = UInt16(data, 17)
	pk.AttackEV = UInt16(data, 19)
	pk.DefenseEV = UInt16(data, 21)
	pk.SpeedEV = UInt16(data, 23)
	pk.SpecialEV = UInt16(data, 25)
	pk.IVData = UInt16(data, 27)
	pk.Move1PP = data[29]
	pk.Move2PP = data[30]
	pk.Move3PP = data[31]
	pk.Move4PP = data[32]

	pk.CalculateIV()
	pk.CalculateStats()

	return &pk
}

func Gen1CalculateHpStat(base float64, iv float64, ev float64, level float64) float64 {
	return ((((base+iv)*2 + (math.Sqrt(ev) / 4)) * level) / 100) + level + 10
}

func Gen1CalculateOtherStat(base float64, iv float64, ev float64, level float64) float64 {
	return ((((base+iv)*2 + (math.Sqrt(ev) / 4)) * level) / 100) + 5
}

func (pk *PK1G) CalculateIV() {
	pk.SpecialIV = uint8(pk.IVData & 0x000F)
	pk.SpeedIV = uint8(pk.IVData & 0x00F0 >> 4)
	pk.DefenseIV = uint8(pk.IVData & 0x0F00 >> 8)
	pk.AttackIV = uint8(pk.IVData & 0xF000 >> 12)

	pk.HpIV = (pk.AttackIV&0x1)<<3 | (pk.DefenseIV&0x1)>>2 | (pk.SpeedIV&0x1)>>1 | (pk.SpecialIV & 0x1)
}

func (pk *PK1G) CalculateStats() {
	var base Gen1Base
	if !pk.IsMissingNo() {
		base = Gen1BaseStats[pk.Id()-1]
	} else {
		base = Gen1BaseStats[0]
	}

	pk.Hp = uint8(Gen1CalculateHpStat(float64(base.Hp), float64(pk.HpIV), float64(pk.HpEV), float64(pk.Lvl)))
	pk.Attack = uint8(Gen1CalculateOtherStat(float64(base.Attack), float64(pk.AttackIV), float64(pk.AttackEV), float64(pk.Lvl)))
	pk.Defense = uint8(Gen1CalculateOtherStat(float64(base.Defense), float64(pk.DefenseIV), float64(pk.DefenseEV), float64(pk.Lvl)))
	pk.Speed = uint8(Gen1CalculateOtherStat(float64(base.Speed), float64(pk.SpeedIV), float64(pk.SpeedEV), float64(pk.Lvl)))
	pk.Special = uint8(Gen1CalculateOtherStat(float64(base.Special), float64(pk.SpecialIV), float64(pk.SpecialEV), float64(pk.Lvl)))
}

func (pk *PK1G) Gen() uint8 {
	return pk.Generation
}

func (pk *PK1G) Id() uint16 {
	return uint16(Gen1IdMap[pk.Species])
}

func (pk *PK1G) Form() uint8 {
	return 0
}

func (pk *PK1G) Nickname() string {
	return ""
}

func (pk *PK1G) Level() uint8 {
	return pk.Lvl
}

func (pk *PK1G) Held() uint16 {
	return uint16(pk.HeldItem)
}

func (pk *PK1G) IsAsleep() bool {
	return CheckFlag(pk.Status, PK1GStatusAsleep)
}

func (pk *PK1G) IsPoisoned() bool {
	return CheckFlag(pk.Status, PK1GStatusPoisoned)
}

func (pk *PK1G) IsBurned() bool {
	return CheckFlag(pk.Status, PK1GStatusBurned)
}

func (pk *PK1G) IsFrozen() bool {
	return CheckFlag(pk.Status, PK1GStatusFrozen)
}

func (pk *PK1G) IsParalyzed() bool {
	return CheckFlag(pk.Status, PK1GStatusParalyzed)
}

func (pk *PK1G) EV() []uint16 {
	return []uint16{pk.HpEV, pk.AttackEV, pk.DefenseEV, pk.SpeedEV, pk.SpecialEV, pk.SpecialEV}
}

func (pk *PK1G) IV() []uint16 {
	return []uint16{uint16(pk.HpIV), uint16(pk.AttackIV), uint16(pk.DefenseIV), uint16(pk.SpeedIV), uint16(pk.SpecialIV), uint16(pk.SpecialIV)}
}

func (pk *PK1G) Stats() []uint16 {
	return []uint16{uint16(pk.Hp), uint16(pk.Attack), uint16(pk.Defense), uint16(pk.Speed), uint16(pk.Special), uint16(pk.Special)}
}

func (pk *PK1G) Moves() []uint16 {
	return []uint16{uint16(pk.Move1), uint16(pk.Move2),
		uint16(pk.Move3), uint16(pk.Move4)}
}

func (pk *PK1G) PP() []uint8 {
	return []uint8{pk.Move1PP, pk.Move2PP, pk.Move3PP, pk.Move4PP}
}

var (
	Gen1ShinyAttackIV = []uint8{2, 3, 6, 7, 10, 11, 14, 15}
)

func (pk *PK1G) IsShiny() bool {
	i := sort.Search(len(Gen1ShinyAttackIV), func(i int) bool {
		return Gen1ShinyAttackIV[i] >= pk.AttackIV
	})
	if i < len(Gen1ShinyAttackIV) && Gen1ShinyAttackIV[i] == pk.AttackIV {
		return pk.DefenseIV == 10 && pk.SpeedIV == 10 && pk.SpecialIV == 10
	} else {
		return false
	}
}

func (pk *PK1G) IsMissingNo() bool {
	return pk.Id() == 0 || pk.Id() > 151
}

func (pk *PK1G) Bytes() []byte {
	data := make([]byte, 33)

	data[0] = pk.Species
	UInt16ToBytes(data, 1, pk.CurrentHp)
	data[3] = pk.Lvl
	data[4] = pk.Status
	data[5] = pk.Type1
	data[6] = pk.Type2
	data[7] = pk.HeldItem
	data[8] = pk.Move1
	data[9] = pk.Move2
	data[10] = pk.Move3
	data[11] = pk.Move4
	UInt16ToBytes(data, 12, pk.OriginalTrainerID)
	UInt24ToBytes(data, 14, pk.Experience) //TODO repair this
	UInt16ToBytes(data, 17, pk.HpEV)
	UInt16ToBytes(data, 19, pk.AttackEV)
	UInt16ToBytes(data, 21, pk.DefenseEV)
	UInt16ToBytes(data, 23, pk.SpeedEV)
	UInt16ToBytes(data, 25, pk.SpecialEV)
	UInt16ToBytes(data, 27, pk.IVData)
	data[29] = pk.Move1PP
	data[30] = pk.Move2PP
	data[31] = pk.Move3PP
	data[32] = pk.Move4PP

	return data
}

func (pk *PK1G) Upgrade() (PK, error) {
	npk := PK2G{}
	npk.Generation = 2

	npk.Species = Gen1IdMap[pk.Species]
	if val, ok := Gen1HeldItemToGen2[pk.HeldItem]; ok {
		npk.HeldItem = val
	} else {
		npk.HeldItem = pk.HeldItem
	}
	npk.Move1 = pk.Move1
	npk.Move2 = pk.Move2
	npk.Move3 = pk.Move3
	npk.Move4 = pk.Move4
	npk.OriginalTrainerID = pk.OriginalTrainerID
	npk.Experience = pk.Experience
	npk.HpEV = pk.HpEV
	npk.AttackEV = pk.AttackEV
	npk.DefenseEV = pk.DefenseEV
	npk.SpeedEV = pk.SpeedEV
	npk.SpecialEV = pk.SpecialEV
	npk.IVData = pk.IVData
	npk.Move1PP = pk.Move1PP
	npk.Move2PP = pk.Move2PP
	npk.Move3PP = pk.Move3PP
	npk.Move4PP = pk.Move4PP
	npk.Friendship = 0
	npk.Pokerus = 0
	npk.CaughtData = 0
	npk.Lvl = pk.Level()

	npk.CalculateIV()
	npk.CalculateStats()

	return &npk, nil
}

var (
	Gen1IdMap = make([]uint8, 256)
)

func InitGen1IdMap() {
	Gen1IdMap[100] = 39
	Gen1IdMap[101] = 40
	Gen1IdMap[102] = 133
	Gen1IdMap[103] = 136
	Gen1IdMap[104] = 135
	Gen1IdMap[105] = 134
	Gen1IdMap[106] = 66
	Gen1IdMap[107] = 41
	Gen1IdMap[108] = 23
	Gen1IdMap[109] = 46
	Gen1IdMap[110] = 61
	Gen1IdMap[111] = 62
	Gen1IdMap[112] = 13
	Gen1IdMap[113] = 14
	Gen1IdMap[114] = 15
	Gen1IdMap[116] = 85
	Gen1IdMap[117] = 57
	Gen1IdMap[118] = 51
	Gen1IdMap[119] = 49
	Gen1IdMap[120] = 87
	Gen1IdMap[123] = 10
	Gen1IdMap[124] = 11
	Gen1IdMap[125] = 12
	Gen1IdMap[126] = 68
	Gen1IdMap[128] = 55
	Gen1IdMap[129] = 97
	Gen1IdMap[130] = 42
	Gen1IdMap[131] = 150
	Gen1IdMap[132] = 143
	Gen1IdMap[133] = 129
	Gen1IdMap[136] = 89
	Gen1IdMap[138] = 99
	Gen1IdMap[139] = 91
	Gen1IdMap[141] = 101
	Gen1IdMap[142] = 36
	Gen1IdMap[143] = 110
	Gen1IdMap[144] = 53
	Gen1IdMap[145] = 105
	Gen1IdMap[147] = 93
	Gen1IdMap[148] = 63
	Gen1IdMap[149] = 65
	Gen1IdMap[150] = 17
	Gen1IdMap[151] = 18
	Gen1IdMap[152] = 121
	Gen1IdMap[153] = 1
	Gen1IdMap[154] = 3
	Gen1IdMap[155] = 73
	Gen1IdMap[157] = 118
	Gen1IdMap[158] = 119
	Gen1IdMap[163] = 77
	Gen1IdMap[164] = 78
	Gen1IdMap[165] = 19
	Gen1IdMap[166] = 20
	Gen1IdMap[167] = 33
	Gen1IdMap[168] = 30
	Gen1IdMap[169] = 74
	Gen1IdMap[170] = 137
	Gen1IdMap[171] = 142
	Gen1IdMap[173] = 81
	Gen1IdMap[176] = 4
	Gen1IdMap[177] = 7
	Gen1IdMap[178] = 5
	Gen1IdMap[179] = 8
	Gen1IdMap[180] = 6
	Gen1IdMap[185] = 43
	Gen1IdMap[186] = 44
	Gen1IdMap[187] = 45
	Gen1IdMap[188] = 69
	Gen1IdMap[189] = 70
	Gen1IdMap[190] = 71
	Gen1IdMap[1] = 112
	Gen1IdMap[2] = 115
	Gen1IdMap[3] = 32
	Gen1IdMap[4] = 35
	Gen1IdMap[5] = 21
	Gen1IdMap[6] = 100
	Gen1IdMap[7] = 34
	Gen1IdMap[8] = 80
	Gen1IdMap[9] = 2
	Gen1IdMap[10] = 103
	Gen1IdMap[11] = 108
	Gen1IdMap[12] = 102
	Gen1IdMap[13] = 88
	Gen1IdMap[14] = 94
	Gen1IdMap[15] = 29
	Gen1IdMap[16] = 31
	Gen1IdMap[17] = 104
	Gen1IdMap[18] = 111
	Gen1IdMap[19] = 131
	Gen1IdMap[20] = 59
	Gen1IdMap[21] = 151
	Gen1IdMap[22] = 130
	Gen1IdMap[23] = 90
	Gen1IdMap[24] = 72
	Gen1IdMap[25] = 92
	Gen1IdMap[26] = 123
	Gen1IdMap[27] = 120
	Gen1IdMap[28] = 9
	Gen1IdMap[29] = 127
	Gen1IdMap[30] = 114
	Gen1IdMap[33] = 58
	Gen1IdMap[34] = 95
	Gen1IdMap[35] = 22
	Gen1IdMap[36] = 16
	Gen1IdMap[37] = 79
	Gen1IdMap[38] = 64
	Gen1IdMap[39] = 75
	Gen1IdMap[40] = 113
	Gen1IdMap[41] = 67
	Gen1IdMap[42] = 122
	Gen1IdMap[43] = 106
	Gen1IdMap[44] = 107
	Gen1IdMap[45] = 24
	Gen1IdMap[46] = 47
	Gen1IdMap[47] = 54
	Gen1IdMap[48] = 96
	Gen1IdMap[49] = 76
	Gen1IdMap[51] = 126
	Gen1IdMap[53] = 125
	Gen1IdMap[54] = 82
	Gen1IdMap[55] = 109
	Gen1IdMap[57] = 56
	Gen1IdMap[58] = 86
	Gen1IdMap[59] = 50
	Gen1IdMap[60] = 128
	Gen1IdMap[65] = 48
	Gen1IdMap[66] = 149
	Gen1IdMap[70] = 84
	Gen1IdMap[71] = 60
	Gen1IdMap[72] = 124
	Gen1IdMap[73] = 146
	Gen1IdMap[74] = 144
	Gen1IdMap[75] = 145
	Gen1IdMap[76] = 132
	Gen1IdMap[77] = 52
	Gen1IdMap[78] = 98
	Gen1IdMap[82] = 37
	Gen1IdMap[83] = 38
	Gen1IdMap[84] = 25
	Gen1IdMap[85] = 26
	Gen1IdMap[88] = 147
	Gen1IdMap[89] = 148
	Gen1IdMap[90] = 140
	Gen1IdMap[91] = 141
	Gen1IdMap[92] = 116
	Gen1IdMap[93] = 117
	Gen1IdMap[96] = 27
	Gen1IdMap[97] = 28
	Gen1IdMap[98] = 138
	Gen1IdMap[99] = 139
}

var (
	Gen1HeldItemToGen2 = make(map[uint8]uint8)
)

func InitGen1HeldItem() {
	Gen1HeldItemToGen2[0x19] = 0x92
	Gen1HeldItemToGen2[0x2D] = 0x53
	Gen1HeldItemToGen2[0x32] = 0xAE
	Gen1HeldItemToGen2[0x5A] = 0xAD
	Gen1HeldItemToGen2[0x64] = 0xAD
	Gen1HeldItemToGen2[0x78] = 0xAD
	Gen1HeldItemToGen2[0x7F] = 0xAD
	Gen1HeldItemToGen2[0xBE] = 0xAD
	Gen1HeldItemToGen2[0xFF] = 0xAD
}

type Gen1Base struct {
	Hp      int `json:"hp"`
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	Speed   int `json:"speed"`
	Special int `json:"special"`
}

var (
	Gen1BaseStats = []Gen1Base{
		{45, 49, 49, 45, 65},
		{60, 62, 63, 60, 80},
		{80, 82, 83, 80, 100},
		{39, 52, 43, 65, 50},
		{58, 64, 58, 80, 65},
		{78, 84, 78, 100, 85},
		{44, 48, 65, 43, 50},
		{59, 63, 80, 58, 65},
		{79, 83, 100, 78, 85},
		{45, 30, 35, 45, 20},
		{50, 20, 55, 30, 25},
		{60, 45, 50, 70, 80},
		{40, 35, 30, 50, 20},
		{45, 25, 50, 35, 25},
		{65, 80, 40, 75, 45},
		{40, 45, 40, 56, 35},
		{63, 60, 55, 71, 50},
		{83, 80, 75, 91, 70},
		{30, 56, 35, 72, 25},
		{55, 81, 60, 97, 50},
		{40, 60, 30, 70, 31},
		{65, 90, 65, 100, 61},
		{35, 60, 44, 55, 40},
		{60, 85, 69, 80, 65},
		{35, 55, 30, 90, 50},
		{60, 90, 55, 100, 90},
		{50, 75, 85, 40, 30},
		{75, 100, 110, 65, 55},
		{55, 47, 52, 41, 40},
		{70, 62, 67, 56, 55},
		{90, 82, 87, 76, 75},
		{46, 57, 40, 50, 40},
		{61, 72, 57, 65, 55},
		{81, 92, 77, 85, 75},
		{70, 45, 48, 35, 60},
		{95, 70, 73, 60, 85},
		{38, 41, 40, 65, 65},
		{73, 76, 75, 100, 100},
		{115, 45, 20, 20, 25},
		{140, 70, 45, 45, 50},
		{40, 45, 35, 55, 40},
		{75, 80, 70, 90, 75},
		{45, 50, 55, 30, 75},
		{60, 65, 70, 40, 85},
		{75, 80, 85, 50, 100},
		{35, 70, 55, 25, 55},
		{60, 95, 80, 30, 80},
		{60, 55, 50, 45, 40},
		{70, 65, 60, 90, 90},
		{10, 55, 25, 95, 45},
		{35, 80, 50, 120, 70},
		{40, 45, 35, 90, 40},
		{65, 70, 60, 115, 65},
		{50, 52, 48, 55, 50},
		{80, 82, 78, 85, 80},
		{40, 80, 35, 70, 35},
		{65, 105, 60, 95, 60},
		{55, 70, 45, 60, 50},
		{90, 110, 80, 95, 80},
		{40, 50, 40, 90, 40},
		{65, 65, 65, 90, 50},
		{90, 85, 95, 70, 70},
		{25, 20, 15, 90, 105},
		{40, 35, 30, 105, 120},
		{55, 50, 45, 120, 135},
		{70, 80, 50, 35, 35},
		{80, 100, 70, 45, 50},
		{90, 130, 80, 55, 65},
		{50, 75, 35, 40, 70},
		{65, 90, 50, 55, 85},
		{80, 105, 65, 70, 100},
		{40, 40, 35, 70, 100},
		{80, 70, 65, 100, 120},
		{40, 80, 100, 20, 30},
		{55, 95, 115, 35, 45},
		{80, 110, 130, 45, 55},
		{50, 85, 55, 90, 65},
		{65, 100, 70, 105, 80},
		{90, 65, 65, 15, 40},
		{95, 75, 110, 30, 80},
		{25, 35, 70, 45, 95},
		{50, 60, 95, 70, 120},
		{52, 65, 55, 60, 58},
		{35, 85, 45, 75, 35},
		{60, 110, 70, 100, 60},
		{65, 45, 55, 45, 70},
		{90, 70, 80, 70, 95},
		{80, 80, 50, 25, 40},
		{105, 105, 75, 50, 65},
		{30, 65, 100, 40, 45},
		{50, 95, 180, 70, 85},
		{30, 35, 30, 80, 100},
		{45, 50, 45, 95, 115},
		{60, 65, 60, 110, 130},
		{35, 45, 160, 70, 30},
		{60, 48, 45, 42, 90},
		{85, 73, 70, 67, 115},
		{30, 105, 90, 50, 25},
		{55, 130, 115, 75, 50},
		{40, 30, 50, 100, 55},
		{60, 50, 70, 140, 80},
		{60, 40, 80, 40, 60},
		{95, 95, 85, 55, 125},
		{50, 50, 95, 35, 40},
		{60, 80, 110, 45, 50},
		{50, 120, 53, 87, 35},
		{50, 105, 79, 76, 35},
		{90, 55, 75, 30, 60},
		{40, 65, 95, 35, 60},
		{65, 90, 120, 60, 85},
		{80, 85, 95, 25, 30},
		{105, 130, 120, 40, 45},
		{250, 5, 5, 50, 105},
		{65, 55, 115, 60, 100},
		{105, 95, 80, 90, 40},
		{30, 40, 70, 60, 70},
		{55, 65, 95, 85, 95},
		{45, 67, 60, 63, 50},
		{80, 92, 65, 68, 80},
		{30, 45, 55, 85, 70},
		{60, 75, 85, 115, 100},
		{40, 45, 65, 90, 100},
		{70, 110, 80, 105, 55},
		{65, 50, 35, 95, 95},
		{65, 83, 57, 105, 85},
		{65, 95, 57, 93, 85},
		{65, 125, 100, 85, 55},
		{75, 100, 95, 110, 70},
		{20, 10, 55, 80, 20},
		{95, 125, 79, 81, 100},
		{130, 85, 80, 60, 95},
		{48, 48, 48, 48, 48},
		{55, 55, 50, 55, 65},
		{130, 65, 60, 65, 110},
		{65, 65, 60, 130, 110},
		{65, 130, 60, 65, 110},
		{65, 60, 70, 40, 75},
		{35, 40, 100, 35, 90},
		{70, 60, 125, 55, 115},
		{30, 80, 90, 55, 45},
		{60, 115, 105, 80, 70},
		{80, 105, 65, 130, 60},
		{160, 110, 65, 30, 65},
		{90, 85, 100, 85, 125},
		{90, 90, 85, 100, 125},
		{90, 100, 90, 90, 125},
		{41, 64, 45, 50, 50},
		{61, 84, 65, 70, 70},
		{91, 134, 95, 80, 100},
		{106, 110, 90, 130, 154},
		{100, 100, 100, 100, 100},
	}
)
