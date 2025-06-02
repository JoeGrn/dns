package tests

import (
	"bytes"
	"reflect"
	"testing"

	dns "github.com/joegrn/dns/pkg"
)

func TestReadDNSHeader(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    dns.DNSHeader
		wantErr bool
	}{
		{
			name: "Valid header",
			input: []byte{
				0x12, 0x34, // ID: 0x1234
				0x01, 0x00, // Flags: 0x0100
				0x00, 0x01, // QDCount: 1
				0x00, 0x00, // ANCount: 0
				0x00, 0x00, // NSCount: 0
				0x00, 0x00, // ARCount: 0
			},
			want: dns.DNSHeader{
				ID:      0x1234,
				Flags:   0x0100,
				QDCount: 1,
				ANCount: 0,
				NSCount: 0,
				ARCount: 0,
			},
			wantErr: false,
		},
		{
			name: "Another valid header",
			input: []byte{
				0xAB, 0xCD, // ID: 0xABCD
				0x81, 0x80, // Flags: 0x8180 (standard response)
				0x00, 0x01, // QDCount: 1
				0x00, 0x01, // ANCount: 1
				0x00, 0x02, // NSCount: 2
				0x00, 0x03, // ARCount: 3
			},
			want: dns.DNSHeader{
				ID:      0xABCD,
				Flags:   0x8180,
				QDCount: 1,
				ANCount: 1,
				NSCount: 2,
				ARCount: 3,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.input)
			got, err := dns.ReadDNSHeader(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadDNSHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadDNSHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteDNSHeader(t *testing.T) {
	tests := []struct {
		name    string
		header  dns.DNSHeader
		want    []byte
		wantErr bool
	}{
		{
			name: "Standard header",
			header: dns.DNSHeader{
				ID:      0x1234,
				Flags:   0x0100,
				QDCount: 1,
				ANCount: 0,
				NSCount: 0,
				ARCount: 0,
			},
			want: []byte{
				0x12, 0x34, // ID: 0x1234
				0x01, 0x00, // Flags: 0x0100
				0x00, 0x01, // QDCount: 1
				0x00, 0x00, // ANCount: 0
				0x00, 0x00, // NSCount: 0
				0x00, 0x00, // ARCount: 0
			},
			wantErr: false,
		},
		{
			name: "Response header",
			header: dns.DNSHeader{
				ID:      0xABCD,
				Flags:   0x8180,
				QDCount: 1,
				ANCount: 1,
				NSCount: 2,
				ARCount: 3,
			},
			want: []byte{
				0xAB, 0xCD, // ID: 0xABCD
				0x81, 0x80, // Flags: 0x8180
				0x00, 0x01, // QDCount: 1
				0x00, 0x01, // ANCount: 1
				0x00, 0x02, // NSCount: 2
				0x00, 0x03, // ARCount: 3
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := new(bytes.Buffer)
			err := dns.WriteDNSHeader(buffer, tt.header)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteDNSHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(buffer.Bytes(), tt.want) {
				t.Errorf("WriteDNSHeader() = %v, want %v", buffer.Bytes(), tt.want)
			}
		})
	}
}
