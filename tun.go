package main

import (
	"github.com/songgao/water"
)

func ifceOuter(ifce *water.Interface, ch <-chan []byte) {
	for {
		data := <-ch
		ifce.Write(data)
	}
}

func ifceInner(ifce *water.Interface, ch chan<- []byte) {
	buffer := make([]byte, 65536)
	for {
		n, err := ifce.Read(buffer)
		if err != nil {
			logger.Warnf("Read from ifce failed, err=%+v", err)
		}
		data := make([]byte, n)
		copy(data, buffer)
		ch <- data
	}
}

func udpOuter(udp *UdpConn, ch <-chan []byte) {
	for {
		data := <-ch
		udp.Write(data)
	}
}

func udpInner(udp *UdpConn, ch chan<- []byte) {
	buffer := make([]byte, 65536)
	for {
		data := udp.Read(buffer)
		if data == nil {
			continue
		}
		cp := make([]byte, len(data))
		copy(cp, data)
		ch <- cp
	}
}

func Reader(ifce *water.Interface, udp *UdpConn) {
	buffer := make([]byte, 65536)
	for {
		n, err := ifce.Read(buffer)
		if err != nil {
			logger.Warnf("Read from tun failed, err=%+v", err)
			continue
		}
		udp.Write((buffer[:n]))
	}
}

func Writer(ifce *water.Interface, udp *UdpConn) {
	buffer := make([]byte, 655336)
	for {
		data := udp.Read(buffer)
		if data == nil {
			continue
		}
		n, err := ifce.Write(data)
		if err != nil {
			logger.Warnf("Write to tun failed, err=%+v", err)
			continue
		}
		if n != len(data) {
			logger.Warnf("Cannot write all data to tun")
		}
	}
}

func main() {
	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("Interface Name: %s\n", ifce.Name())
	udpConn := CreateUdpConn()
	inCh := make(chan []byte, 65536)
	outCh := make(chan []byte, 65536)
	go ifceOuter(ifce, outCh)
	go ifceInner(ifce, inCh)
	go udpOuter(udpConn, inCh)
	udpInner(udpConn, outCh)
}
