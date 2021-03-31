package main

import (
	"log"
	"net"
)

type UdpConn struct {
	outConn net.Conn
	inConn  *net.UDPConn
}

func CreateUdpConn() *UdpConn {
	OutConn, err := net.Dial("udp", Config.ToAddr)
	if err != nil {
		log.Fatalf("Dial udp failed, err=%+v", err)
	}
	localConn, err := net.ListenUDP("udp",
		&net.UDPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: Config.SelfPort,
		})
	if err != nil {
		log.Fatalf("Listen udp failed, err=%+v", err)
	}
	return &UdpConn{
		inConn:  localConn,
		outConn: OutConn,
	}
}

func (u *UdpConn) Read(buffer []byte) []byte {
	n, err := u.inConn.Read(buffer)
	if err != nil {
		logger.Warnf("Read failed, err=%+v", err)
		return nil
	}
	return buffer[:n]
}

func (u *UdpConn) Write(data []byte) {
	n, err := u.outConn.Write((data))
	if err != nil {
		logger.Warnf("Write to udp failed, err=%+v", err)
	}
	if n != len(data) {
		logger.Warnf("Cannot write all data to udp")
	}
}
