package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var apnMccMncMember = func() {
	Member("apn", String, "Access Point Name", func() {
		Format("hostname")
		Example("example.com")
	})
	Member("mcc", String, "Mobile Country Code", func() {
		Pattern(`^[0-9]{3}$`)
		Default("440")
		Example("440")
	})
	Member("mnc", String, "Mobile Network Code", func() {
		Pattern(`^[0-9]{2,3}$`)
		Default("10")
		Example("10")
	})
}

var imsiEbiMember = func() {
	Member("imsi", String, "", func() {
		Pattern(`^[0-9]{14,15}$`)
		Example("440100123456780")
	})
	Param("ebi", Integer, "EPS Bearer ID", func() {
		Minimum(5)
		Maximum(15)
		Default(5)
		Example(5)
	})
}

var msisdnMeiMember = func() {
	Member("msisdn", String, "", func() {
		Pattern(`^[0-9]{12,15}$`)
		Example("8101012345678")
	})
	Member("mei", String, "Mobile Equipment Identifier", func() {
		Pattern(`^[0-9]{15,16}$`)
		Example("1212345612345612")
	})
}

var gtpSessionIDMember = func() {
	Member("sid", Integer, "Session ID in this SGW", func() {
		Minimum(0)
		Example(1)
	})
}

var gtpSessionStateMember = func() {
	Member("state", String, "GTP session's state", func() {
		Enum("idle", "sending Create Sessin Request", "connected")
		Example("idle")
	})
}

var sgwAddrMember = func() {
	Member("sgwAddr", String, "SGW GTPv2-C loopback address", func() {
		Format("ipv4")
		Example("127.0.0.1")
	})
}
