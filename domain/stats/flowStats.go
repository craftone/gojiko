package stats

import (
	"context"

	"github.com/craftone/gojiko/config"
)

type FlowStats struct {
	*absStats
}

func NewFlowStats(ctx context.Context) *FlowStats {
	return &FlowStats{newAbsStats(ctx, config.StatsFlowBufferLength())}
}
