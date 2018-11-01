package main

import (
	"encoding/json"
	"gopkg.in/antage/eventsource.v1"
	"time"
)

type Event interface {
	Data() string
	Type() string
}

type SlotUpdateEvent struct {
	Location CompleteLocation `json:"location"`
	PKData   *PKData          `json:"data"`
}

func (e SlotUpdateEvent) Data() string {
	d, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(d)
}
func (e SlotUpdateEvent) Type() string {
	return "slot_update"
}

type BoxRenameEvent struct {
	Location CompleteLocation `json:"location"`
	Name     string           `json:"name"`
}

func (e BoxRenameEvent) Data() string {
	d, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(d)
}
func (e BoxRenameEvent) Type() string {
	return "box_rename"
}

type PCRenameEvent struct {
	Location CompleteLocation `json:"location"`
	Name     string           `json:"name"`
}

func (e PCRenameEvent) Data() string {
	d, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(d)
}
func (e PCRenameEvent) Type() string {
	return "pc_rename"
}

type BoxCreateEvent struct {
	Location CompleteLocation `json:"location"`
}

func (e BoxCreateEvent) Data() string {
	d, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(d)
}
func (e BoxCreateEvent) Type() string {
	return "box_create"
}

type ComputerCreateEvent struct {
	Id  uint8 `json:"id"`
	Gen uint8 `json:"generation"`
}

func (e ComputerCreateEvent) Data() string {
	d, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(d)
}
func (e ComputerCreateEvent) Type() string {
	return "pc_create"
}

var (
	eventSource = eventsource.New(nil, nil)
)

func emit(event Event) {
	eventSource.SendEventMessage(event.Data(), event.Type(), time.Now().String())
}
