package domain

import (
	"net"

	"github.com/sirupsen/logrus"
)

type opSPgw struct {
	parent        SPgwIf
	numOfSessions int
	fromReceiver  chan UDPpacket
}

func newOpSPgw(parent SPgwIf, raddr net.UDPAddr) (*opSPgw, error) {
	log.WithFields(logrus.Fields{
		"SPgwAddr":         parent.UDPAddr(),
		"OppositeSPgwAddr": raddr,
	}).Info("New OpSPgw has created.")

	return &opSPgw{
		parent:       parent,
		fromReceiver: make(chan UDPpacket),
	}, nil
}
