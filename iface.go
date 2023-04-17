package main

import (
	"fmt"
	"log"
	"net"

	"github.com/songgao/water"
)

type vlanFace struct {
	vid            int
	BroadcastTopic string
	EthTopic       string
	Face           *water.Interface
}

func initIface(vid int) (*vlanFace, error) {
	// create a TAP interface
	config := water.Config{
		DeviceType: water.TAP,
	}
	config.Name = fmt.Sprintf("vnats%d", vid)
	ifce, err := water.New(config)
	if err != nil {
		return nil, err
	}
	// get ethernet address of the interface we just created
	var ownEth net.HardwareAddr
	nifces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, nifce := range nifces {
		if nifce.Name == config.Name {
			ownEth = nifce.HardwareAddr
			break
		}
	}
	if len(ownEth) == 0 {
		log.Fatal("failed to get own ethernet address")
	}

	// config
	return &vlanFace{
		vid:            vid,
		BroadcastTopic: _getBroadcastTopic(vid),
		EthTopic:       _getEthTopic(vid, ownEth),
		Face:           ifce,
	}, nil
}

func (m *vlanFace) IsTopic(topic string) bool {
	return (m.BroadcastTopic == topic || m.EthTopic == topic)
}

func _getEthTopic(vid int, eth net.HardwareAddr) string {
	return fmt.Sprintf("vvvv.xxxx.%d.%x", vid, eth)
}

func _getBroadcastTopic(vid int) string {
	return fmt.Sprintf("vvvv.xxxx.%d", vid)
}
