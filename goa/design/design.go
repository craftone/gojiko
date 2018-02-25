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
		Param("sgwAddr", String, "SGW GTPv2-C loopback address", func() {
			Format("ipv4")
			Example("127.0.0.1")
		})
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
			Param("sid", Integer, "Session ID")
		})
		Response(OK)
		Response(NotFound, ErrorMedia)
	})

	Action("showByIMSIandEBI", func() {
		Description("Show the gtp session by IMSI and EBI")
		Routing(GET("/imsi/:imsi/ebi/:ebi"))
		Params(func() {
			Param("imsi")
			Param("ebi")
		})
		Response(OK)
		Response(NotFound, ErrorMedia)
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
