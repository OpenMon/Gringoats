package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"io/ioutil"
	"net/http"
	"strconv"
)

func FindComputer(w *http.ResponseWriter, ps *httprouter.Params) *Computer {
	id, err := strconv.Atoi(ps.ByName("pcId"))
	if err != nil {
		http.Error(*w, "not a number", http.StatusBadRequest)
		log.Error(err.Error())
		return nil
	}

	if len(BankData.Computers) < id {
		http.Error(*w, "pc not found", http.StatusNotFound)
		log.Error(err.Error())
		return nil
	}
	return &BankData.Computers[id]
}

func HandleComputers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make([]*ComputerData, len(BankData.Computers))
	for i, c := range BankData.Computers {
		data[i] = c.GetComputerData()
	}

	s, err := json.Marshal(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	fmt.Fprint(w, string(s))
}

func HandleComputer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := FindComputer(&w, &ps)
	if c == nil {
		return
	}
	data := c.GetComputerData()

	s, err := json.Marshal(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	fmt.Fprint(w, string(s))
}

func HandleBoxes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := FindComputer(&w, &ps)
	if c == nil {
		return
	}
	data := make([]*BoxData, len(c.Boxes))

	for i, br := range c.Boxes {
		data[i] = br.GetBoxData()
	}

	s, err := json.Marshal(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	fmt.Fprint(w, string(s))
}

func HandleBox(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := FindComputer(&w, &ps)
	if c == nil {
		return
	}

	boxId, err := strconv.Atoi(ps.ByName("boxId"))
	if err != nil {
		http.Error(w, "not a number", http.StatusBadRequest)
		log.Error(err.Error())
		return
	}
	if len(c.Boxes) < boxId {
		http.Error(w, "box not found", http.StatusNotFound)
		return
	}
	data := c.Boxes[boxId].GetBoxData()

	s, err := json.Marshal(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	fmt.Fprint(w, string(s))
}

func HandlePokemonMove(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := FindComputer(&w, &ps)
	if c == nil {
		return
	}

	boxFromId, err := strconv.Atoi(ps.ByName("boxFromId"))
	if err != nil {
		http.Error(w, "not a number", http.StatusBadRequest)
		log.Error(err.Error())
		return
	}
	if len(c.Boxes) < boxFromId {
		http.Error(w, "box not found", http.StatusNotFound)
		return
	}
	slotFromId, err := strconv.Atoi(ps.ByName("slotFromId"))
	if err != nil {
		http.Error(w, "not a number", http.StatusBadRequest)
		log.Error(err.Error())
		return
	}
	if slotFromId >= 30 {
		http.Error(w, "slot not found", http.StatusNotFound)
		return
	}

	boxToId, err := strconv.Atoi(ps.ByName("boxToId"))
	if err != nil {
		http.Error(w, "not a number", http.StatusBadRequest)
		log.Error(err.Error())
		return
	}
	if len(c.Boxes) < boxToId {
		http.Error(w, "box not found", http.StatusNotFound)
		return
	}
	slotToId, err := strconv.Atoi(ps.ByName("slotToId"))
	if err != nil {
		http.Error(w, "not a number", http.StatusBadRequest)
		log.Error(err.Error())
		return
	}
	if slotToId >= 30 {
		http.Error(w, "slot not found", http.StatusNotFound)
		return
	}

	tpk := c.Boxes[boxToId].GetBox().Content[slotToId].Pokemon
	c.SetPK(&PCLocation{boxToId, slotToId}, c.Boxes[boxFromId].GetBox().Content[slotFromId].Pokemon)
	c.SetPK(&PCLocation{boxFromId, slotFromId}, tpk)
}

func HandleBoxRename(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := FindComputer(&w, &ps)
	if c == nil {
		return
	}

	boxId, err := strconv.Atoi(ps.ByName("boxId"))
	if err != nil {
		http.Error(w, "not a number", http.StatusBadRequest)
		log.Error(err.Error())
		return
	}
	if len(c.Boxes) < boxId {
		http.Error(w, "box not found", http.StatusNotFound)
		return
	}
	name := ps.ByName("name")
	if len(name) > 10 {
		http.Error(w, "name too long", http.StatusBadRequest)
		return
	}

	c.RenameBox(boxId, name)
}

func HandleComputerRename(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := FindComputer(&w, &ps)
	if c == nil {
		return
	}

	name := ps.ByName("name")
	if len(name) > 10 {
		http.Error(w, "name too long", http.StatusBadRequest)
		return
	}

	c.RenamePC(name)
}

func HandleSlotDump(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := FindComputer(&w, &ps)
	if c == nil {
		return
	}

	boxId, err := strconv.Atoi(ps.ByName("boxId"))
	if err != nil {
		http.Error(w, "not a number", http.StatusBadRequest)
		log.Error(err.Error())
		return
	}
	if len(c.Boxes) < boxId {
		http.Error(w, "box not found", http.StatusNotFound)
		return
	}
	slotId, err := strconv.Atoi(ps.ByName("slotId"))
	if err != nil {
		http.Error(w, "not a number", http.StatusBadRequest)
		log.Error(err.Error())
		return
	}
	if slotId >= 30 {
		http.Error(w, "slot not found", http.StatusNotFound)
		return
	}

	slot := &c.Boxes[boxId].GetBox().Content[slotId]
	if slot.Pokemon == nil {
		http.Error(w, "empty slot", http.StatusNotFound)
		return
	}
	w.Write(slot.Pokemon.Bytes())
}

func HandleComputerSend(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !IsValidSavSize(r.ContentLength) {
		http.Error(w, "invalid sav size", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	savVersion := DetectSavVersion(data)
	if savVersion == SavInvalid {
		http.Error(w, "invalid sav", http.StatusBadRequest)
		return
	}

	if savVersion == SavGen1 {
		//TODO CreateComputerGen1(data)
	}
}

func HandleSendG1PK1(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := FindComputer(&w, &ps)
	if c == nil {
		return
	}
	if c.Gen != 1 {
		http.Error(w, "wrong pc generation", http.StatusBadRequest)
		return
	}

	pkd := make([]byte, 69)
	ln, err := r.Body.Read(pkd)
	if err != nil && ln != len(pkd) {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}
	if ln != len(pkd) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	pk := newPK1G(pkd[3:])
	loc := c.AddPokemon(pk)

	s, err := json.Marshal(loc)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	fmt.Fprint(w, string(s))
}

func HandleSendG1Raw(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := FindComputer(&w, &ps)
	if c == nil {
		return
	}
	if c.Gen != 1 {
		http.Error(w, "wrong pc generation", http.StatusBadRequest)
		return
	}

	pkd := make([]byte, 33)
	ln, err := r.Body.Read(pkd)
	if err != nil && ln != len(pkd) {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}
	if ln != len(pkd) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	pk := newPK1G(pkd)
	loc := c.AddPokemon(pk)

	s, err := json.Marshal(loc)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	fmt.Fprint(w, string(s))
}

func HandleSendG1PHBankGB(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := FindComputer(&w, &ps)
	if c == nil {
		return
	}
	if c.Gen != 1 {
		http.Error(w, "wrong pc generation", http.StatusBadRequest)
		return
	}

	bank := make([]byte, 36352)
	ln, err := r.Body.Read(bank)
	if err != nil && ln != len(bank) {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	if !bytes.Equal(bank[:4], []byte{'B', 'K', 'G', 'B'}) {
		http.Error(w, "bad magic", http.StatusBadRequest)
		return
	}

	version := bank[4] | bank[5]<<8 | bank[6]<<16 | bank[7]<<24
	if version != 160 {
		http.Error(w, "incompatible version", http.StatusBadRequest)
		return
	}

	adds := make([]*PCLocation, 0)
	offsetBox1 := 0x100
	boxSize := 2 + (32)*(1+(0x21)+(0xB)*2)
	for i := 0; i < 20; i++ {
		boxOffset := offsetBox1 + i*boxSize
		for j := 0; j < 32; j++ {
			pkmOffset := boxOffset + 33*(j+1) + 1

			if bank[pkmOffset] != 0x00 && bank[pkmOffset] != 0xFF {
				adds = append(adds, c.AddPokemon(newPK1G(bank[pkmOffset:pkmOffset+33])))
			}
		}
	}

	s, err := json.Marshal(adds)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	fmt.Fprint(w, string(s))
}

func HandleSendG2PK2(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := FindComputer(&w, &ps)
	if c == nil {
		return
	}
	if c.Gen != 2 {
		http.Error(w, "wrong pc generation", http.StatusBadRequest)
		return
	}

	pkd := make([]byte, 73)
	ln, err := r.Body.Read(pkd)
	if err != nil && ln != len(pkd) {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}
	if ln != len(pkd) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	pk := newPK2G(pkd[3:])
	loc := c.AddPokemon(pk)

	s, err := json.Marshal(loc)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	fmt.Fprint(w, string(s))
}

func HandleSendG2Raw(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := FindComputer(&w, &ps)
	if c == nil {
		return
	}
	if c.Gen != 2 {
		http.Error(w, "wrong pc generation", http.StatusBadRequest)
		return
	}

	pkd := make([]byte, 32)
	ln, err := r.Body.Read(pkd)
	if err != nil && ln != len(pkd) {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}
	if ln != len(pkd) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	pk := newPK2G(pkd)
	loc := c.AddPokemon(pk)

	s, err := json.Marshal(loc)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	fmt.Fprint(w, string(s))
}

func HandleSendG3Raw(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := FindComputer(&w, &ps)
	if c == nil {
		return
	}
	if c.Gen != 3 {
		http.Error(w, "wrong pc generation", http.StatusBadRequest)
		return
	}

	pkd := make([]byte, 80)
	ln, err := r.Body.Read(pkd)
	if err != nil && ln != len(pkd) {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}
	if ln != len(pkd) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	pk := newPK3G(pkd)
	fmt.Print(pk)
	loc := c.AddPokemon(pk)

	s, err := json.Marshal(loc)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	fmt.Fprint(w, string(s))
}

func InitHTTPServer() {
	router := httprouter.New()
	//router.GET("/", Index)
	router.GET("/pc", HandleComputers)
	router.GET("/pc/:pcId", HandleComputer)
	router.GET("/pc/:pcId/rename/:name", HandleComputerRename)
	router.GET("/pc/:pcId/box", HandleBoxes)
	router.GET("/pc/:pcId/box/:boxId", HandleBox)
	router.GET("/pc/:pcId/box/:boxId/rename/:name", HandleBoxRename)
	router.GET("/pc/:pcId/move/:boxFromId/:slotFromId/:boxToId/:slotToId", HandlePokemonMove)
	router.GET("/pc/:pcId/box/:boxId/dump/:slotId", HandleSlotDump)
	//TODO historicize removals
	//router.GET("/pc/:pcId/box/:boxId/poll/:slotId", HandleSlotPoll)
	router.POST("/sav/send", HandleComputerSend)
	router.POST("/pc/:pcId/send/1g/pk1", HandleSendG1PK1)
	router.POST("/pc/:pcId/send/1g/raw", HandleSendG1Raw)
	router.POST("/pc/:pcId/send/1g/phbankgb", HandleSendG1PHBankGB)
	router.POST("/pc/:pcId/send/2g/pk2", HandleSendG2PK2)
	router.POST("/pc/:pcId/send/2g/raw", HandleSendG2Raw)
	//router.POST("/pc/:pcId/send/3g/pkm", HandleSendG3Pkm)
	router.POST("/pc/:pcId/send/3g/raw", HandleSendG3Raw)

	defer eventSource.Close()
	router.Handler("GET", "/event", eventSource)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
	})
	http.ListenAndServe(":8080", c.Handler(router))
}
