package ie

import (
	"errors"

	"github.com/craftone/gojiko/util"
)

type CauseType int

const (
	CauseTypeRequestInitial CauseType = iota
	CauseTypeAcceptance
	CauseTypeRetryableRejection
	CauseTypeRejection
	CauseTypeOther
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
	CauseAllDynamicAddressesAreOccupied       CauseValue = 84
	CauseUEContextWithoutTFTAlreadyActivated  CauseValue = 85
	CauseProtocolTypeNotSupported             CauseValue = 86
	CauseUENotResponding                      CauseValue = 87
	CauseUERefuses                            CauseValue = 88
	CauseServiceDenied                        CauseValue = 89
	CauseUnableToPageUE                       CauseValue = 90
	CauseNoMemoryAvailable                    CauseValue = 91
	CauseUserAuthenticationFailed             CauseValue = 92
	CauseAPNAccessDeniedNoSubscription        CauseValue = 93
	CauseRequestRejectedReasonNotSpecified    CauseValue = 94
	CausePTMSISignatureMismatch               CauseValue = 95
	CauseIMSINotKnown                         CauseValue = 96
	CauseSemanticErrorInTheTADOperation       CauseValue = 97
	CauseSyntacticErrorInTheTADOperation      CauseValue = 98
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
		log.Panic("Invalid type")
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

var causeToStrMap = map[CauseValue]string{
	// Request / Initial message

	CauseLocalDetach:                          "Local Detach",
	CauseCompleteDetach:                       "Complete Detach",
	CauseRATChangedFrom3GPPtoNon3GPP:          "RAT change from 3GPP to Non-3GPP",
	CauseISRDeactivation:                      "ISR deactivation",
	CauseErrorIndicationReceivedFromRNCeNodeB: "Error Indication received from RNC/eNodeB",
	CauseIMSIDetachOnly:                       "IMSI Detach Only",
	CauseReactivationRequested:                "Reactivation Requested",
	CausePDNReconnectionToThisAPNDisallowed:   "PDN reconnection to this APN disallowed",
	CauseAccessChangedFromNon3GPPto3GPP:       "Access changed from Non-3GPP to 3GPP",
	CausePDNConnectionInactivityTimerExpires:  "PDN connection inactivity timer expires",

	// Acceptance in a Response / triggered message.

	CauseRequestAccepted:                        "Request accepted",
	CauseRequestAcceptedPartially:               "Request accepted partially",
	CauseNewPDNTypeDueToNetworkPreference:       "New PDN type due to network preference",
	CauseNewPDNTypeDueToSingleAddressBearerOnly: "New PDN type due to single address bearer only",

	// Rejection in a Response / triggered message.

	CauseContextNotFound:               "Context Not Found",
	CauseInvalidMessageFormat:          "Invalid Message Format",
	CauseVersionNotSupportedByNextPeer: "Version not supported by next peer",
	CauseInvalidLength:                 "Invalid length",
	CauseServiceNotSupported:           "Service not supported",
	CauseMandatoryIEIncorrect:          "Mandatory IE incorrect",
	CauseMandatoryIEMissing:            "Mandatory IE missing",

	CauseSystemFailure:                   "System failure",
	CauseNoResourcesAvailable:            "No resources available",
	CauseSemanticErrorInTheTFTOperation:  "Semantic error in the TFT operation",
	CauseSyntacticErrorInTheTFTOperation: "Syntactic error in the TFT operation",
	CauseSemanticErrorsInPacketFilters:   "Semantic errors in packet filter(s)",
	CauseSyntacticErrorsInPacketFilters:  "Syntactic errors in packet filter(s)",
	CauseMissingOrUnknownAPN:             "Missing or unknown APN",

	CauseGREKeyNotFound:                       "GRE key not found",
	CauseRelocationFailure:                    "Relocation failure",
	CauseDeniedInRAT:                          "Denied in RAT",
	CausePreferredPDNTypeNotSupported:         "Preferred PDN type not supported",
	CauseAllDynamicAddressesAreOccupied:       "All dynamic addresses are occupied",
	CauseUEContextWithoutTFTAlreadyActivated:  "UE context without TFT already activated",
	CauseProtocolTypeNotSupported:             "Protocol type not supported",
	CauseUENotResponding:                      "UE not responding",
	CauseUERefuses:                            "UE refuses",
	CauseServiceDenied:                        "Service denied",
	CauseUnableToPageUE:                       "Unable to page UE",
	CauseNoMemoryAvailable:                    "No memory available",
	CauseUserAuthenticationFailed:             "User authentication failed",
	CauseAPNAccessDeniedNoSubscription:        "APN access denied – no subscription",
	CauseRequestRejectedReasonNotSpecified:    "Request rejected (reason not specified)",
	CausePTMSISignatureMismatch:               "P-TMSI Signature mismatch",
	CauseIMSINotKnown:                         "IMSI not known",
	CauseSemanticErrorInTheTADOperation:       "Semantic error in the TAD operation",
	CauseSyntacticErrorInTheTADOperation:      "Syntactic error in the TAD operation",
	CauseRemotePeerNotResponding:              "Remote peer not responding",
	CauseCollisionWithNetworkInitiatedRequest: "Collision with network initiated request",
	CauseUnableToPageUEDueToSuspension:        "Unable to page UE due to Suspension",
	CauseConditionalIEMissing:                 "Conditional IE missing",

	CauseAPNRestrictionTypeIncompatibleWithCurrentlyActivePDNConnection: "APN Restriction type Incompatible with currently active PDN connection",
}

func (c CauseValue) Type() CauseType {
	switch {
	case byte(c) == 0:
		return CauseTypeOther
	case byte(c) <= 15:
		return CauseTypeRequestInitial
	case byte(c) <= 63:
		return CauseTypeAcceptance
	case c == CauseNoResourcesAvailable ||
		c == CauseAllDynamicAddressesAreOccupied ||
		c == CauseNoMemoryAvailable ||
		c == CauseMissingOrUnknownAPN ||
		c == CauseAPNAccessDeniedNoSubscription ||
		c == CauseRequestRejectedReasonNotSpecified:
		return CauseTypeRetryableRejection
	case byte(c) <= 239:
		return CauseTypeRejection
	default:
		return CauseTypeRequestInitial
	}
}

func (c CauseValue) Detail() string {
	if str, ok := causeToStrMap[c]; ok {
		return str
	}
	return "Unknown cause"
}

func (t CauseType) String() string {
	switch t {
	case CauseTypeRequestInitial:
		return "Request Initial"
	case CauseTypeAcceptance:
		return "Acceptance"
	case CauseTypeRetryableRejection:
		return "Retryable Rejection"
	case CauseTypeRejection:
		return "Rejection"
	default:
		return "Other"
	}
}
