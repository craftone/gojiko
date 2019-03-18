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
			Member("pseudoSgwDataAddr", String,
				"Specify when using external pseudo SGW-DATA",
				func() {
					Format("ipv4")
					Example("127.0.0.1")
				})
			Member("pseudoSgwDataTEID", Integer,
				`Specify when using external pseudo SGW-DATA which tunnel's TEID has already determined.
If 0 is specified, TEID is generated automatically.
If pseudoSgwDataAddr is not specified, this attribute is ignored.`,
				func() {
					Minimum(0)
					Maximum(0xFFFFFFFF)
					Default(0)
					Example(1)
				})
			ratTypeValueMember()
			Member("tai", TAI)
			Member("ecgi", ECGI)
			Required("apn", "mcc", "mnc", "msisdn", "mei", "imsi", "ebi")
		})
		Response(OK, GtpV2CResponseMedia)
		Response(BadRequest, ErrorMedia)
		Response(NotFound, ErrorMedia)
		Response(Conflict, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
		Response(ServiceUnavailable, GtpV2CCauseMedia)
		Response(GatewayTimeout, ErrorMedia)
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
		Response(Conflict, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("trackingAreaUpdateWithoutSgwRelocation", func() {
		Description("Update the gtp session by IMSI and EBI")
		Routing(PUT("/imsi/:imsi/ebi/:ebi"))
		Params(func() {
			imsiEbiMember()
		})
		Payload(func() {
			Member("tai", TAI)
			Member("ecgi", ECGI)
			ratTypeValueMember()
			Required("tai", "ecgi")
		})
		Response(OK, GtpV2CResponseMedia)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(Conflict, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
		Response(ServiceUnavailable, GtpV2CCauseMedia)
		Response(GatewayTimeout, ErrorMedia)
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
		Attribute("tai", TAI)
		Attribute("ecgi", ECGI)
		Attribute("ratType", RatType)
		Required("apn", "sid", "imsi", "mcc", "mnc", "mei", "mnc", "msisdn", "tai", "ecgi", "ratType")
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
		Attribute("tai")
		Attribute("ecgi")
		Attribute("ratType")
	})
})

var GtpV2CCauseMedia = MediaType("application/vnd.gtpv2c.cause+json", func() {
	Description("GTPv2-C Cause")
	Attributes(func() {
		Attribute("type", String, "Type of return code from PGW", func() {
			Example("OK")
		})
		Attribute("value", Integer, "GTPv2-C response Cause Value", func() {
			Example(16)
		})
		Attribute("detail", String, "Detail of return code from PGW", func() {
			Example("Request accepted")
		})
		Required("type", "value", "detail")
	})
	View("default", func() {
		Attribute("type")
		Attribute("value")
		Attribute("detail")
	})
})

var GtpV2CResponseMedia = MediaType("application/vnd.gtpv2c.csres+json", func() {
	Description("GTPv2-C Crease Session Response")
	Attributes(func() {
		Attribute("cause", GtpV2CCauseMedia)
		Attribute("sessionInfo", GtpSessionMedia)
		Required("cause", "sessionInfo")
	})
	View("default", func() {
		Attribute("cause")
		Attribute("sessionInfo")
	})

})
