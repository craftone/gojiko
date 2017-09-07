package ie

import (
	"log"
)

func getBit(b byte, idx uint) bool {
	if idx < 0 || idx > 7 {
		log.Fatal("the idx must be from 0 to 7")
	}
	if (b>>idx)&1 == 0 {
		return false
	}
	return true
}

func setBit(b byte, idx uint, flag bool) byte {
	if idx < 0 || idx > 7 {
		log.Fatal("the idx must be from 0 to 7")
	}
	if flag {
		return b | (1 << idx)
	}
	return b &^ (1 << idx)
}
