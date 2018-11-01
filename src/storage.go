package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"strings"
	"sync"
)

type Slot struct {
	Lock    sync.Mutex
	Pokemon PK
}

func (s *Slot) IsEmpty() bool {
	return s.Pokemon == nil
}

func (c *Computer) SetPK(loc *PCLocation, pk PK) {
	box := c.Boxes[loc.Box].GetBox()
	s := &box.Content[loc.Slot]
	s.Lock.Lock()
	s.Pokemon = pk

	filePath := box.Computer.GetDirectory() + fmt.Sprintf("/%d", loc.Box)
	file, err := os.OpenFile(filePath, os.O_WRONLY, os.ModeExclusive)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	var byt []byte
	if pk != nil {
		byt = pk.Bytes()
	} else {
		byt = make([]byte, c.GetPKSize())
	}
	ln, err := file.WriteAt(byt, int64(100+(loc.Slot*c.GetPKSize())))
	if err != nil || ln != c.GetPKSize() {
		panic(err)
	}

	emit(&SlotUpdateEvent{
		CompleteLocation{int(c.Id), loc.Box, loc.Slot},
		GetPKData(s.Pokemon)})

	s.Lock.Unlock()
}

type BoxData struct {
	Id    uint8       `json:"id"`
	Name  string      `json:"name"`
	Theme uint8       `json:"theme"`
	Slots [30]*PKData `json:"slots"`
}

type Box struct {
	Computer *Computer
	Name     string
	Theme    uint8
	Content  [30]Slot
}

type BoxRef struct {
	Lock     sync.Mutex
	Computer *Computer
	Id       uint8
	Box      *Box
}

func (c *Computer) RenameBox(boxId int, name string) {
	br := c.Boxes[boxId]
	box := br.GetBox()
	br.Lock.Lock()

	filePath := box.Computer.GetDirectory() + fmt.Sprintf("/%d", boxId)
	file, err := os.OpenFile(filePath, os.O_WRONLY, os.ModeExclusive)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	tmpName := make([]byte, 10)
	copy(tmpName, []byte(name))

	ln, err := file.WriteAt(tmpName, 0)
	if err != nil || ln != len(tmpName) {
		panic(err)
	}

	box.Name = name

	emit(&BoxRenameEvent{
		CompleteLocation{int(c.Id), boxId, 0},
		name})

	br.Lock.Unlock()
}

func (br *BoxRef) GetBox() *Box {
	if br.Box != nil {
		return br.Box
	}

	dir := br.Computer.GetDirectory()
	filePath := fmt.Sprintf("%s/%d", dir, br.Id)
	boxData, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	name := string(boxData[:10])
	if name[0] == 0 {
		name = fmt.Sprintf("Box %d", br.Id+1)
	} else {
		name = strings.TrimRight(name, "\x00")
	}

	box := Box{br.Computer, name, 0, [30]Slot{}}
	null := make([]byte, br.Computer.GetPKSize())
	for i := 0; i < 30; i++ {
		offset := 100 + (i * br.Computer.GetPKSize())
		pkd := boxData[offset : offset+br.Computer.GetPKSize()]
		if !bytes.Equal(null, pkd) {
			box.Content[i].Pokemon = br.Computer.ReadPK(pkd)
		} else {
			box.Content[i].Pokemon = nil
		}
	}

	br.Lock.Lock()
	br.Box = &box
	br.Lock.Unlock()
	return br.Box
}

func (br *BoxRef) GetBoxData() *BoxData {
	b := br.GetBox()
	pkdata := [30]*PKData{}
	for j, s := range b.Content {
		if !s.IsEmpty() {
			pkdata[j] = GetPKData(s.Pokemon)
		} else {
			pkdata[j] = nil
		}
	}
	return &BoxData{br.Id, b.Name, b.Theme, pkdata}
}

type ComputerData struct {
	Id    uint8  `json:"id"`
	Gen   uint8  `json:"generation"`
	Name  string `json:"name"`
	Boxes int    `json:"boxes"`
}

type Computer struct {
	Lock  sync.Mutex
	Id    uint8
	Gen   uint8
	Name  string
	Boxes []BoxRef
}

func (c *Computer) GetPKSize() int {
	switch c.Gen {
	case 1:
		return 33
	case 2:
		return 32
	case 3:
		return 80
	}
	return 0
}

func (c *Computer) ReadPK(data []byte) PK {
	switch c.Gen {
	case 1:
		return newPK1G(data)
	case 2:
		return newPK2G(data)
	case 3:
		return newPK3G(data)
	}
	return nil
}

func (c *Computer) GetComputerData() *ComputerData {
	return &ComputerData{
		c.Id, c.Gen, c.Name, len(c.Boxes),
	}
}

func (c *Computer) CreateBox(locked bool) {
	dir := c.GetDirectory()
	id := uint8(len(c.Boxes))

	filePath := fmt.Sprintf("%s/%d", dir, id)
	err := ioutil.WriteFile(filePath, make([]byte, 100+(30*c.GetPKSize())), 0600)
	if err != nil {
		panic(err)
	}

	if !locked {
		c.Lock.Lock()
	}
	c.Boxes = append(c.Boxes, BoxRef{sync.Mutex{}, c, id, nil})

	emit(&BoxCreateEvent{
		CompleteLocation{int(c.Id), len(c.Boxes) - 1, 0},
	})

	if !locked {
		c.Lock.Unlock()
	}
}

func (c *Computer) FirstAvailableSlot(locked bool) *PCLocation {
	for i, boxRef := range c.Boxes {
		box := boxRef.GetBox()
		for j, slot := range box.Content {
			if slot.IsEmpty() {
				return &PCLocation{i, j}
			}
		}
	}
	c.CreateBox(locked)
	return c.FirstAvailableSlot(locked)
}

type Bank struct {
	Computers []Computer
}

type CompleteLocation struct {
	PC   int `json:"pc"`
	Box  int `json:"box"`
	Slot int `json:"slot"`
}

type PCLocation struct {
	Box  int `json:"box"`
	Slot int `json:"slot"`
}

var (
	BankData *Bank
)

func (c *Computer) AddPokemon(pk PK) *PCLocation {
	c.Lock.Lock()
	loc := c.FirstAvailableSlot(true)
	c.SetPK(loc, pk)
	c.Lock.Unlock()
	return loc
}

func GetAppDirectory() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	dir := usr.HomeDir + "/.gpkb"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0700)
	}

	return dir
}

func GetBankDirectory() string {
	return GetAppDirectory() + "/bank"
}

func GetPCSDirectory() string {
	return GetAppDirectory() + "/bank/pc"
}

func (pc *Computer) GetDirectory() string {
	return GetAppDirectory() + fmt.Sprintf("/bank/pc/%d", pc.Id)
}

func (pc *Computer) WriteMetaData() {
	pc.Lock.Lock()
	file := pc.GetDirectory() + "/.meta"
	meta := make([]byte, 100)
	meta[0] = pc.Gen

	ioutil.WriteFile(file, meta, 0700)
	pc.Lock.Unlock()
}

func (c *Computer) RenamePC(name string) {
	c.Lock.Lock()

	filePath := c.GetDirectory() + "/.meta"
	file, err := os.OpenFile(filePath, os.O_WRONLY, os.ModeExclusive)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	tmpName := make([]byte, 10)
	copy(tmpName, []byte(name))

	ln, err := file.WriteAt(tmpName, 1)
	if err != nil || ln != len(tmpName) {
		panic(err)
	}

	c.Name = name

	emit(&PCRenameEvent{
		CompleteLocation{int(c.Id), 0, 0},
		name})

	c.Lock.Unlock()
}

func (pc *Computer) ReadMetaData() {
	file := pc.GetDirectory() + "/.meta"
	meta, err := ioutil.ReadFile(file)
	if err == nil {
		pc.Gen = meta[0]
	}

	name := string(meta[1:11])
	if name[0] == 0 {
		name = fmt.Sprintf("PC %d", pc.Id+1)
	} else {
		name = strings.TrimRight(name, "\x00")
	}
	pc.Name = name
}

func CreateComputer(id uint8, gen uint8) *Computer {
	c := Computer{sync.Mutex{}, id, gen, "", make([]BoxRef, 0)}
	os.Mkdir(c.GetDirectory(), 0700)
	c.WriteMetaData()
	c.CreateBox(false)

	emit(&ComputerCreateEvent{
		id, gen,
	})

	return &c
}

func CreateStorage() {
	os.Mkdir(GetBankDirectory(), 0700)
	os.Mkdir(GetPCSDirectory(), 0700)

	CreateComputer(0, 1)
}

func LoadStorage(path string) {
	BankData = &Bank{}

	pcs, err := ioutil.ReadDir(GetPCSDirectory())
	if err != nil {
		panic(err)
	}

	for _, pcf := range pcs {
		id, err := strconv.Atoi(pcf.Name())
		if err != nil {
			panic(err)
		}
		c := Computer{sync.Mutex{}, uint8(id), 1, "", make([]BoxRef, 0)}
		c.ReadMetaData()

		boxes, err := ioutil.ReadDir(c.GetDirectory())
		if err != nil {
			panic(err)
		}

		for _, boxf := range boxes {
			boxId, err := strconv.Atoi(boxf.Name())
			if err != nil {
				continue
			}
			boxRef := BoxRef{sync.Mutex{}, &c, uint8(boxId), nil}
			c.Boxes = append(c.Boxes, boxRef)
		}

		BankData.Computers = append(BankData.Computers, c)
	}
}

func InitStorage() {
	storage := GetBankDirectory()
	if _, err := os.Stat(storage); os.IsNotExist(err) {
		CreateStorage()
	}
	LoadStorage(storage)
}
