package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var SendStatsForMachie = Type("SendStatForMachie", func() {
	Member("bitrate", Number, "Send bitrate")
	Member("bytes", Integer, "Send bytes")
	Member("packets", Integer, "Send packets")
	Member("skippedBytes", Integer, "Skipped bytes")
	Member("skippedPackets", Integer, "Skipped packets")
	Required("bitrate", "bytes", "packets", "skippedBytes", "skippedPackets")
})

var SendStatsForHuman = Type("SendStatForHuman", func() {
	Member("bitrate", String, "Send bitrate")
	Member("bytes", String, "Send bytes")
	Member("packets", String, "Send packets")
	Member("skippedBytes", String, "Skipped bytes")
	Member("skippedPackets", String, "Skipped packets")
	Required("bitrate", "bytes", "packets", "skippedBytes", "skippedPackets")
})

var RecvStatsForMachie = Type("RecvStatForMachie", func() {
	Member("bitrate", Number, "Received bitrate")
	Member("bytes", Integer, "Received bytes")
	Member("packets", Integer, "Received packets")
	Member("invalidBytes", Integer, "Received invalid bytes")
	Member("invalidPackets", Integer, "Received invalid packets")
	Required("bitrate", "bytes", "packets", "invalidBytes", "invalidPackets")
})

var RecvStatsForHuman = Type("RecvStatForHuman", func() {
	Member("bitrate", String, "Received bitrate")
	Member("bytes", String, "Received bytes")
	Member("packets", String, "Received packets")
	Member("invalidBytes", String, "Received invalid bytes")
	Member("invalidPackets", String, "Received invalid packets")
	Required("bitrate", "bytes", "packets", "invalidBytes", "invalidPackets")
})

var SendRecvStats = Type("SendRecvStatistics", func() {
	Member("startTime", DateTime)
	Member("endTime", DateTime)
	Member("duration", Number, "seconds")
	Member("sendStats", SendStatsForMachie, "send statistics for machines")
	Member("sendStatsHumanize", SendStatsForHuman, "send statistics for humans")
	Member("recvStats", RecvStatsForMachie, "receive statistics for machines")
	Member("recvStatsHumanize", RecvStatsForHuman, "receive statistics for humans")
	Required("startTime", "duration", "sendStats", "sendStatsHumanize", "recvStats", "recvStatsHumanize")
})
