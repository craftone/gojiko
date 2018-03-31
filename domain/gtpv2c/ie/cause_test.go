package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCause(t *testing.T) {
	cause, err := NewCause(0, CauseRequestAccepted, true, false, true, nil)
	assert.Equal(t, causeNum, cause.header.typeNum)
	assert.Equal(t, CauseRequestAccepted, cause.Value())
	assert.Equal(t, true, cause.Pce())
	assert.Equal(t, false, cause.Bce())
	assert.Equal(t, true, cause.Cs())
	assert.Nil(t, cause.offendingIe)
	assert.Nil(t, err)

	offendingIe := &header{2, 0, 3}
	cause, _ = NewCause(1, CauseNoResourcesAvailable, false, true, false, offendingIe)
	assert.Equal(t, CauseNoResourcesAvailable, cause.value)
	assert.Equal(t, false, cause.Pce())
	assert.Equal(t, true, cause.Bce())
	assert.Equal(t, false, cause.Cs())
	assert.Equal(t, ieTypeNum(2), cause.offendingIe.typeNum)
	assert.Equal(t, byte(3), cause.offendingIe.instance)
}

func TestCause_Marshal(t *testing.T) {
	cause, _ := NewCause(1, CauseNoResourcesAvailable, true, false, true, nil)
	causeBin := cause.Marshal()
	assert.Equal(t, []byte{2, 0, 2, 1, 73, 5}, causeBin)

	causeOff, _ := NewCause(1, 2, true, true, true, &header{2, 0, 3})
	causeOffBin := causeOff.Marshal()
	assert.Equal(t, []byte{2, 0, 6, 1, 2, 7, 2, 0, 0, 3}, causeOffBin)
}

func TestUnmarshal_cause(t *testing.T) {
	causeOrg, _ := NewCause(1, CauseNoResourcesAvailable, true, false, true, nil)
	causeBin := causeOrg.Marshal()
	msg, tail, err := Unmarshal(causeBin, CreateSessionRequest)
	cause := msg.(*Cause)
	assert.Equal(t, byte(1), cause.instance)
	assert.Equal(t, CauseNoResourcesAvailable, cause.value)
	assert.Equal(t, true, cause.pce)
	assert.Equal(t, false, cause.bce)
	assert.Equal(t, true, cause.cs)
	assert.Nil(t, cause.offendingIe)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	causeOffOrg, _ := NewCause(1, CauseNoResourcesAvailable, true, false, true, &header{2, 0, 3})
	causeOffBin := causeOffOrg.Marshal()
	msg, tail, err = Unmarshal(causeOffBin, CreateSessionRequest)
	cause = msg.(*Cause)
	assert.Equal(t, byte(1), cause.instance)
	assert.Equal(t, CauseNoResourcesAvailable, cause.value)
	assert.Equal(t, true, cause.pce)
	assert.Equal(t, false, cause.bce)
	assert.Equal(t, true, cause.cs)
	assert.Equal(t, ieTypeNum(2), cause.offendingIe.typeNum)
	assert.Equal(t, byte(3), cause.offendingIe.instance)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}

func TestCauseType_String(t *testing.T) {
	tests := []struct {
		name string
		t    CauseType
		want string
	}{
		{t: CauseTypeRequestInitial, want: "Request Initial"},
		{t: CauseTypeAcceptance, want: "Acceptance"},
		{t: CauseTypeRetryableRejection, want: "Retryable Rejection"},
		{t: CauseTypeRejection, want: "Rejection"},
		{t: CauseTypeOther, want: "Other"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
				t.Errorf("CauseType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCauseValue_Type(t *testing.T) {
	tests := []struct {
		name string
		c    CauseValue
		want CauseType
	}{
		{"", CauseValue(0), CauseTypeOther},
		{"", CauseValue(1), CauseTypeRequestInitial},
		{"", CauseValue(2), CauseTypeRequestInitial},
		{"", CauseValue(15), CauseTypeRequestInitial},
		{"", CauseValue(16), CauseTypeAcceptance},
		{"", CauseValue(63), CauseTypeAcceptance},
		{"", CauseValue(64), CauseTypeRejection},
		{"", CauseValue(73), CauseTypeRetryableRejection},
		{"", CauseValue(239), CauseTypeRejection},
		{"", CauseValue(240), CauseTypeRequestInitial},
		{"", CauseValue(255), CauseTypeRequestInitial},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Type(); got != tt.want {
				t.Errorf("CauseValue.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCauseValue_Detail(t *testing.T) {
	tests := []struct {
		name string
		c    CauseValue
		want string
	}{
		{"", CauseValue(0), "Unknown cause"},
		{"", CauseValue(1), "Unknown cause"},
		{"", CauseValue(2), "Local Detach"},
		{"", CauseValue(15), "Unknown cause"},
		{"", CauseValue(16), "Request accepted"},
		{"", CauseValue(63), "Unknown cause"},
		{"", CauseValue(64), "Context Not Found"},
		{"", CauseValue(239), "Unknown cause"},
		{"", CauseValue(240), "Unknown cause"},
		{"", CauseValue(255), "Unknown cause"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Detail(); got != tt.want {
				t.Errorf("CauseValue.Detail() = %v, want %v", got, tt.want)
			}
		})
	}
}
