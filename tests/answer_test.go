package tests

import (
	"bytes"
	"reflect"
	"testing"
	
	dns "github.com/joegrn/dns/pkg"
)

func TestWriteDNSAnswer(t *testing.T) {
	tests := []struct {
		name    string
		answer  dns.DNSAnswer
		want    []byte
		wantErr bool
	}{
		{
			name: "A record answer",
			answer: dns.DNSAnswer{
				Name:     []byte{3, 'w', 'w', 'w', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0},
				Type:     1,                      // A record
				Class:    1,                      // IN class
				TTL:      300,                    // 5 minutes
				RDLength: 4,                      // IPv4 address is 4 bytes
				RData:    []byte{192, 168, 1, 1}, // 192.168.1.1
			},
			want: []byte{
				3, 'w', 'w', 'w', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0, // www.example.com
				0x00, 0x01, // TYPE = A
				0x00, 0x01, // CLASS = IN
				0x00, 0x00, 0x01, 0x2c, // TTL = 300
				0x00, 0x04, // RDLENGTH = 4
				192, 168, 1, 1, // RDATA (IP address)
			},
			wantErr: false,
		},
		{
			name: "AAAA record answer",
			answer: dns.DNSAnswer{
				Name:     []byte{3, 'w', 'w', 'w', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0},
				Type:     28,                                                                                                     // AAAA record
				Class:    1,                                                                                                      // IN class
				TTL:      3600,                                                                                                   // 1 hour
				RDLength: 16,                                                                                                     // IPv6 address is 16 bytes
				RData:    []byte{0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, // 2001:db8::1
			},
			want: []byte{
				3, 'w', 'w', 'w', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0, // www.example.com
				0x00, 0x1c, // TYPE = AAAA
				0x00, 0x01, // CLASS = IN
				0x00, 0x00, 0x0e, 0x10, // TTL = 3600
				0x00, 0x10, // RDLENGTH = 16
				0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // RDATA (IPv6 address)
			},
			wantErr: false,
		},
		{
			name: "Empty answer",
			answer: dns.DNSAnswer{
				Name:     []byte{0}, // root domain
				Type:     1,         // A record
				Class:    1,         // IN class
				TTL:      0,         // no caching
				RDLength: 0,         // no data
				RData:    []byte{},
			},
			want: []byte{
				0,          // root domain
				0x00, 0x01, // TYPE = A
				0x00, 0x01, // CLASS = IN
				0x00, 0x00, 0x00, 0x00, // TTL = 0
				0x00, 0x00, // RDLENGTH = 0
				// no RDATA
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := new(bytes.Buffer)
			err := dns.WriteDNSAnswer(buffer, tt.answer)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteDNSAnswer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(buffer.Bytes(), tt.want) {
				t.Errorf("WriteDNSAnswer() = %v, want %v", buffer.Bytes(), tt.want)
			}
		})
	}
}
