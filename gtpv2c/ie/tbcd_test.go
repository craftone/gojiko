package ie

import (
	"reflect"
	"testing"
)

func Test_parseTBCD(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    tbcd
		wantErr bool
	}{
		{"", args{"1"}, tbcd([]byte{0xf1}), false},
		{"", args{"01"}, tbcd([]byte{0x10}), false},
		{"", args{"012345678901234"}, tbcd([]byte{
			0x10, 0x32, 0x54, 0x76, 0x98, 0x10, 0x32, 0xf4,
		}), false},
		{"blank", args{""}, tbcd([]byte{}), false},
		{"error", args{"0123456789a"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTBCD(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTBCD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTBCD() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tbcd_String(t *testing.T) {
	tests := []struct {
		name string
		t    tbcd
		want string
	}{
		{"", tbcd([]byte{0xf1}), "1"},
		{"", tbcd([]byte{0x10}), "01"},
		{"", tbcd([]byte{
			0x10, 0x32, 0x54, 0x76, 0x98, 0x10, 0x32, 0xf4,
		}), "012345678901234"},
		{"blank", tbcd([]byte{}), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
				t.Errorf("tbcd.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unmarshalTbcd(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"", args{[]byte{0xf1}}, "1", false},
		{"", args{[]byte{0x10}}, "01", false},
		{"", args{[]byte{0x1a}}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unmarshalTbcd(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalTbcd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("unmarshalTbcd() = %v, want %v", got, tt.want)
			}
		})
	}
}
