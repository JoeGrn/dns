package dns

import (
	"bytes"
	"fmt"
	"net"
)

func HandleDnsRequest(udpConn *net.UDPConn, source *net.UDPAddr, requestBuffer []byte) []byte {
	reader := bytes.NewReader(requestBuffer)

	queryHeader, err := ReadDNSHeader(reader)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	reader.Seek(12, 0)

	queryQuestion, err := ReadDNSQuestion(reader)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var responseBuffer = new(bytes.Buffer)

	responseHeader := DNSHeader{
		ID:      queryHeader.ID,
		Flags:   MarshalFlags(UnmarshalFlags(requestBuffer[2:4])),
		QDCount: 1,
		ANCount: 1,
		NSCount: 0,
		ARCount: 0,
	}

	if err := WriteDNSHeader(responseBuffer, responseHeader); err != nil {
		fmt.Println(err)
		return nil
	}

	if err := WriteDNSQuestion(responseBuffer, queryQuestion); err != nil {
		fmt.Println(err)
		return nil
	}

	answer := DNSAnswer{
		Name:     queryQuestion.QName,
		Type:     1,
		Class:    1,
		TTL:      60,
		RDLength: 4,
		RData:    []byte{0x08, 0x08, 0x08, 0x08},
	}

	if err := WriteDNSAnswer(responseBuffer, answer); err != nil {
		fmt.Println(err)
		return nil
	}

	return responseBuffer.Bytes()
}
