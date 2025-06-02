package tests

import (
	"bytes"
	"reflect"
	"testing"
	
	dns "github.com/joegrn/dns/pkg"
)

func TestReadDNSQuestion(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    dns.DNSQuestion
		wantErr bool
	}{
		{
			name: "Simple domain",
			input: []byte{
				3, 'w', 'w', 'w', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0, // www.example.com
				0x00, 0x01, // QTYPE = A
				0x00, 0x01, // QCLASS = IN
			},
			want: dns.DNSQuestion{
				QName:  []byte{3, 'w', 'w', 'w', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0},
				QType:  1, // A record
				QClass: 1, // IN class
			},
			wantErr: false,
		},
		{
			name: "Another domain with different type",
			input: []byte{
				6, 'g', 'o', 'o', 'g', 'l', 'e', 3, 'c', 'o', 'm', 0, // google.com
				0x00, 0x0f, // QTYPE = MX
				0x00, 0x01, // QCLASS = IN
			},
			want: dns.DNSQuestion{
				QName:  []byte{6, 'g', 'o', 'o', 'g', 'l', 'e', 3, 'c', 'o', 'm', 0},
				QType:  15, // MX record
				QClass: 1,  // IN class
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.input)
			got, err := dns.ReadDNSQuestion(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadDNSQuestion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got.QName, tt.want.QName) {
				t.Errorf("ReadDNSQuestion().QName = %v, want %v", got.QName, tt.want.QName)
			}

			if got.QType != tt.want.QType {
				t.Errorf("ReadDNSQuestion().QType = %v, want %v", got.QType, tt.want.QType)
			}

			if got.QClass != tt.want.QClass {
				t.Errorf("ReadDNSQuestion().QClass = %v, want %v", got.QClass, tt.want.QClass)
			}
		})
	}
}

func TestWriteDNSQuestion(t *testing.T) {
	tests := []struct {
		name     string
		question dns.DNSQuestion
		want     []byte
		wantErr  bool
	}{
		{
			name: "Simple A record question",
			question: dns.DNSQuestion{
				QName:  []byte{3, 'w', 'w', 'w', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0},
				QType:  1, // A record
				QClass: 1, // IN class
			},
			want: []byte{
				3, 'w', 'w', 'w', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0, // www.example.com
				0x00, 0x01, // QTYPE = A
				0x00, 0x01, // QCLASS = IN
			},
			wantErr: false,
		},
		{
			name: "MX record question",
			question: dns.DNSQuestion{
				QName:  []byte{6, 'g', 'o', 'o', 'g', 'l', 'e', 3, 'c', 'o', 'm', 0},
				QType:  15, // MX record
				QClass: 1,  // IN class
			},
			want: []byte{
				6, 'g', 'o', 'o', 'g', 'l', 'e', 3, 'c', 'o', 'm', 0, // google.com
				0x00, 0x0f, // QTYPE = MX
				0x00, 0x01, // QCLASS = IN
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := new(bytes.Buffer)
			err := dns.WriteDNSQuestion(buffer, tt.question)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteDNSQuestion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(buffer.Bytes(), tt.want) {
				t.Errorf("WriteDNSQuestion() = %v, want %v", buffer.Bytes(), tt.want)
			}
		})
	}
}
