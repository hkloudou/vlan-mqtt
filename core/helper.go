package core

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// SetInterfaceStatus brings up or down a network interface.
func SetInterfaceStatus(iName string, up bool, debug bool) error {
	statusString := "down"
	if up {
		statusString = "up"
	}
	// sargs := fmt.Sprintf("link set dev %s %s mtu %d qlen %d", iName, statusString, devMtuSize, devTxQueLen)
	// args := strings.Split(sargs, " ")
	sargs := fmt.Sprintf("link set dev %s %s", iName, statusString)
	args := strings.Split(sargs, " ")
	return commandExec("ip", args, debug)
}

// SetDevIP sets the local IP address of a network interface.
func SetDevIP(iName string, localAddr net.IP, broadcast string, debug bool) error {
	// sargs := fmt.Sprintf("%s %s netmask %s", iName, localAddr.String(), net.IP(addr.Mask).String())
	// args := strings.Split(sargs, " ")
	// return commandExec("ifconfig", args, debug)
	sargs := fmt.Sprintf("addr add %s broadcast %s dev %s", localAddr.String(), broadcast, iName)
	args := strings.Split(sargs, " ")
	return commandExec("ip", args, debug)
	// ip addr add 10.1.0.{CHANGEME} broadcast 10.1.255.255 dev vnats0

}

// SetDefaultGateway sets the systems gateway to the IP / device specified.
func SetDefaultGateway(gw, iName string, debug bool) error {
	sargs := fmt.Sprintf("add default gw %s dev %s", gw, iName)
	args := strings.Split(sargs, " ")
	return commandExec("route", args, debug)
}

// AddRoute routes all traffic for addr via interface iName.
func AddRoute(addr *net.IPNet, iName string, debug bool) error {
	// sargs := fmt.Sprintf("add %s gw %s dev %s", addr.String(), viaAddr.String(), iName)
	// args := strings.Split(sargs, " ")
	//sudo ip route add 10.1.0.0/16 dev vnats0
	// return commandExec("route", args, debug)
	//route add 10.1.0.0/16 dev vnats0
	sargs := fmt.Sprintf("route add %s dev %s", addr.String(), iName)
	args := strings.Split(sargs, " ")
	return commandExec("ip", args, debug)
}

// DelRoute deletes the route in the system routing table to a specific destination.
// func DelRoute(addr, viaAddr net.IP, iName string, debug bool) error {
// 	sargs := fmt.Sprintf("del %s gw %s dev %s", addr.String(), viaAddr.String(), iName)
// 	args := strings.Split(sargs, " ")
// 	return commandExec("route", args, debug)
// }

// GetNetGateway return net gateway (default route) and nic.
// Credit: https://github.com/bigeagle/gohop/blob/master/hop/iface.go
func GetNetGateway() (gw, dev string, err error) {
	file, err := os.Open("/proc/net/route")
	if err != nil {
		return "", "", err
	}

	defer file.Close()
	rd := bufio.NewReader(file)

	s2byte := func(s string) byte {
		b, _ := strconv.ParseUint(s, 16, 8)
		return byte(b)
	}

	for {
		line, isPrefix, err := rd.ReadLine()

		if err != nil {
			return "", "", err
		}
		if isPrefix {
			return "", "", errors.New("Parse error: Line too long")
		}
		buf := bytes.NewBuffer(line)
		scanner := bufio.NewScanner(buf)
		scanner.Split(bufio.ScanWords)
		tokens := make([]string, 0, 8)

		for scanner.Scan() {
			tokens = append(tokens, scanner.Text())
		}

		iface := tokens[0]
		dest := tokens[1]
		gw := tokens[2]
		mask := tokens[7]

		if bytes.Equal([]byte(dest), []byte("00000000")) &&
			bytes.Equal([]byte(mask), []byte("00000000")) {
			a := s2byte(gw[6:8])
			b := s2byte(gw[4:6])
			c := s2byte(gw[2:4])
			d := s2byte(gw[0:2])

			ip := net.IPv4(a, b, c, d)

			return ip.String(), iface, nil
		}

	}
}
