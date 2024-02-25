package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func DoBroadcast() {
	sock, err := net.DialUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 8830,
	}, &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 8829,
	})
	if err != nil {
		log.Fatal(err)
	}
	for {
		pingMessage := fmt.Sprintf("P:%d", broadcastPort)
		sock.Write([]byte(pingMessage))
		time.Sleep(45 * time.Second)
	}
}
