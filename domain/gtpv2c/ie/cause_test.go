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

func TestCauseDetail(t *testing.T) {
	type args struct {
		c CauseValue
	}
	tests := []struct {
		name  string
		args  args
		want  CauseType
		want1 string
	}{
		{"", args{CauseValue(0)}, CauseTypeOther, "Unknown cause"},
		{"", args{CauseValue(1)}, CauseTypeRequestInitial, "Unknown cause"},
		{"", args{CauseValue(2)}, CauseTypeRequestInitial, "Local Detach"},
		{"", args{CauseValue(15)}, CauseTypeRequestInitial, "Unknown cause"},
		{"", args{CauseValue(16)}, CauseTypeAcceptance, "Request accepted"},
		{"", args{CauseValue(63)}, CauseTypeAcceptance, "Unknown cause"},
		{"", args{CauseValue(64)}, CauseTypeRejection, "Context Not Found"},
		{"", args{CauseValue(239)}, CauseTypeRejection, "Unknown cause"},
		{"", args{CauseValue(240)}, CauseTypeRequestInitial, "Unknown cause"},
		{"", args{CauseValue(255)}, CauseTypeRequestInitial, "Unknown cause"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := CauseDetail(tt.args.c)
			if got != tt.want {
				t.Errorf("CauseDetail() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CauseDetail() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
