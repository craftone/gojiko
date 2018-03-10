package stats

import (
	"fmt"

	humanize "github.com/dustin/go-humanize"
)

func FormatSIFloat(value float64, unit string) string {
	siFloat, siPrefix := humanize.ComputeSI(value)
	return fmt.Sprintf("%.1f %s%s", siFloat, siPrefix, unit)
}

func FormatSIUint(value uint64, unit string) string {
	siFloat, siPrefix := humanize.ComputeSI(float64(value))
	if siPrefix == "" {
		return fmt.Sprintf("%d %s", value, unit)
	}
	return fmt.Sprintf("%.1f %s%s", siFloat, siPrefix, unit)
}
