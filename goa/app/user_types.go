// Code generated by goagen v1.3.0, DO NOT EDIT.
//
// API "gojiko api": Application User Types
//
// Command:
// $ goagen
// --design=github.com/craftone/gojiko/goa/design
// --out=$(GOPATH)/src/github.com/craftone/gojiko/goa
// --version=v1.3.0

package app

import (
	"github.com/goadesign/goa"
	"time"
)

// recvStatForHuman user type.
type recvStatForHuman struct {
	// Received bitrate
	Bitrate *string `form:"bitrate,omitempty" json:"bitrate,omitempty" xml:"bitrate,omitempty"`
	// Received bytes
	Bytes *string `form:"bytes,omitempty" json:"bytes,omitempty" xml:"bytes,omitempty"`
	// Received invalid bytes
	InvalidBytes *string `form:"invalidBytes,omitempty" json:"invalidBytes,omitempty" xml:"invalidBytes,omitempty"`
	// Received invalid packets
	InvalidPackets *string `form:"invalidPackets,omitempty" json:"invalidPackets,omitempty" xml:"invalidPackets,omitempty"`
	// Received packets
	Packets *string `form:"packets,omitempty" json:"packets,omitempty" xml:"packets,omitempty"`
}

// Validate validates the recvStatForHuman type instance.
func (ut *recvStatForHuman) Validate() (err error) {
	if ut.Bitrate == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "bitrate"))
	}
	if ut.Bytes == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "bytes"))
	}
	if ut.Packets == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "packets"))
	}
	if ut.InvalidBytes == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "invalidBytes"))
	}
	if ut.InvalidPackets == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "invalidPackets"))
	}
	return
}

// Publicize creates RecvStatForHuman from recvStatForHuman
func (ut *recvStatForHuman) Publicize() *RecvStatForHuman {
	var pub RecvStatForHuman
	if ut.Bitrate != nil {
		pub.Bitrate = *ut.Bitrate
	}
	if ut.Bytes != nil {
		pub.Bytes = *ut.Bytes
	}
	if ut.InvalidBytes != nil {
		pub.InvalidBytes = *ut.InvalidBytes
	}
	if ut.InvalidPackets != nil {
		pub.InvalidPackets = *ut.InvalidPackets
	}
	if ut.Packets != nil {
		pub.Packets = *ut.Packets
	}
	return &pub
}

// RecvStatForHuman user type.
type RecvStatForHuman struct {
	// Received bitrate
	Bitrate string `form:"bitrate" json:"bitrate" xml:"bitrate"`
	// Received bytes
	Bytes string `form:"bytes" json:"bytes" xml:"bytes"`
	// Received invalid bytes
	InvalidBytes string `form:"invalidBytes" json:"invalidBytes" xml:"invalidBytes"`
	// Received invalid packets
	InvalidPackets string `form:"invalidPackets" json:"invalidPackets" xml:"invalidPackets"`
	// Received packets
	Packets string `form:"packets" json:"packets" xml:"packets"`
}

// Validate validates the RecvStatForHuman type instance.
func (ut *RecvStatForHuman) Validate() (err error) {
	if ut.Bitrate == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "bitrate"))
	}
	if ut.Bytes == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "bytes"))
	}
	if ut.Packets == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "packets"))
	}
	if ut.InvalidBytes == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "invalidBytes"))
	}
	if ut.InvalidPackets == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "invalidPackets"))
	}
	return
}

// recvStatForMachie user type.
type recvStatForMachie struct {
	// Received bitrate
	Bitrate *float64 `form:"bitrate,omitempty" json:"bitrate,omitempty" xml:"bitrate,omitempty"`
	// Received bytes
	Bytes *int `form:"bytes,omitempty" json:"bytes,omitempty" xml:"bytes,omitempty"`
	// Received invalid bytes
	InvalidBytes *int `form:"invalidBytes,omitempty" json:"invalidBytes,omitempty" xml:"invalidBytes,omitempty"`
	// Received invalid packets
	InvalidPackets *int `form:"invalidPackets,omitempty" json:"invalidPackets,omitempty" xml:"invalidPackets,omitempty"`
	// Received packets
	Packets *int `form:"packets,omitempty" json:"packets,omitempty" xml:"packets,omitempty"`
}

// Validate validates the recvStatForMachie type instance.
func (ut *recvStatForMachie) Validate() (err error) {
	if ut.Bitrate == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "bitrate"))
	}
	if ut.Bytes == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "bytes"))
	}
	if ut.Packets == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "packets"))
	}
	if ut.InvalidBytes == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "invalidBytes"))
	}
	if ut.InvalidPackets == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "invalidPackets"))
	}
	return
}

// Publicize creates RecvStatForMachie from recvStatForMachie
func (ut *recvStatForMachie) Publicize() *RecvStatForMachie {
	var pub RecvStatForMachie
	if ut.Bitrate != nil {
		pub.Bitrate = *ut.Bitrate
	}
	if ut.Bytes != nil {
		pub.Bytes = *ut.Bytes
	}
	if ut.InvalidBytes != nil {
		pub.InvalidBytes = *ut.InvalidBytes
	}
	if ut.InvalidPackets != nil {
		pub.InvalidPackets = *ut.InvalidPackets
	}
	if ut.Packets != nil {
		pub.Packets = *ut.Packets
	}
	return &pub
}

// RecvStatForMachie user type.
type RecvStatForMachie struct {
	// Received bitrate
	Bitrate float64 `form:"bitrate" json:"bitrate" xml:"bitrate"`
	// Received bytes
	Bytes int `form:"bytes" json:"bytes" xml:"bytes"`
	// Received invalid bytes
	InvalidBytes int `form:"invalidBytes" json:"invalidBytes" xml:"invalidBytes"`
	// Received invalid packets
	InvalidPackets int `form:"invalidPackets" json:"invalidPackets" xml:"invalidPackets"`
	// Received packets
	Packets int `form:"packets" json:"packets" xml:"packets"`
}

// Validate validates the RecvStatForMachie type instance.
func (ut *RecvStatForMachie) Validate() (err error) {

	return
}

// sendRecvStatistics user type.
type sendRecvStatistics struct {
	// seconds
	Duration *float64   `form:"duration,omitempty" json:"duration,omitempty" xml:"duration,omitempty"`
	EndTime  *time.Time `form:"endTime,omitempty" json:"endTime,omitempty" xml:"endTime,omitempty"`
	// receive statistics for machines
	RecvStats *recvStatForMachie `form:"recvStats,omitempty" json:"recvStats,omitempty" xml:"recvStats,omitempty"`
	// receive statistics for humans
	RecvStatsHumanize *recvStatForHuman `form:"recvStatsHumanize,omitempty" json:"recvStatsHumanize,omitempty" xml:"recvStatsHumanize,omitempty"`
	// send statistics for machines
	SendStats *sendStatForMachie `form:"sendStats,omitempty" json:"sendStats,omitempty" xml:"sendStats,omitempty"`
	// send statistics for humans
	SendStatsHumanize *sendStatForHuman `form:"sendStatsHumanize,omitempty" json:"sendStatsHumanize,omitempty" xml:"sendStatsHumanize,omitempty"`
	StartTime         *time.Time        `form:"startTime,omitempty" json:"startTime,omitempty" xml:"startTime,omitempty"`
}

// Validate validates the sendRecvStatistics type instance.
func (ut *sendRecvStatistics) Validate() (err error) {
	if ut.StartTime == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "startTime"))
	}
	if ut.Duration == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "duration"))
	}
	if ut.SendStats == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "sendStats"))
	}
	if ut.SendStatsHumanize == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "sendStatsHumanize"))
	}
	if ut.RecvStats == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "recvStats"))
	}
	if ut.RecvStatsHumanize == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "recvStatsHumanize"))
	}
	if ut.RecvStats != nil {
		if err2 := ut.RecvStats.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if ut.RecvStatsHumanize != nil {
		if err2 := ut.RecvStatsHumanize.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if ut.SendStats != nil {
		if err2 := ut.SendStats.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if ut.SendStatsHumanize != nil {
		if err2 := ut.SendStatsHumanize.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// Publicize creates SendRecvStatistics from sendRecvStatistics
func (ut *sendRecvStatistics) Publicize() *SendRecvStatistics {
	var pub SendRecvStatistics
	if ut.Duration != nil {
		pub.Duration = *ut.Duration
	}
	if ut.EndTime != nil {
		pub.EndTime = ut.EndTime
	}
	if ut.RecvStats != nil {
		pub.RecvStats = ut.RecvStats.Publicize()
	}
	if ut.RecvStatsHumanize != nil {
		pub.RecvStatsHumanize = ut.RecvStatsHumanize.Publicize()
	}
	if ut.SendStats != nil {
		pub.SendStats = ut.SendStats.Publicize()
	}
	if ut.SendStatsHumanize != nil {
		pub.SendStatsHumanize = ut.SendStatsHumanize.Publicize()
	}
	if ut.StartTime != nil {
		pub.StartTime = *ut.StartTime
	}
	return &pub
}

// SendRecvStatistics user type.
type SendRecvStatistics struct {
	// seconds
	Duration float64    `form:"duration" json:"duration" xml:"duration"`
	EndTime  *time.Time `form:"endTime,omitempty" json:"endTime,omitempty" xml:"endTime,omitempty"`
	// receive statistics for machines
	RecvStats *RecvStatForMachie `form:"recvStats" json:"recvStats" xml:"recvStats"`
	// receive statistics for humans
	RecvStatsHumanize *RecvStatForHuman `form:"recvStatsHumanize" json:"recvStatsHumanize" xml:"recvStatsHumanize"`
	// send statistics for machines
	SendStats *SendStatForMachie `form:"sendStats" json:"sendStats" xml:"sendStats"`
	// send statistics for humans
	SendStatsHumanize *SendStatForHuman `form:"sendStatsHumanize" json:"sendStatsHumanize" xml:"sendStatsHumanize"`
	StartTime         time.Time         `form:"startTime" json:"startTime" xml:"startTime"`
}

// Validate validates the SendRecvStatistics type instance.
func (ut *SendRecvStatistics) Validate() (err error) {

	if ut.SendStats == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "sendStats"))
	}
	if ut.SendStatsHumanize == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "sendStatsHumanize"))
	}
	if ut.RecvStats == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "recvStats"))
	}
	if ut.RecvStatsHumanize == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "recvStatsHumanize"))
	}
	if ut.RecvStatsHumanize != nil {
		if err2 := ut.RecvStatsHumanize.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if ut.SendStatsHumanize != nil {
		if err2 := ut.SendStatsHumanize.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// sendStatForHuman user type.
type sendStatForHuman struct {
	// Send bitrate
	Bitrate *string `form:"bitrate,omitempty" json:"bitrate,omitempty" xml:"bitrate,omitempty"`
	// Send bytes
	Bytes *string `form:"bytes,omitempty" json:"bytes,omitempty" xml:"bytes,omitempty"`
	// Send packets
	Packets *string `form:"packets,omitempty" json:"packets,omitempty" xml:"packets,omitempty"`
	// Skipped bytes
	SkippedBytes *string `form:"skippedBytes,omitempty" json:"skippedBytes,omitempty" xml:"skippedBytes,omitempty"`
	// Skipped packets
	SkippedPackets *string `form:"skippedPackets,omitempty" json:"skippedPackets,omitempty" xml:"skippedPackets,omitempty"`
}

// Validate validates the sendStatForHuman type instance.
func (ut *sendStatForHuman) Validate() (err error) {
	if ut.Bitrate == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "bitrate"))
	}
	if ut.Bytes == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "bytes"))
	}
	if ut.Packets == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "packets"))
	}
	if ut.SkippedBytes == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "skippedBytes"))
	}
	if ut.SkippedPackets == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "skippedPackets"))
	}
	return
}

// Publicize creates SendStatForHuman from sendStatForHuman
func (ut *sendStatForHuman) Publicize() *SendStatForHuman {
	var pub SendStatForHuman
	if ut.Bitrate != nil {
		pub.Bitrate = *ut.Bitrate
	}
	if ut.Bytes != nil {
		pub.Bytes = *ut.Bytes
	}
	if ut.Packets != nil {
		pub.Packets = *ut.Packets
	}
	if ut.SkippedBytes != nil {
		pub.SkippedBytes = *ut.SkippedBytes
	}
	if ut.SkippedPackets != nil {
		pub.SkippedPackets = *ut.SkippedPackets
	}
	return &pub
}

// SendStatForHuman user type.
type SendStatForHuman struct {
	// Send bitrate
	Bitrate string `form:"bitrate" json:"bitrate" xml:"bitrate"`
	// Send bytes
	Bytes string `form:"bytes" json:"bytes" xml:"bytes"`
	// Send packets
	Packets string `form:"packets" json:"packets" xml:"packets"`
	// Skipped bytes
	SkippedBytes string `form:"skippedBytes" json:"skippedBytes" xml:"skippedBytes"`
	// Skipped packets
	SkippedPackets string `form:"skippedPackets" json:"skippedPackets" xml:"skippedPackets"`
}

// Validate validates the SendStatForHuman type instance.
func (ut *SendStatForHuman) Validate() (err error) {
	if ut.Bitrate == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "bitrate"))
	}
	if ut.Bytes == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "bytes"))
	}
	if ut.Packets == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "packets"))
	}
	if ut.SkippedBytes == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "skippedBytes"))
	}
	if ut.SkippedPackets == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "skippedPackets"))
	}
	return
}

// sendStatForMachie user type.
type sendStatForMachie struct {
	// Send bitrate
	Bitrate *float64 `form:"bitrate,omitempty" json:"bitrate,omitempty" xml:"bitrate,omitempty"`
	// Send bytes
	Bytes *int `form:"bytes,omitempty" json:"bytes,omitempty" xml:"bytes,omitempty"`
	// Send packets
	Packets *int `form:"packets,omitempty" json:"packets,omitempty" xml:"packets,omitempty"`
	// Skipped bytes
	SkippedBytes *int `form:"skippedBytes,omitempty" json:"skippedBytes,omitempty" xml:"skippedBytes,omitempty"`
	// Skipped packets
	SkippedPackets *int `form:"skippedPackets,omitempty" json:"skippedPackets,omitempty" xml:"skippedPackets,omitempty"`
}

// Validate validates the sendStatForMachie type instance.
func (ut *sendStatForMachie) Validate() (err error) {
	if ut.Bitrate == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "bitrate"))
	}
	if ut.Bytes == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "bytes"))
	}
	if ut.Packets == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "packets"))
	}
	if ut.SkippedBytes == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "skippedBytes"))
	}
	if ut.SkippedPackets == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "skippedPackets"))
	}
	return
}

// Publicize creates SendStatForMachie from sendStatForMachie
func (ut *sendStatForMachie) Publicize() *SendStatForMachie {
	var pub SendStatForMachie
	if ut.Bitrate != nil {
		pub.Bitrate = *ut.Bitrate
	}
	if ut.Bytes != nil {
		pub.Bytes = *ut.Bytes
	}
	if ut.Packets != nil {
		pub.Packets = *ut.Packets
	}
	if ut.SkippedBytes != nil {
		pub.SkippedBytes = *ut.SkippedBytes
	}
	if ut.SkippedPackets != nil {
		pub.SkippedPackets = *ut.SkippedPackets
	}
	return &pub
}

// SendStatForMachie user type.
type SendStatForMachie struct {
	// Send bitrate
	Bitrate float64 `form:"bitrate" json:"bitrate" xml:"bitrate"`
	// Send bytes
	Bytes int `form:"bytes" json:"bytes" xml:"bytes"`
	// Send packets
	Packets int `form:"packets" json:"packets" xml:"packets"`
	// Skipped bytes
	SkippedBytes int `form:"skippedBytes" json:"skippedBytes" xml:"skippedBytes"`
	// Skipped packets
	SkippedPackets int `form:"skippedPackets" json:"skippedPackets" xml:"skippedPackets"`
}

// Validate validates the SendStatForMachie type instance.
func (ut *SendStatForMachie) Validate() (err error) {

	return
}

// udpEchoFlowPayload user type.
type udpEchoFlowPayload struct {
	// ECHO destination IPv4 address
	DestAddr *string `form:"destAddr,omitempty" json:"destAddr,omitempty" xml:"destAddr,omitempty"`
	// ECHO destination UDP port
	DestPort *int `form:"destPort,omitempty" json:"destPort,omitempty" xml:"destPort,omitempty"`
	// Number of send packets
	NumOfSend *int `form:"numOfSend,omitempty" json:"numOfSend,omitempty" xml:"numOfSend,omitempty"`
	// Receive packet size (including IP header)
	RecvPacketSize *int `form:"recvPacketSize,omitempty" json:"recvPacketSize,omitempty" xml:"recvPacketSize,omitempty"`
	// Send packet size (including IP header)
	SendPacketSize *int `form:"sendPacketSize,omitempty" json:"sendPacketSize,omitempty" xml:"sendPacketSize,omitempty"`
	// ECHO source UDP port
	SourcePort *int `form:"sourcePort,omitempty" json:"sourcePort,omitempty" xml:"sourcePort,omitempty"`
	// Target bitrate(bps) in SGi not S5/S8
	TargetBps *int `form:"targetBps,omitempty" json:"targetBps,omitempty" xml:"targetBps,omitempty"`
	// Type of service
	Tos *int `form:"tos,omitempty" json:"tos,omitempty" xml:"tos,omitempty"`
	// Time To Live
	TTL *int `form:"ttl,omitempty" json:"ttl,omitempty" xml:"ttl,omitempty"`
}

// Finalize sets the default values for udpEchoFlowPayload type instance.
func (ut *udpEchoFlowPayload) Finalize() {
	var defaultDestPort = 7777
	if ut.DestPort == nil {
		ut.DestPort = &defaultDestPort
	}
	var defaultSourcePort = 7777
	if ut.SourcePort == nil {
		ut.SourcePort = &defaultSourcePort
	}
	var defaultTos = 0
	if ut.Tos == nil {
		ut.Tos = &defaultTos
	}
	var defaultTTL = 255
	if ut.TTL == nil {
		ut.TTL = &defaultTTL
	}
}

// Validate validates the udpEchoFlowPayload type instance.
func (ut *udpEchoFlowPayload) Validate() (err error) {
	if ut.DestAddr == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "destAddr"))
	}
	if ut.SendPacketSize == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "sendPacketSize"))
	}
	if ut.TargetBps == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "targetBps"))
	}
	if ut.NumOfSend == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "numOfSend"))
	}
	if ut.RecvPacketSize == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "recvPacketSize"))
	}
	if ut.DestAddr != nil {
		if err2 := goa.ValidateFormat(goa.FormatIPv4, *ut.DestAddr); err2 != nil {
			err = goa.MergeErrors(err, goa.InvalidFormatError(`request.destAddr`, *ut.DestAddr, goa.FormatIPv4, err2))
		}
	}
	if ut.DestPort != nil {
		if *ut.DestPort < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.destPort`, *ut.DestPort, 0, true))
		}
	}
	if ut.DestPort != nil {
		if *ut.DestPort > 65535 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.destPort`, *ut.DestPort, 65535, false))
		}
	}
	if ut.NumOfSend != nil {
		if *ut.NumOfSend < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.numOfSend`, *ut.NumOfSend, 1, true))
		}
	}
	if ut.RecvPacketSize != nil {
		if *ut.RecvPacketSize < 38 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.recvPacketSize`, *ut.RecvPacketSize, 38, true))
		}
	}
	if ut.RecvPacketSize != nil {
		if *ut.RecvPacketSize > 1460 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.recvPacketSize`, *ut.RecvPacketSize, 1460, false))
		}
	}
	if ut.SendPacketSize != nil {
		if *ut.SendPacketSize < 38 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.sendPacketSize`, *ut.SendPacketSize, 38, true))
		}
	}
	if ut.SendPacketSize != nil {
		if *ut.SendPacketSize > 1460 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.sendPacketSize`, *ut.SendPacketSize, 1460, false))
		}
	}
	if ut.SourcePort != nil {
		if *ut.SourcePort < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.sourcePort`, *ut.SourcePort, 0, true))
		}
	}
	if ut.SourcePort != nil {
		if *ut.SourcePort > 65535 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.sourcePort`, *ut.SourcePort, 65535, false))
		}
	}
	if ut.TargetBps != nil {
		if *ut.TargetBps < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.targetBps`, *ut.TargetBps, 1, true))
		}
	}
	if ut.TargetBps != nil {
		if *ut.TargetBps > 100000000000 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.targetBps`, *ut.TargetBps, 100000000000, false))
		}
	}
	if ut.Tos != nil {
		if *ut.Tos < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.tos`, *ut.Tos, 0, true))
		}
	}
	if ut.Tos != nil {
		if *ut.Tos > 255 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.tos`, *ut.Tos, 255, false))
		}
	}
	if ut.TTL != nil {
		if *ut.TTL < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.ttl`, *ut.TTL, 0, true))
		}
	}
	if ut.TTL != nil {
		if *ut.TTL > 255 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.ttl`, *ut.TTL, 255, false))
		}
	}
	return
}

// Publicize creates UDPEchoFlowPayload from udpEchoFlowPayload
func (ut *udpEchoFlowPayload) Publicize() *UDPEchoFlowPayload {
	var pub UDPEchoFlowPayload
	if ut.DestAddr != nil {
		pub.DestAddr = *ut.DestAddr
	}
	if ut.DestPort != nil {
		pub.DestPort = *ut.DestPort
	}
	if ut.NumOfSend != nil {
		pub.NumOfSend = *ut.NumOfSend
	}
	if ut.RecvPacketSize != nil {
		pub.RecvPacketSize = *ut.RecvPacketSize
	}
	if ut.SendPacketSize != nil {
		pub.SendPacketSize = *ut.SendPacketSize
	}
	if ut.SourcePort != nil {
		pub.SourcePort = *ut.SourcePort
	}
	if ut.TargetBps != nil {
		pub.TargetBps = *ut.TargetBps
	}
	if ut.Tos != nil {
		pub.Tos = *ut.Tos
	}
	if ut.TTL != nil {
		pub.TTL = *ut.TTL
	}
	return &pub
}

// UDPEchoFlowPayload user type.
type UDPEchoFlowPayload struct {
	// ECHO destination IPv4 address
	DestAddr string `form:"destAddr" json:"destAddr" xml:"destAddr"`
	// ECHO destination UDP port
	DestPort int `form:"destPort" json:"destPort" xml:"destPort"`
	// Number of send packets
	NumOfSend int `form:"numOfSend" json:"numOfSend" xml:"numOfSend"`
	// Receive packet size (including IP header)
	RecvPacketSize int `form:"recvPacketSize" json:"recvPacketSize" xml:"recvPacketSize"`
	// Send packet size (including IP header)
	SendPacketSize int `form:"sendPacketSize" json:"sendPacketSize" xml:"sendPacketSize"`
	// ECHO source UDP port
	SourcePort int `form:"sourcePort" json:"sourcePort" xml:"sourcePort"`
	// Target bitrate(bps) in SGi not S5/S8
	TargetBps int `form:"targetBps" json:"targetBps" xml:"targetBps"`
	// Type of service
	Tos int `form:"tos" json:"tos" xml:"tos"`
	// Time To Live
	TTL int `form:"ttl" json:"ttl" xml:"ttl"`
}

// Validate validates the UDPEchoFlowPayload type instance.
func (ut *UDPEchoFlowPayload) Validate() (err error) {
	if ut.DestAddr == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "destAddr"))
	}

	if err2 := goa.ValidateFormat(goa.FormatIPv4, ut.DestAddr); err2 != nil {
		err = goa.MergeErrors(err, goa.InvalidFormatError(`type.destAddr`, ut.DestAddr, goa.FormatIPv4, err2))
	}
	if ut.DestPort < 0 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.destPort`, ut.DestPort, 0, true))
	}
	if ut.DestPort > 65535 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.destPort`, ut.DestPort, 65535, false))
	}
	if ut.NumOfSend < 1 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.numOfSend`, ut.NumOfSend, 1, true))
	}
	if ut.RecvPacketSize < 38 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.recvPacketSize`, ut.RecvPacketSize, 38, true))
	}
	if ut.RecvPacketSize > 1460 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.recvPacketSize`, ut.RecvPacketSize, 1460, false))
	}
	if ut.SendPacketSize < 38 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.sendPacketSize`, ut.SendPacketSize, 38, true))
	}
	if ut.SendPacketSize > 1460 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.sendPacketSize`, ut.SendPacketSize, 1460, false))
	}
	if ut.SourcePort < 0 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.sourcePort`, ut.SourcePort, 0, true))
	}
	if ut.SourcePort > 65535 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.sourcePort`, ut.SourcePort, 65535, false))
	}
	if ut.TargetBps < 1 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.targetBps`, ut.TargetBps, 1, true))
	}
	if ut.TargetBps > 100000000000 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.targetBps`, ut.TargetBps, 100000000000, false))
	}
	if ut.Tos < 0 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.tos`, ut.Tos, 0, true))
	}
	if ut.Tos > 255 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.tos`, ut.Tos, 255, false))
	}
	if ut.TTL < 0 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.ttl`, ut.TTL, 0, true))
	}
	if ut.TTL > 255 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.ttl`, ut.TTL, 255, false))
	}
	return
}

// ecgi user type.
type ecgi struct {
	// E-UTRAN Cell Identifier
	Eci *int `form:"eci,omitempty" json:"eci,omitempty" xml:"eci,omitempty"`
	// Mobile Country Code
	Mcc *string `form:"mcc,omitempty" json:"mcc,omitempty" xml:"mcc,omitempty"`
	// Mobile Network Code
	Mnc *string `form:"mnc,omitempty" json:"mnc,omitempty" xml:"mnc,omitempty"`
}

// Finalize sets the default values for ecgi type instance.
func (ut *ecgi) Finalize() {
	var defaultEci = 1
	if ut.Eci == nil {
		ut.Eci = &defaultEci
	}
	var defaultMcc = "440"
	if ut.Mcc == nil {
		ut.Mcc = &defaultMcc
	}
	var defaultMnc = "10"
	if ut.Mnc == nil {
		ut.Mnc = &defaultMnc
	}
}

// Validate validates the ecgi type instance.
func (ut *ecgi) Validate() (err error) {
	if ut.Eci != nil {
		if *ut.Eci < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.eci`, *ut.Eci, 0, true))
		}
	}
	if ut.Eci != nil {
		if *ut.Eci > 268435455 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.eci`, *ut.Eci, 268435455, false))
		}
	}
	if ut.Mcc != nil {
		if ok := goa.ValidatePattern(`^[0-9]{3}$`, *ut.Mcc); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`request.mcc`, *ut.Mcc, `^[0-9]{3}$`))
		}
	}
	if ut.Mnc != nil {
		if ok := goa.ValidatePattern(`^[0-9]{2,3}$`, *ut.Mnc); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`request.mnc`, *ut.Mnc, `^[0-9]{2,3}$`))
		}
	}
	return
}

// Publicize creates Ecgi from ecgi
func (ut *ecgi) Publicize() *Ecgi {
	var pub Ecgi
	if ut.Eci != nil {
		pub.Eci = *ut.Eci
	}
	if ut.Mcc != nil {
		pub.Mcc = *ut.Mcc
	}
	if ut.Mnc != nil {
		pub.Mnc = *ut.Mnc
	}
	return &pub
}

// Ecgi user type.
type Ecgi struct {
	// E-UTRAN Cell Identifier
	Eci int `form:"eci" json:"eci" xml:"eci"`
	// Mobile Country Code
	Mcc string `form:"mcc" json:"mcc" xml:"mcc"`
	// Mobile Network Code
	Mnc string `form:"mnc" json:"mnc" xml:"mnc"`
}

// Validate validates the Ecgi type instance.
func (ut *Ecgi) Validate() (err error) {
	if ut.Eci < 0 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.eci`, ut.Eci, 0, true))
	}
	if ut.Eci > 268435455 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.eci`, ut.Eci, 268435455, false))
	}
	if ok := goa.ValidatePattern(`^[0-9]{3}$`, ut.Mcc); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`type.mcc`, ut.Mcc, `^[0-9]{3}$`))
	}
	if ok := goa.ValidatePattern(`^[0-9]{2,3}$`, ut.Mnc); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`type.mnc`, ut.Mnc, `^[0-9]{2,3}$`))
	}
	return
}

// fteid user type.
type fteid struct {
	Ipv4 *string `form:"ipv4,omitempty" json:"ipv4,omitempty" xml:"ipv4,omitempty"`
	Teid *string `form:"teid,omitempty" json:"teid,omitempty" xml:"teid,omitempty"`
}

// Validate validates the fteid type instance.
func (ut *fteid) Validate() (err error) {
	if ut.Teid == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "teid"))
	}
	if ut.Ipv4 == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "ipv4"))
	}
	if ut.Ipv4 != nil {
		if err2 := goa.ValidateFormat(goa.FormatIPv4, *ut.Ipv4); err2 != nil {
			err = goa.MergeErrors(err, goa.InvalidFormatError(`request.ipv4`, *ut.Ipv4, goa.FormatIPv4, err2))
		}
	}
	if ut.Teid != nil {
		if ok := goa.ValidatePattern(`^0x[0-9A-F]{8}$`, *ut.Teid); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`request.teid`, *ut.Teid, `^0x[0-9A-F]{8}$`))
		}
	}
	return
}

// Publicize creates Fteid from fteid
func (ut *fteid) Publicize() *Fteid {
	var pub Fteid
	if ut.Ipv4 != nil {
		pub.Ipv4 = *ut.Ipv4
	}
	if ut.Teid != nil {
		pub.Teid = *ut.Teid
	}
	return &pub
}

// Fteid user type.
type Fteid struct {
	Ipv4 string `form:"ipv4" json:"ipv4" xml:"ipv4"`
	Teid string `form:"teid" json:"teid" xml:"teid"`
}

// Validate validates the Fteid type instance.
func (ut *Fteid) Validate() (err error) {
	if ut.Teid == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "teid"))
	}
	if ut.Ipv4 == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "ipv4"))
	}
	if err2 := goa.ValidateFormat(goa.FormatIPv4, ut.Ipv4); err2 != nil {
		err = goa.MergeErrors(err, goa.InvalidFormatError(`type.ipv4`, ut.Ipv4, goa.FormatIPv4, err2))
	}
	if ok := goa.ValidatePattern(`^0x[0-9A-F]{8}$`, ut.Teid); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`type.teid`, ut.Teid, `^0x[0-9A-F]{8}$`))
	}
	return
}

// gtpSessionFTEIDs user type.
type gtpSessionFTEIDs struct {
	PgwCtrlFTEID *fteid `form:"pgwCtrlFTEID,omitempty" json:"pgwCtrlFTEID,omitempty" xml:"pgwCtrlFTEID,omitempty"`
	PgwDataFTEID *fteid `form:"pgwDataFTEID,omitempty" json:"pgwDataFTEID,omitempty" xml:"pgwDataFTEID,omitempty"`
	SgwCtrlFTEID *fteid `form:"sgwCtrlFTEID,omitempty" json:"sgwCtrlFTEID,omitempty" xml:"sgwCtrlFTEID,omitempty"`
	SgwDataFTEID *fteid `form:"sgwDataFTEID,omitempty" json:"sgwDataFTEID,omitempty" xml:"sgwDataFTEID,omitempty"`
}

// Validate validates the gtpSessionFTEIDs type instance.
func (ut *gtpSessionFTEIDs) Validate() (err error) {
	if ut.PgwCtrlFTEID != nil {
		if err2 := ut.PgwCtrlFTEID.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if ut.PgwDataFTEID != nil {
		if err2 := ut.PgwDataFTEID.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if ut.SgwCtrlFTEID != nil {
		if err2 := ut.SgwCtrlFTEID.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if ut.SgwDataFTEID != nil {
		if err2 := ut.SgwDataFTEID.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// Publicize creates GtpSessionFTEIDs from gtpSessionFTEIDs
func (ut *gtpSessionFTEIDs) Publicize() *GtpSessionFTEIDs {
	var pub GtpSessionFTEIDs
	if ut.PgwCtrlFTEID != nil {
		pub.PgwCtrlFTEID = ut.PgwCtrlFTEID.Publicize()
	}
	if ut.PgwDataFTEID != nil {
		pub.PgwDataFTEID = ut.PgwDataFTEID.Publicize()
	}
	if ut.SgwCtrlFTEID != nil {
		pub.SgwCtrlFTEID = ut.SgwCtrlFTEID.Publicize()
	}
	if ut.SgwDataFTEID != nil {
		pub.SgwDataFTEID = ut.SgwDataFTEID.Publicize()
	}
	return &pub
}

// GtpSessionFTEIDs user type.
type GtpSessionFTEIDs struct {
	PgwCtrlFTEID *Fteid `form:"pgwCtrlFTEID,omitempty" json:"pgwCtrlFTEID,omitempty" xml:"pgwCtrlFTEID,omitempty"`
	PgwDataFTEID *Fteid `form:"pgwDataFTEID,omitempty" json:"pgwDataFTEID,omitempty" xml:"pgwDataFTEID,omitempty"`
	SgwCtrlFTEID *Fteid `form:"sgwCtrlFTEID,omitempty" json:"sgwCtrlFTEID,omitempty" xml:"sgwCtrlFTEID,omitempty"`
	SgwDataFTEID *Fteid `form:"sgwDataFTEID,omitempty" json:"sgwDataFTEID,omitempty" xml:"sgwDataFTEID,omitempty"`
}

// Validate validates the GtpSessionFTEIDs type instance.
func (ut *GtpSessionFTEIDs) Validate() (err error) {
	if ut.PgwCtrlFTEID != nil {
		if err2 := ut.PgwCtrlFTEID.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if ut.PgwDataFTEID != nil {
		if err2 := ut.PgwDataFTEID.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if ut.SgwCtrlFTEID != nil {
		if err2 := ut.SgwCtrlFTEID.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if ut.SgwDataFTEID != nil {
		if err2 := ut.SgwDataFTEID.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ratType user type.
type ratType struct {
	RatType *string `form:"ratType,omitempty" json:"ratType,omitempty" xml:"ratType,omitempty"`
	// Default is 6 : E-UTRAN (WB-E-UTRAN)
	RatTypeValue *int `form:"ratTypeValue,omitempty" json:"ratTypeValue,omitempty" xml:"ratTypeValue,omitempty"`
}

// Finalize sets the default values for ratType type instance.
func (ut *ratType) Finalize() {
	var defaultRatType = "<unkown>"
	if ut.RatType == nil {
		ut.RatType = &defaultRatType
	}
	var defaultRatTypeValue = 6
	if ut.RatTypeValue == nil {
		ut.RatTypeValue = &defaultRatTypeValue
	}
}

// Validate validates the ratType type instance.
func (ut *ratType) Validate() (err error) {
	if ut.RatTypeValue != nil {
		if *ut.RatTypeValue < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.ratTypeValue`, *ut.RatTypeValue, 0, true))
		}
	}
	if ut.RatTypeValue != nil {
		if *ut.RatTypeValue > 255 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.ratTypeValue`, *ut.RatTypeValue, 255, false))
		}
	}
	return
}

// Publicize creates RatType from ratType
func (ut *ratType) Publicize() *RatType {
	var pub RatType
	if ut.RatType != nil {
		pub.RatType = *ut.RatType
	}
	if ut.RatTypeValue != nil {
		pub.RatTypeValue = *ut.RatTypeValue
	}
	return &pub
}

// RatType user type.
type RatType struct {
	RatType string `form:"ratType" json:"ratType" xml:"ratType"`
	// Default is 6 : E-UTRAN (WB-E-UTRAN)
	RatTypeValue int `form:"ratTypeValue" json:"ratTypeValue" xml:"ratTypeValue"`
}

// Validate validates the RatType type instance.
func (ut *RatType) Validate() (err error) {
	if ut.RatTypeValue < 0 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.ratTypeValue`, ut.RatTypeValue, 0, true))
	}
	if ut.RatTypeValue > 255 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.ratTypeValue`, ut.RatTypeValue, 255, false))
	}
	return
}

// tai user type.
type tai struct {
	// Mobile Country Code
	Mcc *string `form:"mcc,omitempty" json:"mcc,omitempty" xml:"mcc,omitempty"`
	// Mobile Network Code
	Mnc *string `form:"mnc,omitempty" json:"mnc,omitempty" xml:"mnc,omitempty"`
	// Tracking Area Code
	Tac *int `form:"tac,omitempty" json:"tac,omitempty" xml:"tac,omitempty"`
}

// Finalize sets the default values for tai type instance.
func (ut *tai) Finalize() {
	var defaultMcc = "440"
	if ut.Mcc == nil {
		ut.Mcc = &defaultMcc
	}
	var defaultMnc = "10"
	if ut.Mnc == nil {
		ut.Mnc = &defaultMnc
	}
	var defaultTac = 1
	if ut.Tac == nil {
		ut.Tac = &defaultTac
	}
}

// Validate validates the tai type instance.
func (ut *tai) Validate() (err error) {
	if ut.Mcc != nil {
		if ok := goa.ValidatePattern(`^[0-9]{3}$`, *ut.Mcc); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`request.mcc`, *ut.Mcc, `^[0-9]{3}$`))
		}
	}
	if ut.Mnc != nil {
		if ok := goa.ValidatePattern(`^[0-9]{2,3}$`, *ut.Mnc); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`request.mnc`, *ut.Mnc, `^[0-9]{2,3}$`))
		}
	}
	if ut.Tac != nil {
		if *ut.Tac < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.tac`, *ut.Tac, 0, true))
		}
	}
	if ut.Tac != nil {
		if *ut.Tac > 65535 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`request.tac`, *ut.Tac, 65535, false))
		}
	}
	return
}

// Publicize creates Tai from tai
func (ut *tai) Publicize() *Tai {
	var pub Tai
	if ut.Mcc != nil {
		pub.Mcc = *ut.Mcc
	}
	if ut.Mnc != nil {
		pub.Mnc = *ut.Mnc
	}
	if ut.Tac != nil {
		pub.Tac = *ut.Tac
	}
	return &pub
}

// Tai user type.
type Tai struct {
	// Mobile Country Code
	Mcc string `form:"mcc" json:"mcc" xml:"mcc"`
	// Mobile Network Code
	Mnc string `form:"mnc" json:"mnc" xml:"mnc"`
	// Tracking Area Code
	Tac int `form:"tac" json:"tac" xml:"tac"`
}

// Validate validates the Tai type instance.
func (ut *Tai) Validate() (err error) {
	if ok := goa.ValidatePattern(`^[0-9]{3}$`, ut.Mcc); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`type.mcc`, ut.Mcc, `^[0-9]{3}$`))
	}
	if ok := goa.ValidatePattern(`^[0-9]{2,3}$`, ut.Mnc); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`type.mnc`, ut.Mnc, `^[0-9]{2,3}$`))
	}
	if ut.Tac < 0 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.tac`, ut.Tac, 0, true))
	}
	if ut.Tac > 65535 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`type.tac`, ut.Tac, 65535, false))
	}
	return
}
