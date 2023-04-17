package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/hkloudou/xlib/xcolor"
	"github.com/hkloudou/xlib/xflag"
	"github.com/hkloudou/xtransport"
	"github.com/hkloudou/xtransport/packets/mqtt"
	transport "github.com/hkloudou/xtransport/transports/tcp"
	"github.com/songgao/water/waterutil"
)

func main() {
	app := xflag.NewApp()
	app.Flags = append(app.Flags, &xflag.IntFlag{
		Name:     "vid",
		Required: true,
		Value:    0,
		Usage:    "vlan id",
	})

	app.Flags = append(app.Flags, &xflag.StringFlag{
		Name:     "mqtt",
		Required: true,
		Value:    "",
		Usage:    "mqtt server",
	})
	app.Action = func(ctx *xflag.Context) error {
		f, err := initIface(ctx.Int("vid"))
		tran := transport.NewTransport("tcp", xtransport.Secure(false))
		cli, err := tran.Dial(ctx.String("mqtt"), xtransport.WithTimeout(5*time.Second))
		if err != nil {
			log.Panicln(xcolor.Red("dial"), err)
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
		connPacket.ProtocolName = "MQTT"
		connPacket.ProtocolVersion = 4
		err = cli.Send(connPacket)
		if err != nil {
			log.Panicln(xcolor.Red("Send"), err)
		}

		go func() {
			var frame [1500]byte
			for {
				// read frame from interface
				n, err := f.Face.Read(frame[:])
				if err != nil {
					log.Fatal(err)
				}
				frame2 := frame[:n]
				// the topic to publish to
				dst := waterutil.MACDestination(frame2)
				log.Println(xcolor.Green("DD"), fmt.Sprintf("%x=>%x", waterutil.MACSource(frame2), waterutil.MACDestination(frame2)))
				var pubTopic string
				if waterutil.IsBroadcast(dst) {
					pubTopic = f.BroadcastTopic
				} else {
					pubTopic = fmt.Sprintf("vvvv.xxxx.%d.%x", ctx.Int("vid"), dst)
				}

				// publish
				pub := mqtt.NewControlPacket(mqtt.Publish).(*mqtt.PublishPacket)
				pub.TopicName = pubTopic
				pub.Payload = frame2
				if err := cli.Send(pub); err != nil {
					log.Fatal(err)
				}
			}
		}()
		for {
			request, err := cli.Recv(func(r io.Reader) (interface{}, error) {
				i, err := mqtt.ReadPacket(r)
				return i, err
			})
			if err != nil {
				log.Panicln(xcolor.Red("Recv"), err)
				break
			}
			log.Println(xcolor.Green("D"), request)
			if request.(mqtt.ControlPacket).Type() <= 0 || request.(mqtt.ControlPacket).Type() >= 14 {
				cli.Close()
				break
			}
			switch request.(mqtt.ControlPacket).Type() {
			case mqtt.Pingreq:
				cli.Send(mqtt.NewControlPacket(mqtt.Pingresp))
				continue
			case mqtt.Publish:
				// request.(mqtt.PublishPacket).Payload
				req := request.(mqtt.PublishPacket)
				if f.IsTopic(req.TopicName) {
					_, err := f.Face.Write(req.Payload)
					if err != nil {
						log.Println(xcolor.Red("PUB"), request)
						break
					}
					continue
				}
				log.Println(xcolor.Red("SE.ERR"), request)
				continue
			case mqtt.Connack:
				//连接成功
				sub := mqtt.NewControlPacket(mqtt.Subscribe).(*mqtt.SubscribePacket)
				sub.Topics = []string{f.BroadcastTopic, f.EthTopic}
				sub.Qoss = []byte{0, 0}
				if err := cli.Send(sub); err != nil {
					log.Println(xcolor.Red("SUB"), request)
					break
				}
				continue
			}
		}
		return nil
	}
	panic(app.Run(os.Args))
}
