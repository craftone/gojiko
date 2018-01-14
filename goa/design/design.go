package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("gojiko api", func() {
	Title("Gojiko API")
	Description(`Gojiko は日本のL2接続MVNOに適した簡素な疑似SGWシミュレータです。
Jmeter等で操作可能なため試験自動化に適しています。`)
	Scheme("http")
	Host("localhost:8080")
})

var _ = Resource("gtpsession", func() {
	BasePath("/gtpsessions")
	DefaultMedia(GtpSessionMedia)

	Action("create", func() {
		Description("Create a new gtp sesseion")
		Routing(POST(""))
		Payload(func() {
			Member("sgwAddr", String, "SGW GTPv2-C loopback address", func() {
				Format("ipv4")
				Default("127.0.0.1")
				Example("127.0.0.1")
			})

			apnMccMncMember()
			msisdnMeiMember()
			imsiEbiMember()
			Required("sgwAddr", "apn", "mcc", "mnc", "msisdn", "mei", "imsi", "ebi")
		})
		Response(OK)
		Response(NotFound)
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
	Description("A gtp session")
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
