package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("gojiko api", func() {
	Title("Gojiko API")
	Description(`Gojiko is a S5 simulator suitable for 
	using by Japanese MVNO connected via L2 layer`)
	Scheme("http")
	Host("localhost:8080")
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

var TAI = Type("tai", func() {
	mccMncMember()
	Member("tac", Integer, "Tracking Area Code", func() {
		Default(1)
		Minimum(0)
		Maximum(0xFFFF)
		Example(1)
	})
})

var ECGI = Type("ecgi", func() {
	mccMncMember()
	Member("eci", Integer, "E-UTRAN Cell Identifier", func() {
		Default(1)
		Minimum(0)
		Maximum(0x0FFFFFFF) // 28bit
		Example(1)
	})
})

var RatType = Type("ratType", func() {
	ratTypeValueMember()
	Attribute("ratType", String, "", func() {
		Example("EUTRAN (WB-E-UTRAN)")
		Default("<unkown>")
	})
})
