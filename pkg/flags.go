package dns

const (
	QR_BIT     = 15
	OPCODE_POS = 11
	AA_BIT     = 10
	TC_BIT     = 9
	RD_BIT     = 8
	RA_BIT     = 7
	Z_POS      = 4
	RCODE_POS  = 0
)

const (
	OPCODE_MASK = 0x7800
	Z_MASK      = 0x0070
	RCODE_MASK  = 0x000F
)

type Flags struct {
	QR     bool
	OpCode uint8
	AA     bool
	TC     bool
	RD     bool
	RA     bool
	Z      uint8
	RCODE  uint8
}

func getBit(flags uint16, pos int) bool {
	return (flags & (1 << pos)) != 0
}

func setBit(flags uint16, pos int, value bool) uint16 {
	if value {
		return flags | (1 << pos)
	}
	return flags & ^(1 << pos)
}

func getField(flags uint16, mask uint16, pos int) uint8 {
	return uint8((flags & mask) >> pos)
}

func setField(flags uint16, value uint8, mask uint16, pos int) uint16 {
	return (flags & ^mask) | (uint16(value) << pos)
}

func UnmarshalFlags(buffer []byte) Flags {
	flags := uint16(buffer[0])<<8 | uint16(buffer[1])

	opcode := getField(flags, OPCODE_MASK, OPCODE_POS)
	var rcode uint8
	if opcode != 0 {
		rcode = 4
	}

	return Flags{
		QR:     getBit(flags, QR_BIT),
		OpCode: opcode,
		AA:     getBit(flags, AA_BIT),
		TC:     getBit(flags, TC_BIT),
		RD:     getBit(flags, RD_BIT),
		RA:     getBit(flags, RA_BIT),
		Z:      getField(flags, Z_MASK, Z_POS),
		RCODE:  rcode,
	}
}

func MarshalFlags(f Flags) uint16 {
	var flags uint16 = 0

	flags = setBit(flags, QR_BIT, f.QR)
	flags = setField(flags, f.OpCode, OPCODE_MASK, OPCODE_POS)
	flags = setBit(flags, AA_BIT, f.AA)
	flags = setBit(flags, TC_BIT, f.TC)
	flags = setBit(flags, RD_BIT, f.RD)
	flags = setBit(flags, RA_BIT, f.RA)
	flags = setField(flags, f.Z, Z_MASK, Z_POS)
	flags = setField(flags, f.RCODE, RCODE_MASK, RCODE_POS)

	return flags
}
