package domain

import "testing"

func TestGsResCode_String(t *testing.T) {
	tests := []struct {
		name string
		c    GsResCode
		want string
	}{
		{c: GsResOK, want: "OK"},
		{c: GsResNG, want: "NG"},
		{c: GsResRetryableNG, want: "Retryable NG"},
		{c: GsResTimeout, want: "Timeout"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("GsResCode.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
