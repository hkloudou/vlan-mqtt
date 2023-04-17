package main

import (
	"io"

	"github.com/hkloudou/xtransport"
	"github.com/hkloudou/xtransport/packets/mqtt"
	transport "github.com/hkloudou/xtransport/transports/tcp"
)

func main() {
	tran := transport.NewTransport("tcp", xtransport.Secure(false))
	cli, err := tran.Dial("broker.emqx.io:1833")
	if err != nil {
		panic(err)
	}
	defer func() {
		if r := recover(); r != nil {
			println(r)
		}
		cli.Close()
	}()

	connPacket := mqtt.NewControlPacket(mqtt.Connect).(*mqtt.ConnectPacket)
	// connPacket.
	connPacket.ClientIdentifier = "mqttx_517cc888"
	connPacket.Keepalive = 60
	connPacket.CleanSession = true
	cli.Send(connPacket)
	for {
		request, err := cli.Recv(func(r io.Reader) (interface{}, error) {
			i, err := mqtt.ReadPacket(r)
			return i, err
		})
		if err != nil {
			return
		}
		if request.(mqtt.ControlPacket).Type() <= 0 || request.(mqtt.ControlPacket).Type() >= 14 {
			cli.Close()
			return
		}
		switch request.(mqtt.ControlPacket).Type() {
		case mqtt.Pingreq:
			cli.Send(mqtt.NewControlPacket(mqtt.Pingresp))
			continue
		}
	}
	// if err := tran.Listen("broker.emqx.io:1833"); err != nil {
	// 	panic(err)
	// }
	// tran.Accept(func(sock xtransport.Socket) {
	// 	defer func() {
	// 		if r := recover(); r != nil {
	// 			println(r)
	// 		}
	// 		sock.Close()
	// 	}()
	// 	for {
	// 		request, err := sock.Recv(func(r io.Reader) (interface{}, error) {
	// 			i, err := mqtt.ReadPacket(r)
	// 			return i, err
	// 		})
	// 		if err != nil {
	// 			return
	// 		}
	// 		if request == nil {
	// 			continue
	// 		}
	// 		// log.Println("recv", request.String())
	// 		if request.(mqtt.ControlPacket).Type() <= 0 || request.(mqtt.ControlPacket).Type() >= 14 {
	// 			sock.Close()
	// 			return
	// 		}
	// 		switch request.(mqtt.ControlPacket).Type() {
	// 		case mqtt.Pingreq:
	// 			sock.Send(mqtt.NewControlPacket(mqtt.Pingresp))
	// 			break
	// 		case mqtt.Connect:
	// 			_hook.OnClientConnect(sock, request.(*mqtt.ConnectPacket))
	// 			break
	// 		case mqtt.Subscribe:
	// 			_hook.OnClientSubcribe(sock, request.(*mqtt.SubscribePacket))
	// 			break
	// 		case mqtt.Unsubscribe:
	// 			_hook.OnClientUnSubcribe(sock, request.(*mqtt.UnsubscribePacket))
	// 			break
	// 		case mqtt.Publish:
	// 			_hook.OnClientPublish(sock, request.(*mqtt.PublishPacket))
	// 			break
	// 		default:
	// 			// return nil, fmt.Errorf("not support packet type:%d", data.Type())
	// 		}
	// 	}
	// })
}
