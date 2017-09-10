package ie

import (
	"reflect"
	"testing"
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
