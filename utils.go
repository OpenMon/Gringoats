package main

import (
	"encoding/binary"
)

func UInt16(data []byte, index int) uint16 {
	return binary.BigEndian.Uint16(data[index : index+2])
}

func UInt16ToBytes(data []byte, index int, attr uint16) {
	binary.BigEndian.PutUint16(data[index:index+2], attr)
}

func UInt24(data []byte, index int) uint32 {
	buff := make([]byte, 4)
	buff[3] = 0
	copy(buff, data[index:index+3])
	return binary.BigEndian.Uint32(buff)
}

func UInt24ToBytes(data []byte, index int, attr uint32) {
	buff := make([]byte, 4)
	binary.BigEndian.PutUint32(buff, attr)
	copy(data[index:index+3], buff[:3])
}

func UInt32(data []byte, index int) uint32 {
	return binary.BigEndian.Uint32(data[index : index+4])
}

func UInt32ToBytes(data []byte, index int, attr uint32) {
	binary.BigEndian.PutUint32(data[index:index+4], attr)
}

func LUInt16(data []byte, index int) uint16 {
	return binary.LittleEndian.Uint16(data[index : index+2])
}

func LUInt16ToBytes(data []byte, index int, attr uint16) {
	binary.LittleEndian.PutUint16(data[index:index+2], attr)
}

func LUInt24(data []byte, index int) uint32 {
	buff := make([]byte, 4)
	buff[3] = 0
	copy(buff, data[index:index+3])
	return binary.LittleEndian.Uint32(buff)
}

func LUInt24ToBytes(data []byte, index int, attr uint32) {
	buff := make([]byte, 4)
	binary.LittleEndian.PutUint32(buff, attr)
	copy(data[index:index+3], buff[:3])
}

func LUInt32(data []byte, index int) uint32 {
	return binary.LittleEndian.Uint32(data[index : index+4])
}

func LUInt32ToBytes(data []byte, index int, attr uint32) {
	binary.LittleEndian.PutUint32(data[index:index+4], attr)
}

func Name(data []byte, index int, size int) string {
	return string(data[index : index+size])
}

func CheckFlag(data uint8, flag uint8) bool {
	return data&flag == flag
}
