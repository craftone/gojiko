package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

const (
	stats_flow_buffer_length_key = "stats.flow.buffer.length"
)

var (
	stats_flow_buffer_length int
)

func loadStatsConf() {
	// Set default values
	viper.SetDefault(stats_flow_buffer_length_key, 100)

	// check and load configs
	err := setStatsFlowBufferLength(viper.GetInt(stats_flow_buffer_length_key))
	if err != nil {
		log.Fatal(err)
	}
}

func setStatsFlowBufferLength(val int) error {
	if val < 0 {
		return fmt.Errorf("Invalid %s : %d", stats_flow_buffer_length_key, val)
	}
	stats_flow_buffer_length = val
	return nil
}

func StatsFlowBufferLength() int {
	return stats_flow_buffer_length
}
