package ie

import (
	"testing"
)

func Test_getBit(t *testing.T) {
	type args struct {
		b   byte
		idx uint
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{0, 0}, false},
		{"", args{1, 0}, true},
		{"", args{2, 0}, false},
		{"", args{3, 1}, true},
		{"", args{0x0f, 7}, false},
		{"", args{0x80, 7}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBit(tt.args.b, tt.args.idx); got != tt.want {
				t.Errorf("getBit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setBit(t *testing.T) {
	type args struct {
		b    byte
		idx  uint
		flag bool
	}
	tests := []struct {
		name string
		args args
		want byte
	}{
		{"", args{0, 0, false}, 0},
		{"", args{0, 0, true}, 1},
		{"", args{1, 0, false}, 0},
		{"", args{1, 0, true}, 1},
		{"", args{0, 7, false}, 0},
		{"", args{0, 7, true}, 0x80},
		{"", args{0x80, 7, false}, 0},
		{"", args{0x80, 7, true}, 0x80},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setBit(tt.args.b, tt.args.idx, tt.args.flag); got != tt.want {
				t.Errorf("setBit() = %v, want %v", got, tt.want)
			}
		})
	}
}
