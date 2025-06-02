package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type DNSQuestion struct {
	QName  []byte
	QType  uint16
	QClass uint16
}

func ReadDNSQuestion(reader *bytes.Reader) (DNSQuestion, error) {
	var question DNSQuestion
	var err error

	question.QName, err = readDomainName(reader)
	if err != nil {
		return question, fmt.Errorf("failed to read domain name: %v", err)
	}

	err = binary.Read(reader, binary.BigEndian, &question.QType)
	if err != nil {
		return question, fmt.Errorf("failed to read QTYPE: %v", err)
	}

	err = binary.Read(reader, binary.BigEndian, &question.QClass)
	if err != nil {
		return question, fmt.Errorf("failed to read QCLASS: %v", err)
	}

	return question, nil
}

func WriteDNSQuestion(buffer *bytes.Buffer, question DNSQuestion) error {
	if _, err := buffer.Write(question.QName); err != nil {
		return fmt.Errorf("failed to write QName: %v", err)
	}

	if err := binary.Write(buffer, binary.BigEndian, question.QType); err != nil {
		return fmt.Errorf("failed to write QType: %v", err)
	}

	if err := binary.Write(buffer, binary.BigEndian, question.QClass); err != nil {
		return fmt.Errorf("failed to write QClass: %v", err)
	}

	return nil
}

func readDomainName(reader *bytes.Reader) ([]byte, error) {
	var name []byte
	for {
		length, err := reader.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("failed to read label length: %v", err)
		}

		if length == 0 {
			break
		}

		if length&0xC0 == 0xC0 {

		}

		label := make([]byte, length)
		if _, err := reader.Read(label); err != nil {
			return nil, fmt.Errorf("failed to read label: %v", err)
		}

		name = append(name, length)
		name = append(name, label...)
	}
	name = append(name, 0)
	return name, nil
}
