package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

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
			gtpSessionIDMember()
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

	Action("deleteByIMSIandEBI", func() {
		Description("Delete the gtp session by IMSI and EBI")
		Routing(DELETE("/imsi/:imsi/ebi/:ebi"))
		Params(func() {
			imsiEbiMember()
		})
		Response(OK, GtpV2CCauseMedia)
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})
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
		gtpSessionStateMember()

		Attribute("fteid", GtpSessionFTEIDs)
		Attribute("paa", String, "PDN Address Allocation", func() {
			Default("0.0.0.0")
			Format("ipv4")
		})

		apnMccMncMember()
		msisdnMeiMember()
		imsiEbiMember()
		Required("apn", "sid", "imsi", "mcc", "mnc", "mei", "mnc", "msisdn")
	})
	View("default", func() {
		Attribute("sid")
		Attribute("apn")
		Attribute("mcc")
		Attribute("mnc")
		Attribute("msisdn")
		Attribute("mei")
		Attribute("imsi")
		Attribute("ebi")
		Attribute("fteid")
		Attribute("paa")
	})
})

var GtpV2CCauseMedia = MediaType("application/vnd.gtpv2c.cause+json", func() {
	Description("GTPv2-C Cause")
	Attributes(func() {
		Attribute("type", String, "Type of return code from PGW", func() {
			Example("OK")
		})
		Attribute("detail", String, "Detail of return code from PGW", func() {
			Example("Request accepted")
		})
		Required("type", "detail")
	})
	View("default", func() {
		Attribute("type")
		Attribute("detail")
	})
})
