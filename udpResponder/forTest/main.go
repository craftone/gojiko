package main

import (
	"encoding/binary"
	"log"
	"net"
	"os"
	"strconv"
)

var Port uint64

func usage() {
	log.Fatal("Usage: go run forTest/*.go SrcAddr DestAddr PacketSize SeqNum")
}

func main() {
	if len(os.Args) != 5 {
		usage()
	}
	srcAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	destAddr, err := net.ResolveUDPAddr("udp", os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	pktSize, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}
	seqNum, err := strconv.Atoi(os.Args[4])
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", srcAddr, destAddr)
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, pktSize-28)
	binary.BigEndian.PutUint16(buf[0:2], uint16(pktSize))
	binary.BigEndian.PutUint64(buf[2:10], uint64(seqNum))
	_, err = conn.Write(buf)
	if err != nil {
		log.Fatal(err)
	}
}
