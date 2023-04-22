package main

import (
	"fmt"
	"log"
	"net"

	"github.com/hkloudou/vlan-mqtt/core"
	"github.com/seancfoley/ipaddress-go/ipaddr"
	"github.com/songgao/water"
)

type vlanFace struct {
	vid            int
	BroadcastTopic string
	EthTopic       string
	Face           *water.Interface
	OwnEth         net.HardwareAddr
}

func initIface(vid int, ipstr string) (*vlanFace, error) {
	localAddr, addr, err := net.ParseCIDR(ipstr)
	// net.par
	if err != nil {
		return nil, err
	}
	// ipaddr.new
	addr2, _ := ipaddr.NewIPAddressFromNetIPNet(addr)
	bcast, _ := addr2.ToIPv4().ToBroadcastAddress()
	// log.Panicln(bcast.GetNetIP().String())

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
	/*
			设置ip
			//ip addr add 10.1.0.{CHANGEME} broadcast 10.1.255.255 dev vnats0
				sargs := fmt.Sprintf("%s %s netmask %s", iName, localAddr.String(), net.IP(addr.Mask).String())
		args := strings.Split(sargs, " ")
		return commandExec("ifconfig", args, debug)
	*/

	// core.SetDevIP(config.Name, net.ParseIP("10.1.0.1"), net.ParseIP("10.1.0.1").DefaultMask().String())

	err = core.SetDevIP(config.Name, localAddr, bcast.GetNetIP().String(), true)
	if err != nil {
		return nil, err
	}
	err = core.SetInterfaceStatus(config.Name, true, true)
	if err != nil {
		return nil, err
	}
	err = core.AddRoute(addr, config.Name, true)
	if err != nil {
		return nil, err
	}
	// config
	return &vlanFace{
		vid:            vid,
		BroadcastTopic: _getBroadcastTopic(vid),
		EthTopic:       _getEthTopic(vid, ownEth),
		OwnEth:         ownEth,
		Face:           ifce,
	}, nil
}

func releaseIface(vid int) {

}

func (m *vlanFace) IsTopic(topic string) bool {
	return (m.BroadcastTopic == topic || m.EthTopic == topic)
}

func _getEthTopic(vid int, eth net.HardwareAddr) string {
	return fmt.Sprintf("vvvv.%d.%x", vid, eth)
}

func _getBroadcastTopic(vid int) string {
	return fmt.Sprintf("vvvv.%d", vid)
}
