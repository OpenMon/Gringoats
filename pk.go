package main

type PK interface {
	Gen() uint8
	Id() uint16
	Form() uint8
	Nickname() string
	Level() uint8
	Held() uint16
	IsAsleep() bool
	IsPoisoned() bool
	IsBurned() bool
	IsFrozen() bool
	IsParalyzed() bool
	EV() []uint16
	IV() []uint16
	Stats() []uint16
	Moves() []uint16
	PP() []uint8
	IsShiny() bool
	IsMissingNo() bool
	Bytes() []byte
}

type PKData struct {
	Gen       uint8        `json:"generation"`
	Id        uint16       `json:"id"`
	Form      uint8        `json:"form"`
	Shiny     bool         `json:"shiny"`
	Level     uint8        `json:"level"`
	MissingNo bool         `json:"missing-no"`
	Held      uint16       `json:"held_item"`
	Moves     []PKDataMove `json:"moves"`
	EV        []uint16     `json:"ev"`
	IV        []uint16     `json:"iv"`
	Stats     []uint16     `json:"stats"`
}

type PKDataMove struct {
	Id uint16 `json:"id"`
	PP uint8  `json:"pp"`
}

func GetPKData(pk PK) *PKData {
	if pk == nil {
		return nil
	}
	pkpp := pk.PP()
	pkmv := pk.Moves()

	pkd := PKData{
		pk.Gen(),
		pk.Id(),
		pk.Form(),
		pk.IsShiny(),
		pk.Level(),
		pk.IsMissingNo(),
		pk.Held(),
		[]PKDataMove{
			{pkmv[0], pkpp[0]},
			{pkmv[1], pkpp[1]},
			{pkmv[2], pkpp[2]},
			{pkmv[3], pkpp[3]},
		},
		pk.EV(),
		pk.IV(),
		pk.Stats(),
	}

	return &pkd
}

type PKDowngradeCompatible interface {
	Downgrade() (PK, error)
}

type PKUpgradeCompatible interface {
	Upgrade() (PK, error)
}
