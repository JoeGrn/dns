package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type DNSHeader struct {
	ID      uint16
	Flags   uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

func ReadDNSHeader(reader *bytes.Reader) (DNSHeader, error) {
	var header DNSHeader
	err := binary.Read(reader, binary.BigEndian, &header)
	if err != nil {
		return header, fmt.Errorf("failed to read DNS header: %v", err)
	}
	return header, nil
}

func WriteDNSHeader(buffer *bytes.Buffer, header DNSHeader) error {
	err := binary.Write(buffer, binary.BigEndian, header)
	if err != nil {
		return fmt.Errorf("failed to write DNS header: %v", err)
	}
	return nil
}
