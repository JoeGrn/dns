package tests

import (
	"reflect"
	"testing"

	dns "github.com/joegrn/dns/pkg"
)

func TestUnmarshalFlags(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  dns.Flags
	}{
		{
			name:  "Query flags",
			input: []byte{0x00, 0x00}, // Standard query flags (all zeros)
			want: dns.Flags{
				QR:     false,
				OpCode: 0,
				AA:     false,
				TC:     false,
				RD:     false,
				RA:     false,
				Z:      0,
				RCODE:  0,
			},
		},
		{
			name:  "Response flags",
			input: []byte{0x81, 0x80}, // Standard response flags
			want: dns.Flags{
				QR:     true,
				OpCode: 0,
				AA:     false,
				TC:     false,
				RD:     true,
				RA:     true,
				Z:      0,
				RCODE:  0,
			},
		},
		{
			name:  "Authoritative answer",
			input: []byte{0x84, 0x00}, // Authoritative answer
			want: dns.Flags{
				QR:     true,
				OpCode: 0,
				AA:     true,
				TC:     false,
				RD:     false,
				RA:     false,
				Z:      0,
				RCODE:  0,
			},
		},
		{
			name:  "Truncated response",
			input: []byte{0x82, 0x00}, // Truncated response
			want: dns.Flags{
				QR:     true,
				OpCode: 0,
				AA:     false,
				TC:     true,
				RD:     false,
				RA:     false,
				Z:      0,
				RCODE:  0,
			},
		},
		{
			name:  "Name error (NXDOMAIN)",
			input: []byte{0x81, 0x83}, // Name error
			want: dns.Flags{
				QR:     true,
				OpCode: 0,
				AA:     false,
				TC:     false,
				RD:     true,
				RA:     true,
				Z:      0,
				RCODE:  3,
			},
		},
		{
			name:  "Non-standard opcode",
			input: []byte{0x28, 0x00}, // Opcode 5 (not standard)
			want: dns.Flags{
				QR:     false,
				OpCode: 5,
				AA:     false,
				TC:     false,
				RD:     false,
				RA:     false,
				Z:      0,
				RCODE:  4, // Per your implementation, non-0 opcode gets RCODE 4
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dns.UnmarshalFlags(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarshalFlags(t *testing.T) {
	tests := []struct {
		name  string
		flags dns.Flags
		want  uint16
	}{
		{
			name: "Standard query",
			flags: dns.Flags{
				QR:     false,
				OpCode: 0,
				AA:     false,
				TC:     false,
				RD:     true,
				RA:     false,
				Z:      0,
				RCODE:  0,
			},
			want: 0x0100, // Just RD bit set
		},
		{
			name: "Standard response",
			flags: dns.Flags{
				QR:     true,
				OpCode: 0,
				AA:     false,
				TC:     false,
				RD:     true,
				RA:     true,
				Z:      0,
				RCODE:  0,
			},
			want: 0x8180, // QR, RD, RA bits set
		},
		{
			name: "Authoritative answer",
			flags: dns.Flags{
				QR:     true,
				OpCode: 0,
				AA:     true,
				TC:     false,
				RD:     false,
				RA:     false,
				Z:      0,
				RCODE:  0,
			},
			want: 0x8400, // QR and AA bits set
		},
		{
			name: "Name error (NXDOMAIN)",
			flags: dns.Flags{
				QR:     true,
				OpCode: 0,
				AA:     false,
				TC:     false,
				RD:     true,
				RA:     true,
				Z:      0,
				RCODE:  3,
			},
			want: 0x8183, // QR, RD, RA bits set and RCODE=3
		},
		{
			name: "With Z bits set",
			flags: dns.Flags{
				QR:     false,
				OpCode: 0,
				AA:     false,
				TC:     false,
				RD:     false,
				RA:     false,
				Z:      3,
				RCODE:  0,
			},
			want: 0x0030, // Z=3 (bits 4-6 set)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dns.MarshalFlags(tt.flags)
			if got != tt.want {
				t.Errorf("MarshalFlags() = 0x%04x, want 0x%04x", got, tt.want)
			}
		})
	}
}

func TestBitManipulation(t *testing.T) {
	// For these tests, we need to directly access the bit manipulation functions,
	// but since they're unexported, we'll test them indirectly through the Marshal/Unmarshal functions
	
	t.Run("getBit", func(t *testing.T) {
		// Test if proper bits are extracted from flags
		flags := dns.Flags{QR: true, RD: true}
		marshaled := dns.MarshalFlags(flags)
		
		// Assert QR and RD bits are set in the marshaled value
		if marshaled&0x8000 == 0 {
			t.Error("QR bit not properly set")
		}
		if marshaled&0x0100 == 0 {
			t.Error("RD bit not properly set")
		}
	})
	
	t.Run("setBit", func(t *testing.T) {
		// Test by unmarshaling flags with specific bits set
		unmarshaled := dns.UnmarshalFlags([]byte{0x81, 0x00}) // QR and RD set
		
		if !unmarshaled.QR {
			t.Error("QR bit not correctly unmarshaled")
		}
		if !unmarshaled.RD {
			t.Error("RD bit not correctly unmarshaled")
		}
	})
	
	t.Run("getField", func(t *testing.T) {
		// Test opcode field extraction
		unmarshaled := dns.UnmarshalFlags([]byte{0x28, 0x00}) // Opcode 5
		
		if unmarshaled.OpCode != 5 {
			t.Errorf("OpCode not correctly extracted, got %d, want 5", unmarshaled.OpCode)
		}
	})
	
	t.Run("setField", func(t *testing.T) {
		// Test setting RCODE field
		flags := dns.Flags{RCODE: 3}
		marshaled := dns.MarshalFlags(flags)
		
		if marshaled&0x000F != 3 {
			t.Errorf("RCODE not correctly set, got %d, want 3", marshaled&0x000F)
		}
	})
}
