package ie

import (
	"errors"
	"log"
)

type causeValue byte

const (
	// Request / Initial message

	CauseLocalDetach                          causeValue = 2
	CauseCompleteDetach                       causeValue = 3
	CauseRATChangedFrom3GPPtoNon3GPP          causeValue = 4
	CauseISRDeactivation                      causeValue = 5
	CauseErrorIndicationReceivedFromRNCeNodeB causeValue = 6
	CauseIMSIDetachOnly                       causeValue = 7
	CauseReactivationRequested                causeValue = 8
	CausePDNReconnectionToThisAPNDisallowed   causeValue = 9
	CauseAccessChangedFromNon3GPPto3GPP       causeValue = 10
	CausePDNConnectionInactivityTimerExpires  causeValue = 11

	// Acceptance in a Response / triggered message.

	CauseRequestAccepted                        causeValue = 16
	CauseRequestAcceptedPartially               causeValue = 17
	CauseNewPDNTypeDueToNetworkPreference       causeValue = 18
	CauseNewPDNTypeDueToSingleAddressBearerOnly causeValue = 19

	// Rejection in a Response / triggered message.

	CauseContextNotFound               causeValue = 64
	CauseInvalidMessageFormat          causeValue = 65
	CauseVersionNotSupportedByNextPeer causeValue = 66
	CauseInvalidLength                 causeValue = 67
	CauseServiceNotSupported           causeValue = 68
	CauseMandatoryIEIncorrect          causeValue = 69
	CauseMandatoryIEMissing            causeValue = 70

	CauseSystemFailure                   causeValue = 72
	CauseNoResourcesAvailable            causeValue = 73
	CauseSemanticErrorInTheTFTOperation  causeValue = 74
	CauseSyntacticErrorInTheTFTOperation causeValue = 75
	CauseSemanticErrorsInPacketFilters   causeValue = 76
	CauseSyntacticErrorsInPacketFilters  causeValue = 77
	CauseMissingOrUnknownAPN             causeValue = 78

	CauseGREKeyNotFound                       causeValue = 80
	CauseRelocationFailure                    causeValue = 81
	CauseDeniedInRAT                          causeValue = 82
	CausePreferredPDNTypeNotSupported         causeValue = 83
	CauseAllDynamicAddressesAreOccupied       causeValue = 85
	CauseUEContextWithoutTFTAlreadyActivated  causeValue = 86
	CauseProtocolTypeNotSupported             causeValue = 87
	CauseUENotResponding                      causeValue = 88
	CauseUERefuses                            causeValue = 89
	CauseUnableToPageUE                       causeValue = 90
	CauseNoMemoryAvailable                    causeValue = 91
	CauseUserAuthenticationFailed             causeValue = 92
	CauseAPNAccessDeniedNoSubscription        causeValue = 93
	CauseRequestRejectedReasonNotSpecified    causeValue = 94
	CausePTMSISignatureMismatch               causeValue = 95
	CauseIMSINotKnown                         causeValue = 96
	CauseSemanticErrorInTheTADOperation       causeValue = 97
	CauseSyntacticErrorInTheTADOperation      causeValue = 98
	CauseServiceDenied                        causeValue = 99
	CauseRemotePeerNotResponding              causeValue = 100
	CauseCollisionWithNetworkInitiatedRequest causeValue = 101
	CauseUnableToPageUEDueToSuspension        causeValue = 102
	CauseConditionalIEMissing                 causeValue = 103

	CauseAPNRestrictionTypeIncompatibleWithCurrentlyActivePDNConnection causeValue = 104
)

type Cause struct {
	header
	value       causeValue
	pce         bool
	bce         bool
	cs          bool
	offendingIe *header
}

func NewCause(instance byte, value causeValue, pce, bce, cs bool, offendingIe *header) (*Cause, error) {
	length := 2
	if offendingIe != nil {
		length = 6
	}

	header, err := newHeader(causeNum, uint16(length), instance)
	if err != nil {
		return nil, err
	}

	return &Cause{
		header:      header,
		value:       value,
		pce:         pce,
		bce:         bce,
		cs:          cs,
		offendingIe: offendingIe,
	}, nil
}

func (c *Cause) Marshal() []byte {
	buf := make([]byte, c.header.length)
	buf[0] = byte(c.value)
	buf[1] = setBit(buf[1], 2, c.pce)
	buf[1] = setBit(buf[1], 1, c.bce)
	buf[1] = setBit(buf[1], 0, c.cs)
	if c.offendingIe != nil {
		buf[2] = byte(c.offendingIe.typeNum)
		buf[5] = c.offendingIe.instance
	}
	return c.header.marshal(buf)
}

func unmarshalCause(h header, buf []byte) (*Cause, error) {
	if h.typeNum != causeNum {
		log.Fatal("Invalud type")
	}

	if len(buf) < 2 {
		return nil, errors.New("too short data")
	}
	value := causeValue(buf[0])
	pce := getBit(buf[1], 2)
	bce := getBit(buf[1], 1)
	cs := getBit(buf[1], 0)

	var offendingIeHeader *header
	if h.length == 6 {
		offendingIeHeader = &header{
			typeNum:  ieTypeNum(buf[2]),
			length:   0,
			instance: buf[5],
		}
	}
	cause, err := NewCause(h.instance, value, pce, bce, cs, offendingIeHeader)
	if err != nil {
		return nil, err
	}
	return cause, nil
}

func (c *Cause) Value() causeValue {
	return c.value
}
func (c *Cause) Pce() bool {
	return c.pce
}
func (c *Cause) Bce() bool {
	return c.bce
}
func (c *Cause) Cs() bool {
	return c.cs
}
