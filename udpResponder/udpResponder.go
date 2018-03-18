package main

import (
	"encoding/binary"
	"log"
	"net"
	"os"
	"time"

	"github.com/urfave/cli"
)

var Port uint64
var DebugMode bool
var QueueSize uint64

type RecvPacket struct {
	raddr          *net.UDPAddr
	sendPacketSize uint16
	seqNum         uint64
}

func main() {
	app := cli.NewApp()
	app.Name = "UDP Responder for gojiko"
	app.Usage = "Run at the host who have the UDP Flow's destination IP"
	app.HideVersion = true
	app.Compiled = time.Now()
	app.Flags = []cli.Flag{
		cli.Uint64Flag{
			Name:        "port, p",
			Usage:       "listen port",
			Value:       7777,
			Destination: &Port,
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "show debug log",
		},
		cli.Uint64Flag{
			Name:        "queue, q",
			Usage:       "Queue size",
			Value:       1000,
			Destination: &QueueSize,
		},
	}
	app.Action = udpResponder

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func udpResponder(c *cli.Context) error {
	if Port > 65535 {
		return cli.NewExitError("Port number should be less than 65536.", 1)
	}
	if c.Bool("debug") {
		DebugMode = true
	}

	udpAddr := &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: int(Port),
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return cli.NewExitError(err, 2)
	}

	buf := make([]byte, 65536)
	if DebugMode {
		log.Print("[[DEBUG MODE]]")
	}

	// Start sender goroutine
	toSender := make(chan RecvPacket, QueueSize)
	go sender(udpConn, toSender)
	go sender(udpConn, toSender)
	go sender(udpConn, toSender)
	go sender(udpConn, toSender)
	go sender(udpConn, toSender)
	go statReporter(5)

	log.Printf("Starting UDP responder [ %s ] ...", udpAddr.String())
	for {
		n, addr, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			return cli.NewExitError(err, 3)
		}
		writeRecv(addr.String(), 1, uint64(n))
		if n < 10 {
			log.Printf("Received a invalid packet from %s : %#v", addr.String(), buf[:n])
			continue
		}
		sendPacketSize := binary.BigEndian.Uint16(buf[0:2])
		seqNum := binary.BigEndian.Uint64(buf[2:10])
		if DebugMode {
			log.Printf("Received a packet from %s, size: %d, response size: %d, sequence number: %d",
				addr.String(), n, sendPacketSize, seqNum)
			log.Printf("body: %#v", buf[:n])
		}
		toSender <- RecvPacket{
			raddr:          addr,
			sendPacketSize: sendPacketSize,
			seqNum:         seqNum,
		}
	}

	return nil
}

func sender(udpConn *net.UDPConn, toSender chan RecvPacket) {
	if DebugMode {
		log.Print("Start Sender ...")
	}
	buf := make([]byte, 65536)

	for recv := range toSender {
		binary.BigEndian.PutUint16(buf[0:2], recv.sendPacketSize)
		binary.BigEndian.PutUint64(buf[2:10], recv.seqNum)
		size := recv.sendPacketSize - 28 // 20:IP header, 8: UDP header
		udpConn.WriteTo(buf[0:size], recv.raddr)
		if DebugMode {
			log.Printf("Send a packet raddr: %s, len: %d, seqNum: %d",
				recv.raddr.String(), recv.sendPacketSize, recv.seqNum)
		}
		writeSend(recv.raddr.String(), 1, uint64(recv.sendPacketSize))
	}
}

func statReporter(interval uint) {
	t := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		select {
		case <-t.C:
			reports := theAddrSendRecvStats.Strings()
			for _, report := range reports {
				log.Print(report)
			}
		}
	}

}
