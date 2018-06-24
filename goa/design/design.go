package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("gojiko api", func() {
	Title("Gojiko API")
	Description(`Gojiko は日本のL2接続MVNOでの利用に適した簡素な疑似SGWシミュレータです。`)
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
