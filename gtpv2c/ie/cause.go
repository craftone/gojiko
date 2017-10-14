package ie

import (
	"errors"
	"log"

	"github.com/craftone/gojiko/util"
)

type CauseValue byte

const (
	// Request / Initial message

	CauseLocalDetach                          CauseValue = 2
	CauseCompleteDetach                       CauseValue = 3
	CauseRATChangedFrom3GPPtoNon3GPP          CauseValue = 4
	CauseISRDeactivation                      CauseValue = 5
	CauseErrorIndicationReceivedFromRNCeNodeB CauseValue = 6
	CauseIMSIDetachOnly                       CauseValue = 7
	CauseReactivationRequested                CauseValue = 8
	CausePDNReconnectionToThisAPNDisallowed   CauseValue = 9
	CauseAccessChangedFromNon3GPPto3GPP       CauseValue = 10
	CausePDNConnectionInactivityTimerExpires  CauseValue = 11

	// Acceptance in a Response / triggered message.

	CauseRequestAccepted                        CauseValue = 16
	CauseRequestAcceptedPartially               CauseValue = 17
	CauseNewPDNTypeDueToNetworkPreference       CauseValue = 18
	CauseNewPDNTypeDueToSingleAddressBearerOnly CauseValue = 19

	// Rejection in a Response / triggered message.

	CauseContextNotFound               CauseValue = 64
	CauseInvalidMessageFormat          CauseValue = 65
	CauseVersionNotSupportedByNextPeer CauseValue = 66
	CauseInvalidLength                 CauseValue = 67
	CauseServiceNotSupported           CauseValue = 68
	CauseMandatoryIEIncorrect          CauseValue = 69
	CauseMandatoryIEMissing            CauseValue = 70

	CauseSystemFailure                   CauseValue = 72
	CauseNoResourcesAvailable            CauseValue = 73
	CauseSemanticErrorInTheTFTOperation  CauseValue = 74
	CauseSyntacticErrorInTheTFTOperation CauseValue = 75
	CauseSemanticErrorsInPacketFilters   CauseValue = 76
	CauseSyntacticErrorsInPacketFilters  CauseValue = 77
	CauseMissingOrUnknownAPN             CauseValue = 78

	CauseGREKeyNotFound                       CauseValue = 80
	CauseRelocationFailure                    CauseValue = 81
	CauseDeniedInRAT                          CauseValue = 82
	CausePreferredPDNTypeNotSupported         CauseValue = 83
	CauseAllDynamicAddressesAreOccupied       CauseValue = 85
	CauseUEContextWithoutTFTAlreadyActivated  CauseValue = 86
	CauseProtocolTypeNotSupported             CauseValue = 87
	CauseUENotResponding                      CauseValue = 88
	CauseUERefuses                            CauseValue = 89
	CauseUnableToPageUE                       CauseValue = 90
	CauseNoMemoryAvailable                    CauseValue = 91
	CauseUserAuthenticationFailed             CauseValue = 92
	CauseAPNAccessDeniedNoSubscription        CauseValue = 93
	CauseRequestRejectedReasonNotSpecified    CauseValue = 94
	CausePTMSISignatureMismatch               CauseValue = 95
	CauseIMSINotKnown                         CauseValue = 96
	CauseSemanticErrorInTheTADOperation       CauseValue = 97
	CauseSyntacticErrorInTheTADOperation      CauseValue = 98
	CauseServiceDenied                        CauseValue = 99
	CauseRemotePeerNotResponding              CauseValue = 100
	CauseCollisionWithNetworkInitiatedRequest CauseValue = 101
	CauseUnableToPageUEDueToSuspension        CauseValue = 102
	CauseConditionalIEMissing                 CauseValue = 103

	CauseAPNRestrictionTypeIncompatibleWithCurrentlyActivePDNConnection CauseValue = 104
)

type Cause struct {
	header
	value       CauseValue
	pce         bool
	bce         bool
	cs          bool
	offendingIe *header
}

func NewCause(instance byte, value CauseValue, pce, bce, cs bool, offendingIe *header) (*Cause, error) {
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
	buf[1] = util.SetBit(buf[1], 2, c.pce)
	buf[1] = util.SetBit(buf[1], 1, c.bce)
	buf[1] = util.SetBit(buf[1], 0, c.cs)
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
	value := CauseValue(buf[0])
	pce := util.GetBit(buf[1], 2)
	bce := util.GetBit(buf[1], 1)
	cs := util.GetBit(buf[1], 0)

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

func (c *Cause) Value() CauseValue {
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
