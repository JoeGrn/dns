package tests

import (
	"net"
	"reflect"
	"testing"

	dns "github.com/joegrn/dns/pkg"
)

func TestHandleDnsRequest(t *testing.T) {
	tests := []struct {
		name    string
		request []byte
		want    []byte
	}{
		{
			name: "Standard A record query",
			request: []byte{
				0x12, 0x34, // ID: 0x1234
				0x01, 0x00, // Flags: standard query
				0x00, 0x01, // QDCOUNT: 1
				0x00, 0x00, // ANCOUNT: 0
				0x00, 0x00, // NSCOUNT: 0
				0x00, 0x00, // ARCOUNT: 0
				3, 'w', 'w', 'w', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0, // www.example.com
				0x00, 0x01, // QTYPE = A
				0x00, 0x01, // QCLASS = IN
			},
			// The response should be the request + an answer section with 8.8.8.8
			want: append(
				[]byte{
					0x12, 0x34, // ID: 0x1234 (preserved from request)
					0x01, 0x00, // Flags: copied from request
					0x00, 0x01, // QDCOUNT: 1
					0x00, 0x01, // ANCOUNT: 1
					0x00, 0x00, // NSCOUNT: 0
					0x00, 0x00, // ARCOUNT: 0
					3, 'w', 'w', 'w', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0,
					0x00, 0x01, // QTYPE = A
					0x00, 0x01, // QCLASS = IN
				},
				// Answer section
				[]byte{
					3, 'w', 'w', 'w', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0,
					0x00, 0x01, // TYPE = A
					0x00, 0x01, // CLASS = IN
					0x00, 0x00, 0x00, 0x3c, // TTL = 60
					0x00, 0x04, // RDLENGTH = 4
					0x08, 0x08, 0x08, 0x08, // RDATA (8.8.8.8)
				}...,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a real UDP connection for the test
			// since our function requires *net.UDPConn specifically
			conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
			if err != nil {
				t.Fatalf("Failed to create UDP conn: %v", err)
			}
			defer conn.Close()

			sourceAddr := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 53000}

			got := dns.HandleDnsRequest(conn, sourceAddr, tt.request)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleDnsRequest() returned incorrect response")
				t.Errorf("Got:  %v", got)
				t.Errorf("Want: %v", tt.want)

				if len(got) != len(tt.want) {
					t.Errorf("Response length differs: got %d bytes, want %d bytes", len(got), len(tt.want))
				} else {
					for i := range got {
						if got[i] != tt.want[i] {
							t.Errorf("Difference at byte %d: got 0x%02x, want 0x%02x", i, got[i], tt.want[i])
						}
					}
				}
			}
		})
	}
}
