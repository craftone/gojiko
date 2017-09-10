package ie

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newMccMnc(t *testing.T) {
	type args struct {
		mcc string
		mnc string
	}
	tests := []struct {
		name    string
		args    args
		want    mccMnc
		wantErr bool
	}{
		{"", args{"440", "10"}, mccMnc{"440", "10", [3]byte{0x44, 0xf0, 0x01}}, false},
		{"", args{"440", "210"}, mccMnc{"440", "210", [3]byte{0x44, 0x00, 0x12}}, false},
		{"error", args{"4401", "210"}, mccMnc{}, true},
		{"error", args{"44", "210"}, mccMnc{}, true},
		{"error", args{"440", "1"}, mccMnc{}, true},
		{"error", args{"440", "3210"}, mccMnc{}, true},
		{"error", args{"44a", "10"}, mccMnc{}, true},
		{"error", args{"440", "*10"}, mccMnc{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newMccMnc(tt.args.mcc, tt.args.mnc)
			if (err != nil) != tt.wantErr {
				t.Errorf("newMccMnc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newMccMnc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUli(t *testing.T) {
	ecgi, _ := NewEcgi("440", "10", 0x22D6600)
	tai, _ := NewTai("440", "10", 0x1421)
	uliArg := UliArg{
		Ecgi: ecgi,
		Tai:  tai,
	}
	uli, err := NewUli(0, uliArg)
	assert.Equal(t, uliNum, uli.typeNum)
	assert.Equal(t, uint16(13), uli.length)
	assert.Equal(t, uint32(0x22D6600), uli.Ecgi.Eci)
	assert.Equal(t, uint16(0x1421), uli.Tai.Tac)
	assert.Nil(t, err)
}

func TestUli_Marshal(t *testing.T) {
	ecgi, _ := NewEcgi("440", "10", 0x22D6600)
	tai, _ := NewTai("440", "10", 0x1421)
	uliArg := UliArg{
		Ecgi: ecgi,
		Tai:  tai,
	}
	uli, _ := NewUli(0, uliArg)
	uliBin := uli.Marshal()
	assert.Equal(t, []byte{
		0x56, 0x00, 0x0D, 0x00, 0x18, 0x44, 0xF0, 0x01,
		0x14, 0x21, 0x44, 0xF0, 0x01, 0x02, 0x2D, 0x66,
		0x00,
	}, uliBin)
}

func TestUnmarshal_uli(t *testing.T) {
	ecgi, _ := NewEcgi("440", "10", 0x22D6600)
	tai, _ := NewTai("440", "10", 0x1421)
	uliArg := UliArg{
		Ecgi: ecgi,
		Tai:  tai,
	}
	uli, _ := NewUli(1, uliArg)
	uliBin := uli.Marshal()
	msg, tail, err := Unmarshal(uliBin)
	uli = msg.(*Uli)
	assert.Equal(t, byte(1), uli.instance)
	assert.Equal(t, "440", uli.Ecgi.Mcc)
	assert.Equal(t, "10", uli.Ecgi.Mnc)
	assert.Equal(t, uint32(0x22D6600), uli.Ecgi.Eci)
	assert.Equal(t, "440", uli.Tai.Mcc)
	assert.Equal(t, "10", uli.Tai.Mnc)
	assert.Equal(t, uint16(0x1421), uli.Tai.Tac)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}
