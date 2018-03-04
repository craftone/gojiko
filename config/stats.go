package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

const (
	stats_flow_buffer_length_str = "stats.flow.buffer.length"
)

var (
	stats_flow_buffer_length int
)

func initStats() {
	// Set default values
	viper.SetDefault(stats_flow_buffer_length_str, 100)

	// check and load configs
	err := setStatsFlowBufferLength(viper.GetInt(stats_flow_buffer_length_str))
	if err != nil {
		log.Fatal(err)
	}
}

func setStatsFlowBufferLength(val int) error {
	if val < 0 {
		return fmt.Errorf("Invalid %s : %d", stats_flow_buffer_length_str, val)
	}
	stats_flow_buffer_length = val
	return nil
}

func StatsFlowBufferLength() int {
	return stats_flow_buffer_length
}
