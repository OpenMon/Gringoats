package main

import (
	"bytes"
	"testing"
)

func TestPK2Serialization(t *testing.T) {
	pk := PK2G{}

	b1 := pk.Bytes()
	b2 := newPK2G(b1).Bytes()

	diff := bytes.Compare(b1, b2)

	if diff != 0 {
		t.Log("PK2 are not the same after serialization/deserialization")
		t.Fail()
	}
}
