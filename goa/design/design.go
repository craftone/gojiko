package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("gojiko api", func() {
	Title("Gojiko API")
	Description(`Gojiko は日本のL2接続MVNOでの利用に適した簡素な疑似SGWシミュレータです。
Jmeter等で操作可能であるため、試験自動化に適しています。`)
	Scheme("http")
	Host("localhost:8080")
})

var _ = Resource("gtpsession", func() {
	BasePath("/sgw/:sgwAddr/gtpsessions")
	DefaultMedia(GtpSessionMedia)
	Params(func() {
		sgwAddrMember()
	})

	Action("create", func() {
		Description("Create a new gtp sesseion")
		Routing(POST(""))
		Payload(func() {
			apnMccMncMember()
			msisdnMeiMember()
			imsiEbiMember()
			Required("apn", "mcc", "mnc", "msisdn", "mei", "imsi", "ebi")
		})
		Response(OK)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("showByID", func() {
		Description("Show the gtp session by session ID")
		Routing(GET("/id/:sid"))
		Params(func() {
			sessionIdMember()
		})
		Response(OK)
		Response(NotFound, ErrorMedia)
	})

	Action("showByIMSIandEBI", func() {
		Description("Show the gtp session by IMSI and EBI")
		Routing(GET("/imsi/:imsi/ebi/:ebi"))
		Params(func() {
			imsiEbiMember()
		})
		Response(OK)
		Response(NotFound, ErrorMedia)
	})
})

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

		Response(OK, UdpEchoFlowMedia)
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})
})

var FTEID = Type("fteid", func() {
	Attribute("teid", String, "", func() {
		Pattern("^0x[0-9A-F]{8}$")
		Example("0x12345678")
	})
	Attribute("ipv4", String, "", func() {
		Format("ipv4")
		Example("127.0.0.1")
	})
	Required("teid", "ipv4")
})

var GtpSessionFTEIDs = Type("gtpSessionFTEIDs", func() {
	Attribute("sgwCtrlFTEID", FTEID)
	Attribute("sgwDataFTEID", FTEID)
	Attribute("pgwCtrlFTEID", FTEID)
	Attribute("pgwDataFTEID", FTEID)
})

var GtpSessionMedia = MediaType("application/vnd.gtpsession+json", func() {
	Description("A GTP session")
	Attributes(func() {
		gtpSessionIDMember()
		gtpSessionStatusMember()

		Attribute("fteid", GtpSessionFTEIDs)

		apnMccMncMember()
		msisdnMeiMember()
		imsiEbiMember()
		Required("id", "apn", "mcc", "mnc", "msisdn", "mei", "imsi", "ebi", "fteid")
	})
	View("default", func() {
		Attribute("id")
		Attribute("apn")
		Attribute("mcc")
		Attribute("mnc")
		Attribute("msisdn")
		Attribute("mei")
		Attribute("imsi")
		Attribute("ebi")
		Attribute("fteid")
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
	Member("recvPacketSize", Integer, "Receive packet size (including IP header)", func() {
		Minimum(38)
		Maximum(1460)
		Example(1460)
	})
	Required("destAddr", "sendPacketSize", "targetBps", "recvPacketSize")
})

var UdpEchoFlowMedia = MediaType("application/vnd.udpechoflow+json", func() {
	Description("A UDP ECHO flow")
	Attribute("UdpEchoFlowArg", UdpEchoFlowPayload)
	View("default", func() {
		Attribute("UdpEchoFlowArg")
	})
})
