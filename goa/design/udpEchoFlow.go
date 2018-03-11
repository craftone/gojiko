package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("udpEchoFlowByIMSIandEBI", func() {
	BasePath("/sgw/:sgwAddr/gtpsessions/imsi/:imsi/ebi/:ebi/udp_echo_flow")
	Params(func() {
		sgwAddrMember()
		imsiEbiMember()
	})

	Action("create", func() {
		Description("Start UDP ECHO flow by IMSI and EBI")
		Routing(POST(""))
		Payload(UdpEchoFlowPayload)

		Response(OK, func() {
			Media(UdpEchoFlowMedia)
		})
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("show", func() {
		Description("Show UDP ECHO flow by IMSI and EBI. The flow is Current flow or last processed flow.")
		Routing(GET(""))
		Response(OK, func() {
			Media(UdpEchoFlowMedia, "withStats")
		})
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("delete", func() {
		Description("End UDP ECHO flow by IMSI and EBI")
		Routing(DELETE(""))
		Response(OK, func() {
			Media(UdpEchoFlowMedia, "withStats")
		})
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})
})

var UdpEchoFlowPayload = Type("UdpEchoFlowPayload", func() {
	Member("destAddr", String, "ECHO destination IPv4 address", func() {
		Format("ipv4")
	})
	Member("destPort", Integer, "ECHO destination UDP port", func() {
		Minimum(0)
		Maximum(65535)
		Default(7777)
		Example(7777)
	})
	Member("sourcePort", Integer, "ECHO source UDP port", func() {
		Minimum(0)
		Maximum(65535)
		Default(7777)
		Example(7777)
	})
	Member("sendPacketSize", Integer, "Send packet size (including IP header)", func() {
		Minimum(38)
		Maximum(1460)
		Example(1460)
	})
	Member("tos", Integer, "Type of service", func() {
		Minimum(0)
		Maximum(255)
		Default(0)
		Example(0)
	})
	Member("ttl", Integer, "Time To Live", func() {
		Minimum(0)
		Maximum(255)
		Default(255)
		Example(255)
	})
	Member("targetBps", Integer, "Target bitrate(bps) in SGi not S5/S8", func() {
		Minimum(1)
		Maximum(100000000000)
		Example(100000000)
	})
	Member("numOfSend", Integer, "Number of send packets", func() {
		Minimum(1)
	})
	Member("recvPacketSize", Integer, "Receive packet size (including IP header)", func() {
		Minimum(38)
		Maximum(1460)
		Example(1460)
	})
	Required("destAddr", "sendPacketSize", "targetBps", "numOfSend", "recvPacketSize")
})

var UdpEchoFlowMedia = MediaType("application/vnd.udpechoflow+json", func() {
	Description("A UDP ECHO flow")
	Attribute("param", UdpEchoFlowPayload)
	Attribute("stats", SendRecvStats)
	View("default", func() {
		Attribute("param")
	})
	View("withStats", func() {
		Attribute("param")
		Attribute("stats")
	})
})
