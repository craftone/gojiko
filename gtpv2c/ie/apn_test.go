package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApn(t *testing.T) {
	apn, _ := NewApn(0, "example.com")
	assert.Equal(t, apnNum, apn.typeNum)
	assert.Equal(t, "example.com", apn.Value)

	_, err := NewApn(0, "")
	assert.Error(t, err)

	_, err = NewApn(0, "example.c*om")
	assert.Error(t, err)
}

func TestApn_Marshal(t *testing.T) {
	apn, _ := NewApn(1, "example.com")
	apnBin := apn.Marshal()
	assert.Equal(t, []byte{0x47, 0, 0x0b, 1, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d}, apnBin)
}

func TestUnmarshal_apn(t *testing.T) {
	apnOrg, _ := NewApn(1, "example.com")
	apnBin := apnOrg.Marshal()
	msg, tail, err := Unmarshal(apnBin, MsToNetwork)
	apn := msg.(*Apn)
	assert.Equal(t, byte(1), apn.instance)
	assert.Equal(t, "example.com", apn.Value)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}

func TestIsValidAPN(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{"a"}, true},
		{"", args{"a*"}, false},
		{"", args{"abcdefghijklmnopqrstuvwxyz" +
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
			"0123456789" + ".-"},
			true,
		},
		{"", args{"a.b"}, true},
		{"", args{".a.b"}, false},
		{"", args{"a.b."}, false},
		{"", args{".a.b."}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidAPN(tt.args.s); got != tt.want {
				t.Errorf("IsValidAPN() = %v, want %v", got, tt.want)
			}
		})
	}
}
