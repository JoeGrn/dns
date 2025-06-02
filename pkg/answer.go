package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type DNSAnswer struct {
	Name     []byte
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLength uint16
	RData    []byte
}

func WriteDNSAnswer(buffer *bytes.Buffer, answer DNSAnswer) error {
	if _, err := buffer.Write(answer.Name); err != nil {
		return fmt.Errorf("failed to write answer name: %v", err)
	}

	if err := binary.Write(buffer, binary.BigEndian, answer.Type); err != nil {
		return fmt.Errorf("failed to write answer type: %v", err)
	}

	if err := binary.Write(buffer, binary.BigEndian, answer.Class); err != nil {
		return fmt.Errorf("failed to write answer class: %v", err)
	}

	if err := binary.Write(buffer, binary.BigEndian, answer.TTL); err != nil {
		return fmt.Errorf("failed to write TTL: %v", err)
	}

	if err := binary.Write(buffer, binary.BigEndian, answer.RDLength); err != nil {
		return fmt.Errorf("failed to write RDLength: %v", err)
	}

	if _, err := buffer.Write(answer.RData); err != nil {
		return fmt.Errorf("failed to write RData: %v", err)
	}

	return nil
}
