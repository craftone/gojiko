package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

const (
	flow_udpecho_waitreceive_key = "flow.UdpEcho.WaitReceive"
)

var (
	flow_udpecho_waitreceive time.Duration
)

func loadFlowConf() {
	// Set default values
	viper.SetDefault(flow_udpecho_waitreceive_key, 1000) // msec

	// check and load configs
	err := setFlowUdpEchoWaitReceive(viper.GetInt(flow_udpecho_waitreceive_key))
	if err != nil {
		log.Fatal(err)
	}
}

func setFlowUdpEchoWaitReceive(val int) error {
	if val < 0 {
		return fmt.Errorf("Invalid %s : %d", flow_udpecho_waitreceive_key, val)
	}
	flow_udpecho_waitreceive = time.Duration(val) * time.Millisecond
	return nil
}

func FlowUdpEchoWaitReceive() time.Duration {
	return flow_udpecho_waitreceive
}
