package main

import (
	"errors"
	"sort"
)

func InitGen2() {
	InitGen1Types()
}

type PK2G struct {
	Generation        uint8
	Species           uint8
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
	Friendship        uint8
	Pokerus           uint8
	CaughtData        uint16
	Lvl               uint8
	HpIV              uint8
	AttackIV          uint8
	DefenseIV         uint8
	SpeedIV           uint8
	SpecialIV         uint8
	Hp                uint8
	Attack            uint8
	Defense           uint8
	Speed             uint8
	SpecialAttack     uint8
	SpecialDefense    uint8
}

func newPK2G(data []byte) *PK2G {
	pk := PK2G{}
	pk.Generation = 2

	pk.Species = data[0]
	pk.HeldItem = data[1]
	pk.Move1 = data[2]
	pk.Move2 = data[3]
	pk.Move3 = data[4]
	pk.Move4 = data[5]
	pk.OriginalTrainerID = UInt16(data, 6)
	pk.Experience = UInt24(data, 8) //TODO repair this
	pk.HpEV = UInt16(data, 11)
	pk.AttackEV = UInt16(data, 13)
	pk.DefenseEV = UInt16(data, 15)
	pk.SpeedEV = UInt16(data, 17)
	pk.SpecialEV = UInt16(data, 19)
	pk.IVData = UInt16(data, 21)
	pk.Move1PP = data[23]
	pk.Move2PP = data[24]
	pk.Move3PP = data[25]
	pk.Move4PP = data[26]
	pk.Friendship = data[27]
	pk.Pokerus = data[28]
	pk.CaughtData = UInt16(data, 29)
	pk.Lvl = data[31]

	pk.CalculateIV()
	pk.CalculateStats()

	return &pk
}

func (pk *PK2G) CalculateIV() {
	pk.SpecialIV = uint8(pk.IVData & 0x000F)
	pk.SpeedIV = uint8(pk.IVData & 0x00F0 >> 4)
	pk.DefenseIV = uint8(pk.IVData & 0x0F00 >> 8)
	pk.AttackIV = uint8(pk.IVData & 0xF000 >> 12)

	pk.HpIV = (pk.AttackIV&0x1)<<3 | (pk.DefenseIV&0x1)>>2 | (pk.SpeedIV&0x1)>>1 | (pk.SpecialIV & 0x1)
}

func (pk *PK2G) CalculateStats() {
	var base Gen2Base
	if !pk.IsMissingNo() {
		base = Gen2BaseStats[pk.Species-1][0]
	} else {
		base = Gen2BaseStats[0][0]
	}

	pk.Hp = uint8(Gen1CalculateHpStat(float64(base.Hp), float64(pk.HpIV), float64(pk.HpEV), float64(pk.Lvl)))
	pk.Attack = uint8(Gen1CalculateOtherStat(float64(base.Attack), float64(pk.AttackIV), float64(pk.AttackEV), float64(pk.Lvl)))
	pk.Defense = uint8(Gen1CalculateOtherStat(float64(base.Defense), float64(pk.DefenseIV), float64(pk.DefenseEV), float64(pk.Lvl)))
	pk.Speed = uint8(Gen1CalculateOtherStat(float64(base.Speed), float64(pk.SpeedIV), float64(pk.SpeedEV), float64(pk.Lvl)))
	pk.SpecialAttack = uint8(Gen1CalculateOtherStat(float64(base.SpecialAttack), float64(pk.SpecialIV), float64(pk.SpecialEV), float64(pk.Lvl)))
	pk.SpecialDefense = uint8(Gen1CalculateOtherStat(float64(base.SpecialDefense), float64(pk.SpecialIV), float64(pk.SpecialEV), float64(pk.Lvl)))
}

func (pk *PK2G) Gen() uint8 {
	return pk.Generation
}

func (pk *PK2G) Id() uint16 {
	return uint16(pk.Species)
}

func (pk *PK2G) Form() uint8 {
	if pk.Species == 201 {
		d := ((pk.IVData & 0x6) >> 1) | ((pk.IVData & 0x60) >> 3) | ((pk.IVData & 0x600) >> 5) | ((pk.IVData & 0x6000) >> 7)
		return uint8(float32(d) / 10)
	}
	return 0
}

func (pk *PK2G) Nickname() string {
	return ""
}

func (pk *PK2G) Level() uint8 {
	return pk.Lvl
}

func (pk *PK2G) Held() uint16 {
	return uint16(pk.HeldItem)
}

func (pk *PK2G) IsAsleep() bool {
	return false
}

func (pk *PK2G) IsPoisoned() bool {
	return false
}

func (pk *PK2G) IsBurned() bool {
	return false
}

func (pk *PK2G) IsFrozen() bool {
	return false
}

func (pk *PK2G) IsParalyzed() bool {
	return false
}

func (pk *PK2G) EV() []uint16 {
	return []uint16{pk.HpEV, pk.AttackEV, pk.DefenseEV, pk.SpeedEV, pk.SpecialEV, pk.SpecialEV}
}

func (pk *PK2G) IV() []uint16 {
	return []uint16{uint16(pk.HpIV), uint16(pk.AttackIV), uint16(pk.DefenseIV), uint16(pk.SpeedIV), uint16(pk.SpecialIV), uint16(pk.SpecialIV)}
}

func (pk *PK2G) Stats() []uint16 {
	return []uint16{uint16(pk.Hp), uint16(pk.Attack), uint16(pk.Defense), uint16(pk.Speed), uint16(pk.SpecialAttack), uint16(pk.SpecialDefense)}
}

func (pk *PK2G) Moves() []uint16 {
	return []uint16{uint16(pk.Move1), uint16(pk.Move2),
		uint16(pk.Move3), uint16(pk.Move4)}
}

func (pk *PK2G) PP() []uint8 {
	return []uint8{pk.Move1PP, pk.Move2PP, pk.Move3PP, pk.Move4PP}
}

func (pk *PK2G) IsShiny() bool {
	i := sort.Search(len(Gen1ShinyAttackIV), func(i int) bool {
		return Gen1ShinyAttackIV[i] >= pk.AttackIV
	})
	if i < len(Gen1ShinyAttackIV) && Gen1ShinyAttackIV[i] == pk.AttackIV {
		return pk.DefenseIV == 10 && pk.SpeedIV == 10 && pk.SpecialIV == 10
	} else {
		return false
	}
}

func (pk *PK2G) IsMissingNo() bool {
	//TODO check id
	return false
}

func (pk *PK2G) Bytes() []byte {
	data := make([]byte, 32)

	data[0] = pk.Species
	data[1] = pk.HeldItem
	data[2] = pk.Move1
	data[3] = pk.Move2
	data[4] = pk.Move3
	data[5] = pk.Move4
	UInt16ToBytes(data, 6, pk.OriginalTrainerID)
	UInt24ToBytes(data, 8, pk.Experience) //TODO repair this
	UInt16ToBytes(data, 11, pk.HpEV)
	UInt16ToBytes(data, 13, pk.AttackEV)
	UInt16ToBytes(data, 15, pk.DefenseEV)
	UInt16ToBytes(data, 17, pk.SpeedEV)
	UInt16ToBytes(data, 19, pk.SpecialEV)
	UInt16ToBytes(data, 21, pk.IVData)
	data[23] = pk.Move1PP
	data[24] = pk.Move2PP
	data[25] = pk.Move3PP
	data[26] = pk.Move4PP
	data[27] = pk.Friendship
	data[28] = pk.Pokerus
	UInt16ToBytes(data, 29, pk.CaughtData)
	data[31] = pk.Lvl

	return data
}

func find1GId(id uint8) (uint8, error) {
	for i := 0; i < len(Gen1IdMap); i++ {
		if Gen1IdMap[i] == id {
			return uint8(i), nil
		}
	}
	return 0, errors.New("incompatible pokemon id")
}

func (pk *PK2G) Downgrade() (PK, error) {
	npk := PK1G{}
	npk.Generation = 1

	id, err := find1GId(pk.Species)
	if err != nil {
		return nil, err
	}

	npk.Species = id
	npk.CurrentHp = 0 //TODO calculate max hp
	npk.Lvl = pk.Lvl
	npk.Status = 0
	npk.Type1 = Gen1Type[pk.Species][0]
	if len(Gen1Type[pk.Species]) > 1 {
		npk.Type2 = Gen1Type[pk.Species][1]
	} else {
		npk.Type2 = npk.Type1
	}
	npk.HeldItem = pk.HeldItem
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
	npk.Move1PP = pk.Move2PP
	npk.Move2PP = pk.Move2PP
	npk.Move3PP = pk.Move2PP
	npk.Move4PP = pk.Move2PP

	npk.CalculateIV()
	pk.CalculateStats()

	return &npk, nil
}

var (
	Gen1Type = make([][]uint8, 256)
)

func InitGen1Types() {
	Gen1Type[63] = []uint8{24}
	Gen1Type[142] = []uint8{5, 2}
	Gen1Type[65] = []uint8{24}
	Gen1Type[24] = []uint8{3}
	Gen1Type[59] = []uint8{20}
	Gen1Type[144] = []uint8{25, 2}
	Gen1Type[15] = []uint8{7, 3}
	Gen1Type[69] = []uint8{22, 3}
	Gen1Type[9] = []uint8{21}
	Gen1Type[1] = []uint8{22, 3}
	Gen1Type[12] = []uint8{7, 2}
	Gen1Type[10] = []uint8{7}
	Gen1Type[113] = []uint8{0}
	Gen1Type[6] = []uint8{20, 2}
	Gen1Type[4] = []uint8{20}
	Gen1Type[5] = []uint8{20}
	Gen1Type[36] = []uint8{}
	Gen1Type[35] = []uint8{}
	Gen1Type[91] = []uint8{21, 25}
	Gen1Type[104] = []uint8{4}
	Gen1Type[87] = []uint8{21, 25}
	Gen1Type[50] = []uint8{4}
	Gen1Type[132] = []uint8{0}
	Gen1Type[85] = []uint8{0, 2}
	Gen1Type[84] = []uint8{0, 2}
	Gen1Type[148] = []uint8{26}
	Gen1Type[149] = []uint8{26, 2}
	Gen1Type[147] = []uint8{26}
	Gen1Type[96] = []uint8{24}
	Gen1Type[51] = []uint8{4}
	Gen1Type[133] = []uint8{0}
	Gen1Type[23] = []uint8{3}
	Gen1Type[125] = []uint8{23}
	Gen1Type[101] = []uint8{23}
	Gen1Type[102] = []uint8{22, 24}
	Gen1Type[103] = []uint8{22, 24}
	Gen1Type[83] = []uint8{0, 2}
	Gen1Type[22] = []uint8{0, 2}
	Gen1Type[136] = []uint8{20}
	Gen1Type[92] = []uint8{8, 3}
	Gen1Type[94] = []uint8{8, 3}
	Gen1Type[74] = []uint8{5, 4}
	Gen1Type[44] = []uint8{22, 3}
	Gen1Type[42] = []uint8{3, 2}
	Gen1Type[118] = []uint8{21}
	Gen1Type[55] = []uint8{21}
	Gen1Type[76] = []uint8{5, 4}
	Gen1Type[75] = []uint8{5, 4}
	Gen1Type[88] = []uint8{3}
	Gen1Type[58] = []uint8{20}
	Gen1Type[130] = []uint8{21, 2}
	Gen1Type[93] = []uint8{8, 3}
	Gen1Type[107] = []uint8{1}
	Gen1Type[106] = []uint8{1}
	Gen1Type[116] = []uint8{21}
	Gen1Type[97] = []uint8{24}
	Gen1Type[2] = []uint8{22, 3}
	Gen1Type[39] = []uint8{0}
	Gen1Type[135] = []uint8{23}
	Gen1Type[124] = []uint8{25, 24}
	Gen1Type[140] = []uint8{5, 21}
	Gen1Type[141] = []uint8{5, 21}
	Gen1Type[64] = []uint8{24}
	Gen1Type[14] = []uint8{7, 3}
	Gen1Type[115] = []uint8{0}
	Gen1Type[99] = []uint8{21}
	Gen1Type[109] = []uint8{3}
	Gen1Type[98] = []uint8{21}
	Gen1Type[131] = []uint8{21, 25}
	Gen1Type[108] = []uint8{0}
	Gen1Type[68] = []uint8{1}
	Gen1Type[67] = []uint8{1}
	Gen1Type[66] = []uint8{1}
	Gen1Type[129] = []uint8{21}
	Gen1Type[126] = []uint8{20}
	Gen1Type[81] = []uint8{23}
	Gen1Type[82] = []uint8{23}
	Gen1Type[56] = []uint8{1}
	Gen1Type[105] = []uint8{4}
	Gen1Type[52] = []uint8{0}
	Gen1Type[11] = []uint8{7}
	Gen1Type[151] = []uint8{24}
	Gen1Type[150] = []uint8{24}
	Gen1Type[146] = []uint8{20, 2}
	Gen1Type[122] = []uint8{24}
	Gen1Type[89] = []uint8{3}
	Gen1Type[34] = []uint8{3, 4}
	Gen1Type[31] = []uint8{3, 4}
	Gen1Type[29] = []uint8{3}
	Gen1Type[32] = []uint8{3}
	Gen1Type[30] = []uint8{3}
	Gen1Type[33] = []uint8{3}
	Gen1Type[38] = []uint8{20}
	Gen1Type[43] = []uint8{22, 3}
	Gen1Type[138] = []uint8{5, 21}
	Gen1Type[139] = []uint8{5, 21}
	Gen1Type[95] = []uint8{5, 4}
	Gen1Type[46] = []uint8{7, 22}
	Gen1Type[47] = []uint8{7, 22}
	Gen1Type[53] = []uint8{0}
	Gen1Type[18] = []uint8{0, 2}
	Gen1Type[17] = []uint8{0, 2}
	Gen1Type[16] = []uint8{0, 2}
	Gen1Type[25] = []uint8{23}
	Gen1Type[127] = []uint8{7}
	Gen1Type[60] = []uint8{21}
	Gen1Type[61] = []uint8{21}
	Gen1Type[62] = []uint8{21, 1}
	Gen1Type[77] = []uint8{20}
	Gen1Type[137] = []uint8{0}
	Gen1Type[57] = []uint8{1}
	Gen1Type[54] = []uint8{21}
	Gen1Type[26] = []uint8{23}
	Gen1Type[78] = []uint8{20}
	Gen1Type[20] = []uint8{0}
	Gen1Type[19] = []uint8{0}
	Gen1Type[112] = []uint8{4, 5}
	Gen1Type[111] = []uint8{4, 5}
	Gen1Type[27] = []uint8{4}
	Gen1Type[28] = []uint8{4}
	Gen1Type[123] = []uint8{7, 2}
	Gen1Type[117] = []uint8{21}
	Gen1Type[119] = []uint8{21}
	Gen1Type[86] = []uint8{21}
	Gen1Type[90] = []uint8{21}
	Gen1Type[80] = []uint8{21, 24}
	Gen1Type[79] = []uint8{21, 24}
	Gen1Type[143] = []uint8{0}
	Gen1Type[21] = []uint8{0, 2}
	Gen1Type[7] = []uint8{21}
	Gen1Type[121] = []uint8{21, 24}
	Gen1Type[120] = []uint8{21}
	Gen1Type[114] = []uint8{22}
	Gen1Type[128] = []uint8{0}
	Gen1Type[72] = []uint8{21, 3}
	Gen1Type[73] = []uint8{21, 3}
	Gen1Type[134] = []uint8{21}
	Gen1Type[49] = []uint8{7, 3}
	Gen1Type[48] = []uint8{7, 3}
	Gen1Type[3] = []uint8{22, 3}
	Gen1Type[71] = []uint8{22, 3}
	Gen1Type[45] = []uint8{22, 3}
	Gen1Type[100] = []uint8{23}
	Gen1Type[37] = []uint8{20}
	Gen1Type[8] = []uint8{21}
	Gen1Type[13] = []uint8{7, 3}
	Gen1Type[70] = []uint8{22, 3}
	Gen1Type[110] = []uint8{3}
	Gen1Type[40] = []uint8{0}
	Gen1Type[145] = []uint8{23, 2}
	Gen1Type[41] = []uint8{3, 2}
}

type Gen2Base struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Hp             int    `json:"hp"`
	Attack         int    `json:"attack"`
	Defense        int    `json:"defense"`
	Speed          int    `json:"speed"`
	SpecialAttack  int    `json:"special-attack"`
	SpecialDefense int    `json:"special-defense"`
}

var (
	Gen2BaseStats = [][]Gen2Base{
		{{1, "Bulbasaur", 45, 49, 49, 45, 65, 65}},
		{{2, "Ivysaur", 60, 62, 63, 60, 80, 80}},
		{{3, "Venusaur", 80, 82, 83, 80, 100, 100}},
		{{4, "Charmander", 39, 52, 43, 65, 60, 50}},
		{{5, "Charmeleon", 58, 64, 58, 80, 80, 65}},
		{{6, "Charizard", 78, 84, 78, 100, 109, 85}},
		{{7, "Squirtle", 44, 48, 65, 43, 50, 64}},
		{{8, "Wartortle", 59, 63, 80, 58, 65, 80}},
		{{9, "Blastoise", 79, 83, 100, 78, 85, 105}},
		{{10, "Caterpie", 45, 30, 35, 45, 20, 20}},
		{{11, "Metapod", 50, 20, 55, 30, 25, 25}},
		{{12, "Butterfree", 60, 45, 50, 70, 80, 80}},
		{{13, "Weedle", 40, 35, 30, 50, 20, 20}},
		{{14, "Kakuna", 45, 25, 50, 35, 25, 25}},
		{{15, "Beedrill", 65, 80, 40, 75, 45, 80}},
		{{16, "Pidgey", 40, 45, 40, 56, 35, 35}},
		{{17, "Pidgeotto", 63, 60, 55, 71, 50, 50}},
		{{18, "Pidgeot", 83, 80, 75, 91, 70, 70}},
		{{19, "Rattata", 30, 56, 35, 72, 25, 35}},
		{{20, "Raticate", 55, 81, 60, 97, 50, 70}},
		{{21, "Spearow", 40, 60, 30, 70, 31, 31}},
		{{22, "Fearow", 65, 90, 65, 100, 61, 61}},
		{{23, "Ekans", 35, 60, 44, 55, 40, 54}},
		{{24, "Arbok", 60, 85, 69, 80, 65, 79}},
		{{25, "Pikachu", 35, 55, 30, 90, 50, 40}},
		{{26, "Raichu", 60, 90, 55, 100, 90, 80}},
		{{27, "Sandshrew", 50, 75, 85, 40, 20, 30}},
		{{28, "Sandslash", 75, 100, 110, 65, 45, 55}},
		{{29, "Nidoran♀", 55, 47, 52, 41, 40, 40}},
		{{30, "Nidorina", 70, 62, 67, 56, 55, 55}},
		{{31, "Nidoqueen", 90, 82, 87, 76, 75, 85}},
		{{32, "Nidoran♂", 46, 57, 40, 50, 40, 40}},
		{{33, "Nidorino", 61, 72, 57, 65, 55, 55}},
		{{34, "Nidoking", 81, 92, 77, 85, 85, 75}},
		{{35, "Clefairy", 70, 45, 48, 35, 60, 65}},
		{{36, "Clefable", 95, 70, 73, 60, 85, 90}},
		{{37, "Vulpix", 38, 41, 40, 65, 50, 65}},
		{{38, "Ninetales", 73, 76, 75, 100, 81, 100}},
		{{39, "Jigglypuff", 115, 45, 20, 20, 45, 25}},
		{{40, "Wigglytuff", 140, 70, 45, 45, 75, 50}},
		{{41, "Zubat", 40, 45, 35, 55, 30, 40}},
		{{42, "Golbat", 75, 80, 70, 90, 65, 75}},
		{{43, "Oddish", 45, 50, 55, 30, 75, 65}},
		{{44, "Gloom", 60, 65, 70, 40, 85, 75}},
		{{45, "Vileplume", 75, 80, 85, 50, 100, 90}},
		{{46, "Paras", 35, 70, 55, 25, 45, 55}},
		{{47, "Parasect", 60, 95, 80, 30, 60, 80}},
		{{48, "Venonat", 60, 55, 50, 45, 40, 55}},
		{{49, "Venomoth", 70, 65, 60, 90, 90, 75}},
		{{50, "Diglett", 10, 55, 25, 95, 35, 45}},
		{{51, "Dugtrio", 35, 80, 50, 120, 50, 70}},
		{{52, "Meowth", 40, 45, 35, 90, 40, 40}},
		{{53, "Persian", 65, 70, 60, 115, 65, 65}},
		{{54, "Psyduck", 50, 52, 48, 55, 65, 50}},
		{{55, "Golduck", 80, 82, 78, 85, 95, 80}},
		{{56, "Mankey", 40, 80, 35, 70, 35, 45}},
		{{57, "Primeape", 65, 105, 60, 95, 60, 70}},
		{{58, "Growlithe", 55, 70, 45, 60, 70, 50}},
		{{59, "Arcanine", 90, 110, 80, 95, 100, 80}},
		{{60, "Poliwag", 40, 50, 40, 90, 40, 40}},
		{{61, "Poliwhirl", 65, 65, 65, 90, 50, 50}},
		{{62, "Poliwrath", 90, 85, 95, 70, 70, 90}},
		{{63, "Abra", 25, 20, 15, 90, 105, 55}},
		{{64, "Kadabra", 40, 35, 30, 105, 120, 70}},
		{{65, "Alakazam", 55, 50, 45, 120, 135, 85}},
		{{66, "Machop", 70, 80, 50, 35, 35, 35}},
		{{67, "Machoke", 80, 100, 70, 45, 50, 60}},
		{{68, "Machamp", 90, 130, 80, 55, 65, 85}},
		{{69, "Bellsprout", 50, 75, 35, 40, 70, 30}},
		{{70, "Weepinbell", 65, 90, 50, 55, 85, 45}},
		{{71, "Victreebel", 80, 105, 65, 70, 100, 60}},
		{{72, "Tentacool", 40, 40, 35, 70, 50, 100}},
		{{73, "Tentacruel", 80, 70, 65, 100, 80, 120}},
		{{74, "Geodude", 40, 80, 100, 20, 30, 30}},
		{{75, "Graveler", 55, 95, 115, 35, 45, 45}},
		{{76, "Golem", 80, 110, 130, 45, 55, 65}},
		{{77, "Ponyta", 50, 85, 55, 90, 65, 65}},
		{{78, "Rapidash", 65, 100, 70, 105, 80, 80}},
		{{79, "Slowpoke", 90, 65, 65, 15, 40, 40}},
		{{80, "Slowbro", 95, 75, 110, 30, 100, 80}},
		{{81, "Magnemite", 25, 35, 70, 45, 95, 55}},
		{{82, "Magneton", 50, 60, 95, 70, 120, 70}},
		{{83, "Farfetch&#x27;d", 52, 65, 55, 60, 58, 62}},
		{{84, "Doduo", 35, 85, 45, 75, 35, 35}},
		{{85, "Dodrio", 60, 110, 70, 100, 60, 60}},
		{{86, "Seel", 65, 45, 55, 45, 45, 70}},
		{{87, "Dewgong", 90, 70, 80, 70, 70, 95}},
		{{88, "Grimer", 80, 80, 50, 25, 40, 50}},
		{{89, "Muk", 105, 105, 75, 50, 65, 100}},
		{{90, "Shellder", 30, 65, 100, 40, 45, 25}},
		{{91, "Cloyster", 50, 95, 180, 70, 85, 45}},
		{{92, "Gastly", 30, 35, 30, 80, 100, 35}},
		{{93, "Haunter", 45, 50, 45, 95, 115, 55}},
		{{94, "Gengar", 60, 65, 60, 110, 130, 75}},
		{{95, "Onix", 35, 45, 160, 70, 30, 45}},
		{{96, "Drowzee", 60, 48, 45, 42, 43, 90}},
		{{97, "Hypno", 85, 73, 70, 67, 73, 115}},
		{{98, "Krabby", 30, 105, 90, 50, 25, 25}},
		{{99, "Kingler", 55, 130, 115, 75, 50, 50}},
		{{100, "Voltorb", 40, 30, 50, 100, 55, 55}},
		{{101, "Electrode", 60, 50, 70, 140, 80, 80}},
		{{102, "Exeggcute", 60, 40, 80, 40, 60, 45}},
		{{103, "Exeggutor", 95, 95, 85, 55, 125, 65}},
		{{104, "Cubone", 50, 50, 95, 35, 40, 50}},
		{{105, "Marowak", 60, 80, 110, 45, 50, 80}},
		{{106, "Hitmonlee", 50, 120, 53, 87, 35, 110}},
		{{107, "Hitmonchan", 50, 105, 79, 76, 35, 110}},
		{{108, "Lickitung", 90, 55, 75, 30, 60, 75}},
		{{109, "Koffing", 40, 65, 95, 35, 60, 45}},
		{{110, "Weezing", 65, 90, 120, 60, 85, 70}},
		{{111, "Rhyhorn", 80, 85, 95, 25, 30, 30}},
		{{112, "Rhydon", 105, 130, 120, 40, 45, 45}},
		{{113, "Chansey", 250, 5, 5, 50, 35, 105}},
		{{114, "Tangela", 65, 55, 115, 60, 100, 40}},
		{{115, "Kangaskhan", 105, 95, 80, 90, 40, 80}},
		{{116, "Horsea", 30, 40, 70, 60, 70, 25}},
		{{117, "Seadra", 55, 65, 95, 85, 95, 45}},
		{{118, "Goldeen", 45, 67, 60, 63, 35, 50}},
		{{119, "Seaking", 80, 92, 65, 68, 65, 80}},
		{{120, "Staryu", 30, 45, 55, 85, 70, 55}},
		{{121, "Starmie", 60, 75, 85, 115, 100, 85}},
		{{122, "Mr. Mime", 40, 45, 65, 90, 100, 120}},
		{{123, "Scyther", 70, 110, 80, 105, 55, 80}},
		{{124, "Jynx", 65, 50, 35, 95, 115, 95}},
		{{125, "Electabuzz", 65, 83, 57, 105, 95, 85}},
		{{126, "Magmar", 65, 95, 57, 93, 100, 85}},
		{{127, "Pinsir", 65, 125, 100, 85, 55, 70}},
		{{128, "Tauros", 75, 100, 95, 110, 40, 70}},
		{{129, "Magikarp", 20, 10, 55, 80, 15, 20}},
		{{130, "Gyarados", 95, 125, 79, 81, 60, 100}},
		{{131, "Lapras", 130, 85, 80, 60, 85, 95}},
		{{132, "Ditto", 48, 48, 48, 48, 48, 48}},
		{{133, "Eevee", 55, 55, 50, 55, 45, 65}},
		{{134, "Vaporeon", 130, 65, 60, 65, 110, 95}},
		{{135, "Jolteon", 65, 65, 60, 130, 110, 95}},
		{{136, "Flareon", 65, 130, 60, 65, 95, 110}},
		{{137, "Porygon", 65, 60, 70, 40, 85, 75}},
		{{138, "Omanyte", 35, 40, 100, 35, 90, 55}},
		{{139, "Omastar", 70, 60, 125, 55, 115, 70}},
		{{140, "Kabuto", 30, 80, 90, 55, 55, 45}},
		{{141, "Kabutops", 60, 115, 105, 80, 65, 70}},
		{{142, "Aerodactyl", 80, 105, 65, 130, 60, 75}},
		{{143, "Snorlax", 160, 110, 65, 30, 65, 110}},
		{{144, "Articuno", 90, 85, 100, 85, 95, 125}},
		{{145, "Zapdos", 90, 90, 85, 100, 125, 90}},
		{{146, "Moltres", 90, 100, 90, 90, 125, 85}},
		{{147, "Dratini", 41, 64, 45, 50, 50, 50}},
		{{148, "Dragonair", 61, 84, 65, 70, 70, 70}},
		{{149, "Dragonite", 91, 134, 95, 80, 100, 100}},
		{{150, "Mewtwo", 106, 110, 90, 130, 154, 90}},
		{{151, "Mew", 100, 100, 100, 100, 100, 100}},
		{{152, "Chikorita", 45, 49, 65, 45, 49, 65}},
		{{153, "Bayleef", 60, 62, 80, 60, 63, 80}},
		{{154, "Meganium", 80, 82, 100, 80, 83, 100}},
		{{155, "Cyndaquil", 39, 52, 43, 65, 60, 50}},
		{{156, "Quilava", 58, 64, 58, 80, 80, 65}},
		{{157, "Typhlosion", 78, 84, 78, 100, 109, 85}},
		{{158, "Totodile", 50, 65, 64, 43, 44, 48}},
		{{159, "Croconaw", 65, 80, 80, 58, 59, 63}},
		{{160, "Feraligatr", 85, 105, 100, 78, 79, 83}},
		{{161, "Sentret", 35, 46, 34, 20, 35, 45}},
		{{162, "Furret", 85, 76, 64, 90, 45, 55}},
		{{163, "Hoothoot", 60, 30, 30, 50, 36, 56}},
		{{164, "Noctowl", 100, 50, 50, 70, 76, 96}},
		{{165, "Ledyba", 40, 20, 30, 55, 40, 80}},
		{{166, "Ledian", 55, 35, 50, 85, 55, 110}},
		{{167, "Spinarak", 40, 60, 40, 30, 40, 40}},
		{{168, "Ariados", 70, 90, 70, 40, 60, 60}},
		{{169, "Crobat", 85, 90, 80, 130, 70, 80}},
		{{170, "Chinchou", 75, 38, 38, 67, 56, 56}},
		{{171, "Lanturn", 125, 58, 58, 67, 76, 76}},
		{{172, "Pichu", 20, 40, 15, 60, 35, 35}},
		{{173, "Cleffa", 50, 25, 28, 15, 45, 55}},
		{{174, "Igglybuff", 90, 30, 15, 15, 40, 20}},
		{{175, "Togepi", 35, 20, 65, 20, 40, 65}},
		{{176, "Togetic", 55, 40, 85, 40, 80, 105}},
		{{177, "Natu", 40, 50, 45, 70, 70, 45}},
		{{178, "Xatu", 65, 75, 70, 95, 95, 70}},
		{{179, "Mareep", 55, 40, 40, 35, 65, 45}},
		{{180, "Flaaffy", 70, 55, 55, 45, 80, 60}},
		{{181, "Ampharos", 90, 75, 75, 55, 115, 90}},
		{{182, "Bellossom", 75, 80, 85, 50, 90, 100}},
		{{183, "Marill", 70, 20, 50, 40, 20, 50}},
		{{184, "Azumarill", 100, 50, 80, 50, 50, 80}},
		{{185, "Sudowoodo", 70, 100, 115, 30, 30, 65}},
		{{186, "Politoed", 90, 75, 75, 70, 90, 100}},
		{{187, "Hoppip", 35, 35, 40, 50, 35, 55}},
		{{188, "Skiploom", 55, 45, 50, 80, 45, 65}},
		{{189, "Jumpluff", 75, 55, 70, 110, 55, 85}},
		{{190, "Aipom", 55, 70, 55, 85, 40, 55}},
		{{191, "Sunkern", 30, 30, 30, 30, 30, 30}},
		{{192, "Sunflora", 75, 75, 55, 30, 105, 85}},
		{{193, "Yanma", 65, 65, 45, 95, 75, 45}},
		{{194, "Wooper", 55, 45, 45, 15, 25, 25}},
		{{195, "Quagsire", 95, 85, 85, 35, 65, 65}},
		{{196, "Espeon", 65, 65, 60, 110, 130, 95}},
		{{197, "Umbreon", 95, 65, 110, 65, 60, 130}},
		{{198, "Murkrow", 60, 85, 42, 91, 85, 42}},
		{{199, "Slowking", 95, 75, 80, 30, 100, 110}},
		{{200, "Misdreavus", 60, 60, 60, 85, 85, 85}},
		{{201, "Unown", 48, 72, 48, 48, 72, 48}},
		{{202, "Wobbuffet", 190, 33, 58, 33, 33, 58}},
		{{203, "Girafarig", 70, 80, 65, 85, 90, 65}},
		{{204, "Pineco", 50, 65, 90, 15, 35, 35}},
		{{205, "Forretress", 75, 90, 140, 40, 60, 60}},
		{{206, "Dunsparce", 100, 70, 70, 45, 65, 65}},
		{{207, "Gligar", 65, 75, 105, 85, 35, 65}},
		{{208, "Steelix", 75, 85, 200, 30, 55, 65}},
		{{209, "Snubbull", 60, 80, 50, 30, 40, 40}},
		{{210, "Granbull", 90, 120, 75, 45, 60, 60}},
		{{211, "Qwilfish", 65, 95, 75, 85, 55, 55}},
		{{212, "Scizor", 70, 130, 100, 65, 55, 80}},
		{{213, "Shuckle", 20, 10, 230, 5, 10, 230}},
		{{214, "Heracross", 80, 125, 75, 85, 40, 95}},
		{{215, "Sneasel", 55, 95, 55, 115, 35, 75}},
		{{216, "Teddiursa", 60, 80, 50, 40, 50, 50}},
		{{217, "Ursaring", 90, 130, 75, 55, 75, 75}},
		{{218, "Slugma", 40, 40, 40, 20, 70, 40}},
		{{219, "Magcargo", 50, 50, 120, 30, 80, 80}},
		{{220, "Swinub", 50, 50, 40, 50, 30, 30}},
		{{221, "Piloswine", 100, 100, 80, 50, 60, 60}},
		{{222, "Corsola", 55, 55, 85, 35, 65, 85}},
		{{223, "Remoraid", 35, 65, 35, 65, 65, 35}},
		{{224, "Octillery", 75, 105, 75, 45, 105, 75}},
		{{225, "Delibird", 45, 55, 45, 75, 65, 45}},
		{{226, "Mantine", 65, 40, 70, 70, 80, 140}},
		{{227, "Skarmory", 65, 80, 140, 70, 40, 70}},
		{{228, "Houndour", 45, 60, 30, 65, 80, 50}},
		{{229, "Houndoom", 75, 90, 50, 95, 110, 80}},
		{{230, "Kingdra", 75, 95, 95, 85, 95, 95}},
		{{231, "Phanpy", 90, 60, 60, 40, 40, 40}},
		{{232, "Donphan", 90, 120, 120, 50, 60, 60}},
		{{233, "Porygon2", 85, 80, 90, 60, 105, 95}},
		{{234, "Stantler", 73, 95, 62, 85, 85, 65}},
		{{235, "Smeargle", 55, 20, 35, 75, 20, 45}},
		{{236, "Tyrogue", 35, 35, 35, 35, 35, 35}},
		{{237, "Hitmontop", 50, 95, 95, 70, 35, 110}},
		{{238, "Smoochum", 45, 30, 15, 65, 85, 65}},
		{{239, "Elekid", 45, 63, 37, 95, 65, 55}},
		{{240, "Magby", 45, 75, 37, 83, 70, 55}},
		{{241, "Miltank", 95, 80, 105, 100, 40, 70}},
		{{242, "Blissey", 255, 10, 10, 55, 75, 135}},
		{{243, "Raikou", 90, 85, 75, 115, 115, 100}},
		{{244, "Entei", 115, 115, 85, 100, 90, 75}},
		{{245, "Suicune", 100, 75, 115, 85, 90, 115}},
		{{246, "Larvitar", 50, 64, 50, 41, 45, 50}},
		{{247, "Pupitar", 70, 84, 70, 51, 65, 70}},
		{{248, "Tyranitar", 100, 134, 110, 61, 95, 100}},
		{{249, "Lugia", 106, 90, 130, 110, 90, 154}},
		{{250, "Ho-Oh", 106, 130, 90, 90, 110, 154}},
		{{251, "Celebi", 100, 100, 100, 100, 100, 100}},
		{{252, "Treecko", 40, 45, 35, 70, 65, 55}},
		{{253, "Grovyle", 50, 65, 45, 95, 85, 65}},
		{{254, "Sceptile", 70, 85, 65, 120, 105, 85}},
		{{255, "Torchic", 45, 60, 40, 45, 70, 50}},
		{{256, "Combusken", 60, 85, 60, 55, 85, 60}},
		{{257, "Blaziken", 80, 120, 70, 80, 110, 70}},
		{{258, "Mudkip", 50, 70, 50, 40, 50, 50}},
		{{259, "Marshtomp", 70, 85, 70, 50, 60, 70}},
		{{260, "Swampert", 100, 110, 90, 60, 85, 90}},
		{{261, "Poochyena", 35, 55, 35, 35, 30, 30}},
		{{262, "Mightyena", 70, 90, 70, 70, 60, 60}},
		{{263, "Zigzagoon", 38, 30, 41, 60, 30, 41}},
		{{264, "Linoone", 78, 70, 61, 100, 50, 61}},
		{{265, "Wurmple", 45, 45, 35, 20, 20, 30}},
		{{266, "Silcoon", 50, 35, 55, 15, 25, 25}},
		{{267, "Beautifly", 60, 70, 50, 65, 90, 50}},
		{{268, "Cascoon", 50, 35, 55, 15, 25, 25}},
		{{269, "Dustox", 60, 50, 70, 65, 50, 90}},
		{{270, "Lotad", 40, 30, 30, 30, 40, 50}},
		{{271, "Lombre", 60, 50, 50, 50, 60, 70}},
		{{272, "Ludicolo", 80, 70, 70, 70, 90, 100}},
		{{273, "Seedot", 40, 40, 50, 30, 30, 30}},
		{{274, "Nuzleaf", 70, 70, 40, 60, 60, 40}},
		{{275, "Shiftry", 90, 100, 60, 80, 90, 60}},
		{{276, "Taillow", 40, 55, 30, 85, 30, 30}},
		{{277, "Swellow", 60, 85, 60, 125, 50, 50}},
		{{278, "Wingull", 40, 30, 30, 85, 55, 30}},
		{{279, "Pelipper", 60, 50, 100, 65, 85, 70}},
		{{280, "Ralts", 28, 25, 25, 40, 45, 35}},
		{{281, "Kirlia", 38, 35, 35, 50, 65, 55}},
		{{282, "Gardevoir", 68, 65, 65, 80, 125, 115}},
		{{283, "Surskit", 40, 30, 32, 65, 50, 52}},
		{{284, "Masquerain", 70, 60, 62, 60, 80, 82}},
		{{285, "Shroomish", 60, 40, 60, 35, 40, 60}},
		{{286, "Breloom", 60, 130, 80, 70, 60, 60}},
		{{287, "Slakoth", 60, 60, 60, 30, 35, 35}},
		{{288, "Vigoroth", 80, 80, 80, 90, 55, 55}},
		{{289, "Slaking", 150, 160, 100, 100, 95, 65}},
		{{290, "Nincada", 31, 45, 90, 40, 30, 30}},
		{{291, "Ninjask", 61, 90, 45, 160, 50, 50}},
		{{292, "Shedinja", 1, 90, 45, 40, 30, 30}},
		{{293, "Whismur", 64, 51, 23, 28, 51, 23}},
		{{294, "Loudred", 84, 71, 43, 48, 71, 43}},
		{{295, "Exploud", 104, 91, 63, 68, 91, 63}},
		{{296, "Makuhita", 72, 60, 30, 25, 20, 30}},
		{{297, "Hariyama", 144, 120, 60, 50, 40, 60}},
		{{298, "Azurill", 50, 20, 40, 20, 20, 40}},
		{{299, "Nosepass", 30, 45, 135, 30, 45, 90}},
		{{300, "Skitty", 50, 45, 45, 50, 35, 35}},
		{{301, "Delcatty", 70, 65, 65, 70, 55, 55}},
		{{302, "Sableye", 50, 75, 75, 50, 65, 65}},
		{{303, "Mawile", 50, 85, 85, 50, 55, 55}},
		{{304, "Aron", 50, 70, 100, 30, 40, 40}},
		{{305, "Lairon", 60, 90, 140, 40, 50, 50}},
		{{306, "Aggron", 70, 110, 180, 50, 60, 60}},
		{{307, "Meditite", 30, 40, 55, 60, 40, 55}},
		{{308, "Medicham", 60, 60, 75, 80, 60, 75}},
		{{309, "Electrike", 40, 45, 40, 65, 65, 40}},
		{{310, "Manectric", 70, 75, 60, 105, 105, 60}},
		{{311, "Plusle", 60, 50, 40, 95, 85, 75}},
		{{312, "Minun", 60, 40, 50, 95, 75, 85}},
		{{313, "Volbeat", 65, 73, 55, 85, 47, 75}},
		{{314, "Illumise", 65, 47, 55, 85, 73, 75}},
		{{315, "Roselia", 50, 60, 45, 65, 100, 80}},
		{{316, "Gulpin", 70, 43, 53, 40, 43, 53}},
		{{317, "Swalot", 100, 73, 83, 55, 73, 83}},
		{{318, "Carvanha", 45, 90, 20, 65, 65, 20}},
		{{319, "Sharpedo", 70, 120, 40, 95, 95, 40}},
		{{320, "Wailmer", 130, 70, 35, 60, 70, 35}},
		{{321, "Wailord", 170, 90, 45, 60, 90, 45}},
		{{322, "Numel", 60, 60, 40, 35, 65, 45}},
		{{323, "Camerupt", 70, 100, 70, 40, 105, 75}},
		{{324, "Torkoal", 70, 85, 140, 20, 85, 70}},
		{{325, "Spoink", 60, 25, 35, 60, 70, 80}},
		{{326, "Grumpig", 80, 45, 65, 80, 90, 110}},
		{{327, "Spinda", 60, 60, 60, 60, 60, 60}},
		{{328, "Trapinch", 45, 100, 45, 10, 45, 45}},
		{{329, "Vibrava", 50, 70, 50, 70, 50, 50}},
		{{330, "Flygon", 80, 100, 80, 100, 80, 80}},
		{{331, "Cacnea", 50, 85, 40, 35, 85, 40}},
		{{332, "Cacturne", 70, 115, 60, 55, 115, 60}},
		{{333, "Swablu", 45, 40, 60, 50, 40, 75}},
		{{334, "Altaria", 75, 70, 90, 80, 70, 105}},
		{{335, "Zangoose", 73, 115, 60, 90, 60, 60}},
		{{336, "Seviper", 73, 100, 60, 65, 100, 60}},
		{{337, "Lunatone", 70, 55, 65, 70, 95, 85}},
		{{338, "Solrock", 70, 95, 85, 70, 55, 65}},
		{{339, "Barboach", 50, 48, 43, 60, 46, 41}},
		{{340, "Whiscash", 110, 78, 73, 60, 76, 71}},
		{{341, "Corphish", 43, 80, 65, 35, 50, 35}},
		{{342, "Crawdaunt", 63, 120, 85, 55, 90, 55}},
		{{343, "Baltoy", 40, 40, 55, 55, 40, 70}},
		{{344, "Claydol", 60, 70, 105, 75, 70, 120}},
		{{345, "Lileep", 66, 41, 77, 23, 61, 87}},
		{{346, "Cradily", 86, 81, 97, 43, 81, 107}},
		{{347, "Anorith", 45, 95, 50, 75, 40, 50}},
		{{348, "Armaldo", 75, 125, 100, 45, 70, 80}},
		{{349, "Feebas", 20, 15, 20, 80, 10, 55}},
		{{350, "Milotic", 95, 60, 79, 81, 100, 125}},
		{{351, "Castform", 70, 70, 70, 70, 70, 70}},
		{{352, "Kecleon", 60, 90, 70, 40, 60, 120}},
		{{353, "Shuppet", 44, 75, 35, 45, 63, 33}},
		{{354, "Banette", 64, 115, 65, 65, 83, 63}},
		{{355, "Duskull", 20, 40, 90, 25, 30, 90}},
		{{356, "Dusclops", 40, 70, 130, 25, 60, 130}},
		{{357, "Tropius", 99, 68, 83, 51, 72, 87}},
		{{358, "Chimecho", 65, 50, 70, 65, 95, 80}},
		{{359, "Absol", 65, 130, 60, 75, 75, 60}},
		{{360, "Wynaut", 95, 23, 48, 23, 23, 48}},
		{{361, "Snorunt", 50, 50, 50, 50, 50, 50}},
		{{362, "Glalie", 80, 80, 80, 80, 80, 80}},
		{{363, "Spheal", 70, 40, 50, 25, 55, 50}},
		{{364, "Sealeo", 90, 60, 70, 45, 75, 70}},
		{{365, "Walrein", 110, 80, 90, 65, 95, 90}},
		{{366, "Clamperl", 35, 64, 85, 32, 74, 55}},
		{{367, "Huntail", 55, 104, 105, 52, 94, 75}},
		{{368, "Gorebyss", 55, 84, 105, 52, 114, 75}},
		{{369, "Relicanth", 100, 90, 130, 55, 45, 65}},
		{{370, "Luvdisc", 43, 30, 55, 97, 40, 65}},
		{{371, "Bagon", 45, 75, 60, 50, 40, 30}},
		{{372, "Shelgon", 65, 95, 100, 50, 60, 50}},
		{{373, "Salamence", 95, 135, 80, 100, 110, 80}},
		{{374, "Beldum", 40, 55, 80, 30, 35, 60}},
		{{375, "Metang", 60, 75, 100, 50, 55, 80}},
		{{376, "Metagross", 80, 135, 130, 70, 95, 90}},
		{{377, "Regirock", 80, 100, 200, 50, 50, 100}},
		{{378, "Regice", 80, 50, 100, 50, 100, 200}},
		{{379, "Registeel", 80, 75, 150, 50, 75, 150}},
		{{380, "Latias", 80, 80, 90, 110, 110, 130}},
		{{381, "Latios", 80, 90, 80, 110, 130, 110}},
		{{382, "Kyogre", 100, 100, 90, 90, 150, 140}},
		{{383, "Groudon", 100, 150, 140, 90, 100, 90}},
		{{384, "Rayquaza", 105, 150, 90, 95, 150, 90}},
		{{385, "Jirachi", 100, 100, 100, 100, 100, 100}},
		{{386, "Deoxys (Normal Forme)", 50, 150, 50, 150, 150, 50},
			{386, "Deoxys (Attack Forme)", 50, 180, 20, 150, 180, 20},
			{386, "Deoxys (Defense Forme)", 50, 70, 160, 90, 70, 160},
			{386, "Deoxys (Speed Forme)", 50, 95, 90, 180, 95, 90}},
		{{387, "Turtwig", 55, 68, 64, 31, 45, 55}},
		{{388, "Grotle", 75, 89, 85, 36, 55, 65}},
		{{389, "Torterra", 95, 109, 105, 56, 75, 85}},
		{{390, "Chimchar", 44, 58, 44, 61, 58, 44}},
		{{391, "Monferno", 64, 78, 52, 81, 78, 52}},
		{{392, "Infernape", 76, 104, 71, 108, 104, 71}},
		{{393, "Piplup", 53, 51, 53, 40, 61, 56}},
		{{394, "Prinplup", 64, 66, 68, 50, 81, 76}},
		{{395, "Empoleon", 84, 86, 88, 60, 111, 101}},
		{{396, "Starly", 40, 55, 30, 60, 30, 30}},
		{{397, "Staravia", 55, 75, 50, 80, 40, 40}},
		{{398, "Staraptor", 85, 120, 70, 100, 50, 50}},
		{{399, "Bidoof", 59, 45, 40, 31, 35, 40}},
		{{400, "Bibarel", 79, 85, 60, 71, 55, 60}},
		{{401, "Kricketot", 37, 25, 41, 25, 25, 41}},
		{{402, "Kricketune", 77, 85, 51, 65, 55, 51}},
		{{403, "Shinx", 45, 65, 34, 45, 40, 34}},
		{{404, "Luxio", 60, 85, 49, 60, 60, 49}},
		{{405, "Luxray", 80, 120, 79, 70, 95, 79}},
		{{406, "Budew", 40, 30, 35, 55, 50, 70}},
		{{407, "Roserade", 60, 70, 55, 90, 125, 105}},
		{{408, "Cranidos", 67, 125, 40, 58, 30, 30}},
		{{409, "Rampardos", 97, 165, 60, 58, 65, 50}},
		{{410, "Shieldon", 30, 42, 118, 30, 42, 88}},
		{{411, "Bastiodon", 60, 52, 168, 30, 47, 138}},
		{{412, "Burmy", 40, 29, 45, 36, 29, 45}},
		{{413, "Wormadam (Plant Cloak)", 60, 59, 85, 36, 79, 105},
			{413, "Wormadam (Sandy Cloak)", 60, 79, 105, 36, 59, 85},
			{413, "Wormadam (Trash Cloak)", 60, 69, 95, 36, 69, 95}},
		{{414, "Mothim", 70, 94, 50, 66, 94, 50}},
		{{415, "Combee", 30, 30, 42, 70, 30, 42}},
		{{416, "Vespiquen", 70, 80, 102, 40, 80, 102}},
		{{417, "Pachirisu", 60, 45, 70, 95, 45, 90}},
		{{418, "Buizel", 55, 65, 35, 85, 60, 30}},
		{{419, "Floatzel", 85, 105, 55, 115, 85, 50}},
		{{420, "Cherubi", 45, 35, 45, 35, 62, 53}},
		{{421, "Cherrim", 70, 60, 70, 85, 87, 78}},
		{{422, "Shellos", 76, 48, 48, 34, 57, 62}},
		{{423, "Gastrodon", 111, 83, 68, 39, 92, 82}},
		{{424, "Ambipom", 75, 100, 66, 115, 60, 66}},
		{{425, "Drifloon", 90, 50, 34, 70, 60, 44}},
		{{426, "Drifblim", 150, 80, 44, 80, 90, 54}},
		{{427, "Buneary", 55, 66, 44, 85, 44, 56}},
		{{428, "Lopunny", 65, 76, 84, 105, 54, 96}},
		{{429, "Mismagius", 60, 60, 60, 105, 105, 105}},
		{{430, "Honchkrow", 100, 125, 52, 71, 105, 52}},
		{{431, "Glameow", 49, 55, 42, 85, 42, 37}},
		{{432, "Purugly", 71, 82, 64, 112, 64, 59}},
		{{433, "Chingling", 45, 30, 50, 45, 65, 50}},
		{{434, "Stunky", 63, 63, 47, 74, 41, 41}},
		{{435, "Skuntank", 103, 93, 67, 84, 71, 61}},
		{{436, "Bronzor", 57, 24, 86, 23, 24, 86}},
		{{437, "Bronzong", 67, 89, 116, 33, 79, 116}},
		{{438, "Bonsly", 50, 80, 95, 10, 10, 45}},
		{{439, "Mime Jr.", 20, 25, 45, 60, 70, 90}},
		{{440, "Happiny", 100, 5, 5, 30, 15, 65}},
		{{441, "Chatot", 76, 65, 45, 91, 92, 42}},
		{{442, "Spiritomb", 50, 92, 108, 35, 92, 108}},
		{{443, "Gible", 58, 70, 45, 42, 40, 45}},
		{{444, "Gabite", 68, 90, 65, 82, 50, 55}},
		{{445, "Garchomp", 108, 130, 95, 102, 80, 85}},
		{{446, "Munchlax", 135, 85, 40, 5, 40, 85}},
		{{447, "Riolu", 40, 70, 40, 60, 35, 40}},
		{{448, "Lucario", 70, 110, 70, 90, 115, 70}},
		{{449, "Hippopotas", 68, 72, 78, 32, 38, 42}},
		{{450, "Hippowdon", 108, 112, 118, 47, 68, 72}},
		{{451, "Skorupi", 40, 50, 90, 65, 30, 55}},
		{{452, "Drapion", 70, 90, 110, 95, 60, 75}},
		{{453, "Croagunk", 48, 61, 40, 50, 61, 40}},
		{{454, "Toxicroak", 83, 106, 65, 85, 86, 65}},
		{{455, "Carnivine", 74, 100, 72, 46, 90, 72}},
		{{456, "Finneon", 49, 49, 56, 66, 49, 61}},
		{{457, "Lumineon", 69, 69, 76, 91, 69, 86}},
		{{458, "Mantyke", 45, 20, 50, 50, 60, 120}},
		{{459, "Snover", 60, 62, 50, 40, 62, 60}},
		{{460, "Abomasnow", 90, 92, 75, 60, 92, 85}},
		{{461, "Weavile", 70, 120, 65, 125, 45, 85}},
		{{462, "Magnezone", 70, 70, 115, 60, 130, 90}},
		{{463, "Lickilicky", 110, 85, 95, 50, 80, 95}},
		{{464, "Rhyperior", 115, 140, 130, 40, 55, 55}},
		{{465, "Tangrowth", 100, 100, 125, 50, 110, 50}},
		{{466, "Electivire", 75, 123, 67, 95, 95, 85}},
		{{467, "Magmortar", 75, 95, 67, 83, 125, 95}},
		{{468, "Togekiss", 85, 50, 95, 80, 120, 115}},
		{{469, "Yanmega", 86, 76, 86, 95, 116, 56}},
		{{470, "Leafeon", 65, 110, 130, 95, 60, 65}},
		{{471, "Glaceon", 65, 60, 110, 65, 130, 95}},
		{{472, "Gliscor", 75, 95, 125, 95, 45, 75}},
		{{473, "Mamoswine", 110, 130, 80, 80, 70, 60}},
		{{474, "Porygon-Z", 85, 80, 70, 90, 135, 75}},
		{{475, "Gallade", 68, 125, 65, 80, 65, 115}},
		{{476, "Probopass", 60, 55, 145, 40, 75, 150}},
		{{477, "Dusknoir", 45, 100, 135, 45, 65, 135}},
		{{478, "Froslass", 70, 80, 70, 110, 80, 70}},
		{{479, "Rotom (Normal Rotom)", 50, 50, 77, 91, 95, 77},
			{479, "Rotom (Heat Rotom)", 50, 65, 107, 86, 105, 107},
			{479, "Rotom (Wash Rotom)", 50, 65, 107, 86, 105, 107},
			{479, "Rotom (Frost Rotom)", 50, 65, 107, 86, 105, 107},
			{479, "Rotom (Fan Rotom)", 50, 65, 107, 86, 105, 107},
			{479, "Rotom (Mow Rotom)", 50, 65, 107, 86, 105, 107}},
		{{480, "Uxie", 75, 75, 130, 95, 75, 130}},
		{{481, "Mesprit", 80, 105, 105, 80, 105, 105}},
		{{482, "Azelf", 75, 125, 70, 115, 125, 70}},
		{{483, "Dialga", 100, 120, 120, 90, 150, 100}},
		{{484, "Palkia", 90, 120, 100, 100, 150, 120}},
		{{485, "Heatran", 91, 90, 106, 77, 130, 106}},
		{{486, "Regigigas", 110, 160, 110, 100, 80, 110}},
		{{487, "Giratina (Altered Forme)", 150, 100, 120, 90, 100, 120},
			{487, "Giratina (Origin Forme)", 150, 120, 100, 90, 120, 100}},
		{{488, "Cresselia", 120, 70, 120, 85, 75, 130}},
		{{489, "Phione", 80, 80, 80, 80, 80, 80}},
		{{490, "Manaphy", 100, 100, 100, 100, 100, 100}},
		{{491, "Darkrai", 70, 90, 90, 125, 135, 90}},
		{{492, "Shaymin (Land Forme)", 100, 100, 100, 100, 100, 100},
			{492, "Shaymin (Sky Forme)", 100, 103, 75, 127, 120, 75}},
		{{493, "Arceus", 120, 120, 120, 120, 120, 120}},
		{{494, "Victini", 100, 100, 100, 100, 100, 100}},
		{{495, "Snivy", 45, 45, 55, 63, 45, 55}},
		{{496, "Servine", 60, 60, 75, 83, 60, 75}},
		{{497, "Serperior", 75, 75, 95, 113, 75, 95}},
		{{498, "Tepig", 65, 63, 45, 45, 45, 45}},
		{{499, "Pignite", 90, 93, 55, 55, 70, 55}},
		{{500, "Emboar", 110, 123, 65, 65, 100, 65}},
		{{501, "Oshawott", 55, 55, 45, 45, 63, 45}},
		{{502, "Dewott", 75, 75, 60, 60, 83, 60}},
		{{503, "Samurott", 95, 100, 85, 70, 108, 70}},
		{{504, "Patrat", 45, 55, 39, 42, 35, 39}},
		{{505, "Watchog", 60, 85, 69, 77, 60, 69}},
		{{506, "Lillipup", 45, 60, 45, 55, 25, 45}},
		{{507, "Herdier", 65, 80, 65, 60, 35, 65}},
		{{508, "Stoutland", 85, 100, 90, 80, 45, 90}},
		{{509, "Purrloin", 41, 50, 37, 66, 50, 37}},
		{{510, "Liepard", 64, 88, 50, 106, 88, 50}},
		{{511, "Pansage", 50, 53, 48, 64, 53, 48}},
		{{512, "Simisage", 75, 98, 63, 101, 98, 63}},
		{{513, "Pansear", 50, 53, 48, 64, 53, 48}},
		{{514, "Simisear", 75, 98, 63, 101, 98, 63}},
		{{515, "Panpour", 50, 53, 48, 64, 53, 48}},
		{{516, "Simipour", 75, 98, 63, 101, 98, 63}},
		{{517, "Munna", 76, 25, 45, 24, 67, 55}},
		{{518, "Musharna", 116, 55, 85, 29, 107, 95}},
		{{519, "Pidove", 50, 55, 50, 43, 36, 30}},
		{{520, "Tranquill", 62, 77, 62, 65, 50, 42}},
		{{521, "Unfezant", 80, 105, 80, 93, 65, 55}},
		{{522, "Blitzle", 45, 60, 32, 76, 50, 32}},
		{{523, "Zebstrika", 75, 100, 63, 116, 80, 63}},
		{{524, "Roggenrola", 55, 75, 85, 15, 25, 25}},
		{{525, "Boldore", 70, 105, 105, 20, 50, 40}},
		{{526, "Gigalith", 85, 135, 130, 25, 60, 70}},
		{{527, "Woobat", 55, 45, 43, 72, 55, 43}},
		{{528, "Swoobat", 67, 57, 55, 114, 77, 55}},
		{{529, "Drilbur", 60, 85, 40, 68, 30, 45}},
		{{530, "Excadrill", 110, 135, 60, 88, 50, 65}},
		{{531, "Audino", 103, 60, 86, 50, 60, 86}},
		{{532, "Timburr", 75, 80, 55, 35, 25, 35}},
		{{533, "Gurdurr", 85, 105, 85, 40, 40, 50}},
		{{534, "Conkeldurr", 105, 140, 95, 45, 55, 65}},
		{{535, "Tympole", 50, 50, 40, 64, 50, 40}},
		{{536, "Palpitoad", 75, 65, 55, 69, 65, 55}},
		{{537, "Seismitoad", 105, 85, 75, 74, 85, 75}},
		{{538, "Throh", 120, 100, 85, 45, 30, 85}},
		{{539, "Sawk", 75, 125, 75, 85, 30, 75}},
		{{540, "Sewaddle", 45, 53, 70, 42, 40, 60}},
		{{541, "Swadloon", 55, 63, 90, 42, 50, 80}},
		{{542, "Leavanny", 75, 103, 80, 92, 70, 70}},
		{{543, "Venipede", 30, 45, 59, 57, 30, 39}},
		{{544, "Whirlipede", 40, 55, 99, 47, 40, 79}},
		{{545, "Scolipede", 60, 90, 89, 112, 55, 69}},
		{{546, "Cottonee", 40, 27, 60, 66, 37, 50}},
		{{547, "Whimsicott", 60, 67, 85, 116, 77, 75}},
		{{548, "Petilil", 45, 35, 50, 30, 70, 50}},
		{{549, "Lilligant", 70, 60, 75, 90, 110, 75}},
		{{550, "Basculin", 70, 92, 65, 98, 80, 55}},
		{{551, "Sandile", 50, 72, 35, 65, 35, 35}},
		{{552, "Krokorok", 60, 82, 45, 74, 45, 45}},
		{{553, "Krookodile", 95, 117, 70, 92, 65, 70}},
		{{554, "Darumaka", 70, 90, 45, 50, 15, 45}},
		{{555, "Darmanitan (Standard Mode)", 105, 140, 55, 95, 30, 55},
			{555, "Darmanitan (Zen Mode)", 105, 30, 105, 55, 140, 105}},
		{{556, "Maractus", 75, 86, 67, 60, 106, 67}},
		{{557, "Dwebble", 50, 65, 85, 55, 35, 35}},
		{{558, "Crustle", 70, 95, 125, 45, 65, 75}},
		{{559, "Scraggy", 50, 75, 70, 48, 35, 70}},
		{{560, "Scrafty", 65, 90, 115, 58, 45, 115}},
		{{561, "Sigilyph", 72, 58, 80, 97, 103, 80}},
		{{562, "Yamask", 38, 30, 85, 30, 55, 65}},
		{{563, "Cofagrigus", 58, 50, 145, 30, 95, 105}},
		{{564, "Tirtouga", 54, 78, 103, 22, 53, 45}},
		{{565, "Carracosta", 74, 108, 133, 32, 83, 65}},
		{{566, "Archen", 55, 112, 45, 70, 74, 45}},
		{{567, "Archeops", 75, 140, 65, 110, 112, 65}},
		{{568, "Trubbish", 50, 50, 62, 65, 40, 62}},
		{{569, "Garbodor", 80, 95, 82, 75, 60, 82}},
		{{570, "Zorua", 40, 65, 40, 65, 80, 40}},
		{{571, "Zoroark", 60, 105, 60, 105, 120, 60}},
		{{572, "Minccino", 55, 50, 40, 75, 40, 40}},
		{{573, "Cinccino", 75, 95, 60, 115, 65, 60}},
		{{574, "Gothita", 45, 30, 50, 45, 55, 65}},
		{{575, "Gothorita", 60, 45, 70, 55, 75, 85}},
		{{576, "Gothitelle", 70, 55, 95, 65, 95, 110}},
		{{577, "Solosis", 45, 30, 40, 20, 105, 50}},
		{{578, "Duosion", 65, 40, 50, 30, 125, 60}},
		{{579, "Reuniclus", 110, 65, 75, 30, 125, 85}},
		{{580, "Ducklett", 62, 44, 50, 55, 44, 50}},
		{{581, "Swanna", 75, 87, 63, 98, 87, 63}},
		{{582, "Vanillite", 36, 50, 50, 44, 65, 60}},
		{{583, "Vanillish", 51, 65, 65, 59, 80, 75}},
		{{584, "Vanilluxe", 71, 95, 85, 79, 110, 95}},
		{{585, "Deerling", 60, 60, 50, 75, 40, 50}},
		{{586, "Sawsbuck", 80, 100, 70, 95, 60, 70}},
		{{587, "Emolga", 55, 75, 60, 103, 75, 60}},
		{{588, "Karrablast", 50, 75, 45, 60, 40, 45}},
		{{589, "Escavalier", 70, 135, 105, 20, 60, 105}},
		{{590, "Foongus", 69, 55, 45, 15, 55, 55}},
		{{591, "Amoonguss", 114, 85, 70, 30, 85, 80}},
		{{592, "Frillish", 55, 40, 50, 40, 65, 85}},
		{{593, "Jellicent", 100, 60, 70, 60, 85, 105}},
		{{594, "Alomomola", 165, 75, 80, 65, 40, 45}},
		{{595, "Joltik", 50, 47, 50, 65, 57, 50}},
		{{596, "Galvantula", 70, 77, 60, 108, 97, 60}},
		{{597, "Ferroseed", 44, 50, 91, 10, 24, 86}},
		{{598, "Ferrothorn", 74, 94, 131, 20, 54, 116}},
		{{599, "Klink", 40, 55, 70, 30, 45, 60}},
		{{600, "Klang", 60, 80, 95, 50, 70, 85}},
		{{601, "Klinklang", 60, 100, 115, 90, 70, 85}},
		{{602, "Tynamo", 35, 55, 40, 60, 45, 40}},
		{{603, "Eelektrik", 65, 85, 70, 40, 75, 70}},
		{{604, "Eelektross", 85, 115, 80, 50, 105, 80}},
		{{605, "Elgyem", 55, 55, 55, 30, 85, 55}},
		{{606, "Beheeyem", 75, 75, 75, 40, 125, 95}},
		{{607, "Litwick", 50, 30, 55, 20, 65, 55}},
		{{608, "Lampent", 60, 40, 60, 55, 95, 60}},
		{{609, "Chandelure", 60, 55, 90, 80, 145, 90}},
		{{610, "Axew", 46, 87, 60, 57, 30, 40}},
		{{611, "Fraxure", 66, 117, 70, 67, 40, 50}},
		{{612, "Haxorus", 76, 147, 90, 97, 60, 70}},
		{{613, "Cubchoo", 55, 70, 40, 40, 60, 40}},
		{{614, "Beartic", 95, 110, 80, 50, 70, 80}},
		{{615, "Cryogonal", 70, 50, 30, 105, 95, 135}},
		{{616, "Shelmet", 50, 40, 85, 25, 40, 65}},
		{{617, "Accelgor", 80, 70, 40, 145, 100, 60}},
		{{618, "Stunfisk", 109, 66, 84, 32, 81, 99}},
		{{619, "Mienfoo", 45, 85, 50, 65, 55, 50}},
		{{620, "Mienshao", 65, 125, 60, 105, 95, 60}},
		{{621, "Druddigon", 77, 120, 90, 48, 60, 90}},
		{{622, "Golett", 59, 74, 50, 35, 35, 50}},
		{{623, "Golurk", 89, 124, 80, 55, 55, 80}},
		{{624, "Pawniard", 45, 85, 70, 60, 40, 40}},
		{{625, "Bisharp", 65, 125, 100, 70, 60, 70}},
		{{626, "Bouffalant", 95, 110, 95, 55, 40, 95}},
		{{627, "Rufflet", 70, 83, 50, 60, 37, 50}},
		{{628, "Braviary", 100, 123, 75, 80, 57, 75}},
		{{629, "Vullaby", 70, 55, 75, 60, 45, 65}},
		{{630, "Mandibuzz", 110, 65, 105, 80, 55, 95}},
		{{631, "Heatmor", 85, 97, 66, 65, 105, 66}},
		{{632, "Durant", 58, 109, 112, 109, 48, 48}},
		{{633, "Deino", 52, 65, 50, 38, 45, 50}},
		{{634, "Zweilous", 72, 85, 70, 58, 65, 70}},
		{{635, "Hydreigon", 92, 105, 90, 98, 125, 90}},
		{{636, "Larvesta", 55, 85, 55, 60, 50, 55}},
		{{637, "Volcarona", 85, 60, 65, 100, 135, 105}},
		{{638, "Cobalion", 91, 90, 129, 108, 90, 72}},
		{{639, "Terrakion", 91, 129, 90, 108, 72, 90}},
		{{640, "Virizion", 91, 90, 72, 108, 90, 129}},
		{{641, "Tornadus (Incarnate Forme)", 79, 115, 70, 111, 125, 80},
			{641, "Tornadus (Therian Forme)", 79, 100, 80, 121, 110, 90}},
		{{642, "Thundurus (Incarnate Forme)", 79, 115, 70, 111, 125, 80},
			{642, "Thundurus (Therian Forme)", 79, 105, 70, 101, 145, 80}},
		{{643, "Reshiram", 100, 120, 100, 90, 150, 120}},
		{{644, "Zekrom", 100, 150, 120, 90, 120, 100}},
		{{645, "Landorus (Incarnate Forme)", 89, 125, 90, 101, 115, 80},
			{645, "Landorus (Therian Forme)", 89, 145, 90, 91, 105, 80}},
		{{646, "Kyurem (Normal Kyurem)", 125, 130, 90, 95, 130, 90},
			{646, "Kyurem (Black Kyurem)", 125, 170, 100, 95, 120, 90},
			{646, "Kyurem (White Kyurem)", 125, 120, 90, 95, 170, 100}},
		{{647, "Keldeo", 91, 72, 90, 108, 129, 90}},
		{{648, "Meloetta (Aria Forme)", 100, 77, 77, 90, 128, 128},
			{648, "Meloetta (Pirouette Forme)", 100, 128, 90, 128, 77, 77}},
		{{649, "Genesect", 71, 120, 95, 99, 120, 95}},
	}
)
