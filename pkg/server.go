package dns

import (
	"fmt"
	"net"
)

const (
	UDPMaxMessageSize uint = 512
)

func Serve() {

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	requestBuffer := make([]byte, UDPMaxMessageSize)

	for {
		size, source, err := udpConn.ReadFromUDP(requestBuffer)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		responseBuffer := HandleDnsRequest(udpConn, source, requestBuffer[:size])

		_, err = udpConn.WriteToUDP(responseBuffer, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
