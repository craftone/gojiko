package util

import (
	"log"
)

// GetBit returns b's bit specified by idx (0 means
// least significant bit and 7 means most significant bit).
// if the bit is 0, it returns false. and the bit is 1, returns true.
func GetBit(b byte, idx uint) bool {
	if idx < 0 || idx > 7 {
		log.Fatal("the idx must be from 0 to 7")
	}
	if (b>>idx)&1 == 0 {
		return false
	}
	return true
}

// SetBit sets b's bit specified by idx (0 means
// least significant bit and 7 means most significant bit)
// to 0 (false) or 1 (true) .
func SetBit(b byte, idx uint, flag bool) byte {
	if idx < 0 || idx > 7 {
		log.Fatal("the idx must be from 0 to 7")
	}
	if flag {
		return b | (1 << idx)
	}
	return b &^ (1 << idx)
}
