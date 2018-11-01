package main

import (
	"bytes"
	"testing"
)

func TestPK1Serialization(t *testing.T) {
	pk := PK1G{}

	b1 := pk.Bytes()
	b2 := newPK1G(b1).Bytes()

	diff := bytes.Compare(b1, b2)

	if diff != 0 {
		t.Log("PK1 are not the same after serialization/deserialization")
		t.Fail()
	}
}
