package main

const (
	ADDR_IMPL = iota
	ADDR_A
	ADDR_IMMEDIATE
	ADDR_ZPG
	ADDR_ZPGX
	ADDR_ZPGY
	ADDR_ABS
	ADDR_ABSX
	ADDR_ABSY
	ADDR_REL
	ADDR_XIND
	ADDR_INDY
	ADDR_IND
)

var colors = []*RGB{
	{0x80, 0x80, 0x80}, {0x00, 0x3D, 0xA6}, {0x00, 0x12, 0xB0}, {0x44, 0x00, 0x96},
	{0xA1, 0x00, 0x5E}, {0xC7, 0x00, 0x28}, {0xBA, 0x06, 0x00}, {0x8C, 0x17, 0x00},
	{0x5C, 0x2F, 0x00}, {0x10, 0x45, 0x00}, {0x05, 0x4A, 0x00}, {0x00, 0x47, 0x2E},
	{0x00, 0x41, 0x66}, {0x00, 0x00, 0x00}, {0x05, 0x05, 0x05}, {0x05, 0x05, 0x05},
	{0xC7, 0xC7, 0xC7}, {0x00, 0x77, 0xFF}, {0x21, 0x55, 0xFF}, {0x82, 0x37, 0xFA},
	{0xEB, 0x2F, 0xB5}, {0xFF, 0x29, 0x50}, {0xFF, 0x22, 0x00}, {0xD6, 0x32, 0x00},
	{0xC4, 0x62, 0x00}, {0x35, 0x80, 0x00}, {0x05, 0x8F, 0x00}, {0x00, 0x8A, 0x55},
	{0x00, 0x99, 0xCC}, {0x21, 0x21, 0x21}, {0x09, 0x09, 0x09}, {0x09, 0x09, 0x09},
	{0xFF, 0xFF, 0xFF}, {0x0F, 0xD7, 0xFF}, {0x69, 0xA2, 0xFF}, {0xD4, 0x80, 0xFF},
	{0xFF, 0x45, 0xF3}, {0xFF, 0x61, 0x8B}, {0xFF, 0x88, 0x33}, {0xFF, 0x9C, 0x12},
	{0xFA, 0xBC, 0x20}, {0x9F, 0xE3, 0x0E}, {0x2B, 0xF0, 0x35}, {0x0C, 0xF0, 0xA4},
	{0x05, 0xFB, 0xFF}, {0x5E, 0x5E, 0x5E}, {0x0D, 0x0D, 0x0D}, {0x0D, 0x0D, 0x0D},
	{0xFF, 0xFF, 0xFF}, {0xA6, 0xFC, 0xFF}, {0xB3, 0xEC, 0xFF}, {0xDA, 0xAB, 0xEB},
	{0xFF, 0xA8, 0xF9}, {0xFF, 0xAB, 0xB3}, {0xFF, 0xD2, 0xB0}, {0xFF, 0xEF, 0xA6},
	{0xFF, 0xF7, 0x9C}, {0xD7, 0xE8, 0x95}, {0xA6, 0xED, 0xAF}, {0xA2, 0xF2, 0xDA},
	{0x99, 0xFF, 0xFC}, {0xDD, 0xDD, 0xDD}, {0x11, 0x11, 0x11}, {0x11, 0x11, 0x11},
}

var cycles = []int{
	/*0x00*/ 7, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 4, 4, 6, 6,
	/*0x10*/ 2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 6, 7,
	/*0x20*/ 6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 4, 4, 6, 6,
	/*0x30*/ 2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 6, 7,
	/*0x40*/ 6, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 3, 4, 6, 6,
	/*0x50*/ 2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 6, 7,
	/*0x60*/ 6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 5, 4, 6, 6,
	/*0x70*/ 2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 6, 7,
	/*0x80*/ 2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	/*0x90*/ 2, 6, 2, 6, 4, 4, 4, 4, 2, 4, 2, 5, 5, 4, 5, 5,
	/*0xA0*/ 2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	/*0xB0*/ 2, 5, 2, 5, 4, 4, 4, 4, 2, 4, 2, 4, 4, 4, 4, 4,
	/*0xC0*/ 2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	/*0xD0*/ 2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	/*0xE0*/ 2, 6, 3, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	/*0xF0*/ 2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
}

var opCodeList = map[int]*OpCode{
	0x00: {
		Base: "BRK",
		Mode: ADDR_IMPL,
	},
	0x01: {
		Base: "ORA",
		Mode: ADDR_XIND,
	},
	0x05: {
		Base: "ORA",
		Mode: ADDR_ZPG,
	},
	0x06: {
		Base: "ASL",
		Mode: ADDR_ZPG,
	},
	0x08: {
		Base: "PHP",
		Mode: ADDR_IMPL,
	},
	0x09: {
		Base: "ORA",
		Mode: ADDR_IMMEDIATE,
	},
	0x0A: {
		Base: "ASL",
		Mode: ADDR_A,
	},
	0x0D: {
		Base: "ORA",
		Mode: ADDR_ABS,
	},
	0x0E: {
		Base: "ASL",
		Mode: ADDR_ABS,
	},
	0x10: {
		Base: "BPL",
		Mode: ADDR_REL,
	},
	0x11: {
		Base: "ORA",
		Mode: ADDR_INDY,
	},
	0x15: {
		Base: "ORA",
		Mode: ADDR_ZPG,
	},
	0x16: {
		Base: "ASL",
		Mode: ADDR_ZPG,
	},
	0x18: {
		Base: "CLC",
		Mode: ADDR_IMPL,
	},
	0x19: {
		Base: "ORA",
		Mode: ADDR_ABSY,
	},
	0x1D: {
		Base: "ORA",
		Mode: ADDR_ABSX,
	},
	0x1E: {
		Base: "ASL",
		Mode: ADDR_ABSX,
	},
	0x20: {
		Base: "JSR",
		Mode: ADDR_ABS,
	},
	0x21: {
		Base: "AND",
		Mode: ADDR_XIND,
	},
	0x24: {
		Base: "BIT",
		Mode: ADDR_ZPG,
	},
	0x25: {
		Base: "AND",
		Mode: ADDR_ZPG,
	},
	0x26: {
		Base: "ROL",
		Mode: ADDR_ZPG,
	},
	0x28: {
		Base: "PLP",
		Mode: ADDR_IMPL,
	},
	0x29: {
		Base: "AND",
		Mode: ADDR_IMMEDIATE,
	},
	0x2A: {
		Base: "ROL",
		Mode: ADDR_A,
	},
	0x2C: {
		Base: "BIT",
		Mode: ADDR_ABS,
	},
	0x2D: {
		Base: "AND",
		Mode: ADDR_ABS,
	},
	0x2E: {
		Base: "ROL",
		Mode: ADDR_ABS,
	},
	0x30: {
		Base: "BMI",
		Mode: ADDR_REL,
	},
	0x31: {
		Base: "AND",
		Mode: ADDR_INDY,
	},
	0x35: {
		Base: "AND",
		Mode: ADDR_ZPGX,
	},
	0x36: {
		Base: "ROL",
		Mode: ADDR_ZPGX,
	},
	0x38: {
		Base: "SEC",
		Mode: ADDR_IMPL,
	},
	0x39: {
		Base: "AND",
		Mode: ADDR_ABSY,
	},
	0x3D: {
		Base: "AND",
		Mode: ADDR_ABSX,
	},
	0x3E: {
		Base: "ROL",
		Mode: ADDR_ABSX,
	},
	0x40: {
		Base: "RTI",
		Mode: ADDR_IMPL,
	},
	0x41: {
		Base: "EOR",
		Mode: ADDR_XIND,
	},
	0x45: {
		Base: "EOR",
		Mode: ADDR_ZPG,
	},
	0x46: {
		Base: "LSR",
		Mode: ADDR_ZPG,
	},
	0x48: {
		Base: "PHA",
		Mode: ADDR_IMPL,
	},
	0x49: {
		Base: "EOR",
		Mode: ADDR_IMMEDIATE,
	},
	0x4A: {
		Base: "LSR",
		Mode: ADDR_A,
	},
	0x4C: {
		Base: "JMP",
		Mode: ADDR_ABS,
	},
	0x4D: {
		Base: "EOR",
		Mode: ADDR_ABS,
	},
	0x4E: {
		Base: "LSR",
		Mode: ADDR_ABS,
	},
	0x50: {
		Base: "BVC",
		Mode: ADDR_REL,
	},
	0x51: {
		Base: "EOR",
		Mode: ADDR_INDY,
	},
	0x55: {
		Base: "EOR",
		Mode: ADDR_ZPGX,
	},
	0x56: {
		Base: "LSR",
		Mode: ADDR_ZPGX,
	},
	0x58: {
		Base: "CLI",
		Mode: ADDR_IMPL,
	},
	0x59: {
		Base: "EOR",
		Mode: ADDR_ABSY,
	},
	0x5D: {
		Base: "EOR",
		Mode: ADDR_ABSX,
	},
	0x5E: {
		Base: "LSR",
		Mode: ADDR_ABSX,
	},
	0x60: {
		Base: "RTS",
		Mode: ADDR_IMPL,
	},
	0x61: {
		Base: "ADC",
		Mode: ADDR_XIND,
	},
	0x65: {
		Base: "ADC",
		Mode: ADDR_ZPG,
	},
	0x66: {
		Base: "ROR",
		Mode: ADDR_ZPG,
	},
	0x68: {
		Base: "PLA",
		Mode: ADDR_IMPL,
	},
	0x69: {
		Base: "ADC",
		Mode: ADDR_IMMEDIATE,
	},
	0x6A: {
		Base: "ROR",
		Mode: ADDR_A,
	},
	0x6C: {
		Base: "JMP",
		Mode: ADDR_IND,
	},
	0x6D: {
		Base: "ADC",
		Mode: ADDR_ABS,
	},
	0x6E: {
		Base: "ROR",
		Mode: ADDR_ABS,
	},
	0x70: {
		Base: "BVS",
		Mode: ADDR_REL,
	},
	0x71: {
		Base: "ADC",
		Mode: ADDR_INDY,
	},
	0x75: {
		Base: "ADC",
		Mode: ADDR_ZPGX,
	},
	0x76: {
		Base: "ROR",
		Mode: ADDR_ZPGX,
	},
	0x78: {
		Base: "SEI",
		Mode: ADDR_IMPL,
	},
	0x79: {
		Base: "ADC",
		Mode: ADDR_ABSY,
	},
	0x7D: {
		Base: "ADC",
		Mode: ADDR_ABSX,
	},
	0x7E: {
		Base: "ROR",
		Mode: ADDR_ABSX,
	},
	0x81: {
		Base: "STA",
		Mode: ADDR_XIND,
	},
	0x84: {
		Base: "STY",
		Mode: ADDR_ZPG,
	},
	0x85: {
		Base: "STA",
		Mode: ADDR_ZPG,
	},
	0x86: {
		Base: "STX",
		Mode: ADDR_ZPG,
	},
	0x88: {
		Base: "DEY",
		Mode: ADDR_IMPL,
	},
	0x8A: {
		Base: "TXA",
		Mode: ADDR_IMPL,
	},
	0x8C: {
		Base: "STY",
		Mode: ADDR_ABS,
	},
	0x8D: {
		Base: "STA",
		Mode: ADDR_ABS,
	},
	0x8E: {
		Base: "STX",
		Mode: ADDR_ABS,
	},
	0x90: {
		Base: "BCC",
		Mode: ADDR_REL,
	},
	0x91: {
		Base: "STA",
		Mode: ADDR_INDY,
	},
	0x94: {
		Base: "STY",
		Mode: ADDR_ZPGX,
	},
	0x95: {
		Base: "STA",
		Mode: ADDR_ZPGX,
	},
	0x96: {
		Base: "STX",
		Mode: ADDR_ZPGY,
	},
	0x98: {
		Base: "TYA",
		Mode: ADDR_IMPL,
	},
	0x99: {
		Base: "STA",
		Mode: ADDR_ABSY,
	},
	0x9A: {
		Base: "TXS",
		Mode: ADDR_IMPL,
	},
	0x9D: {
		Base: "STA",
		Mode: ADDR_ABSX,
	},
	0xA0: {
		Base: "LDY",
		Mode: ADDR_IMMEDIATE,
	},
	0xA1: {
		Base: "LDA",
		Mode: ADDR_XIND,
	},
	0xA2: {
		Base: "LDX",
		Mode: ADDR_IMMEDIATE,
	},
	0xA4: {
		Base: "LDY",
		Mode: ADDR_ZPG,
	},
	0xA5: {
		Base: "LDA",
		Mode: ADDR_ZPG,
	},
	0xA6: {
		Base: "LDX",
		Mode: ADDR_ZPG,
	},
	0xA8: {
		Base: "TAY",
		Mode: ADDR_IMPL,
	},
	0xA9: {
		Base: "LDA",
		Mode: ADDR_IMMEDIATE,
	},
	0xAA: {
		Base: "TAX",
		Mode: ADDR_IMPL,
	},
	0xAC: {
		Base: "LDY",
		Mode: ADDR_ABS,
	},
	0xAD: {
		Base: "LDA",
		Mode: ADDR_ABS,
	},
	0xAE: {
		Base: "LDX",
		Mode: ADDR_ABS,
	},
	0xB0: {
		Base: "BCS",
		Mode: ADDR_REL,
	},
	0xB1: {
		Base: "LDA",
		Mode: ADDR_INDY,
	},
	0xB4: {
		Base: "LDY",
		Mode: ADDR_ZPGX,
	},
	0xB5: {
		Base: "LDA",
		Mode: ADDR_ZPGX,
	},
	0xB6: {
		Base: "LDX",
		Mode: ADDR_ZPGY,
	},
	0xB8: {
		Base: "CLV",
		Mode: ADDR_IMPL,
	},
	0xB9: {
		Base: "LDA",
		Mode: ADDR_ABSY,
	},
	0xBA: {
		Base: "TSX",
		Mode: ADDR_IMPL,
	},
	0xBC: {
		Base: "LDY",
		Mode: ADDR_ABSX,
	},
	0xBD: {
		Base: "LDA",
		Mode: ADDR_ABSX,
	},
	0xBE: {
		Base: "LDX",
		Mode: ADDR_ABSY,
	},
	0xC0: {
		Base: "CPY",
		Mode: ADDR_IMMEDIATE,
	},
	0xC1: {
		Base: "CMP",
		Mode: ADDR_XIND,
	},
	0xC4: {
		Base: "CPY",
		Mode: ADDR_ZPG,
	},
	0xC5: {
		Base: "CMP",
		Mode: ADDR_ZPG,
	},
	0xC6: {
		Base: "DEC",
		Mode: ADDR_ZPG,
	},
	0xC8: {
		Base: "INY",
		Mode: ADDR_IMPL,
	},
	0xC9: {
		Base: "CMP",
		Mode: ADDR_IMMEDIATE,
	},
	0xCA: {
		Base: "DEX",
		Mode: ADDR_IMPL,
	},
	0xCC: {
		Base: "CPY",
		Mode: ADDR_ABS,
	},
	0xCD: {
		Base: "CMP",
		Mode: ADDR_ABS,
	},
	0xCE: {
		Base: "DEC",
		Mode: ADDR_ABS,
	},
	0xD0: {
		Base: "BNE",
		Mode: ADDR_REL,
	},
	0xD1: {
		Base: "CMP",
		Mode: ADDR_INDY,
	},
	0xD5: {
		Base: "CMP",
		Mode: ADDR_ZPGX,
	},
	0xD6: {
		Base: "DEC",
		Mode: ADDR_ZPGX,
	},
	0xD8: {
		Base: "CLD",
		Mode: ADDR_IMPL,
	},
	0xD9: {
		Base: "CMP",
		Mode: ADDR_ABSY,
	},
	0xDD: {
		Base: "CMP",
		Mode: ADDR_ABSX,
	},
	0xDE: {
		Base: "DEC",
		Mode: ADDR_ABSX,
	},
	0xE0: {
		Base: "CPX",
		Mode: ADDR_IMMEDIATE,
	},
	0xE1: {
		Base: "SBC",
		Mode: ADDR_XIND,
	},
	0xE4: {
		Base: "CPX",
		Mode: ADDR_ZPG,
	},
	0xE5: {
		Base: "SBC",
		Mode: ADDR_ZPG,
	},
	0xE6: {
		Base: "INC",
		Mode: ADDR_ZPG,
	},
	0xE8: {
		Base: "INX",
		Mode: ADDR_IMPL,
	},
	0xE9: {
		Base: "SBC",
		Mode: ADDR_IMMEDIATE,
	},
	0xEA: {
		Base: "NOP",
		Mode: ADDR_IMPL,
	},
	0xEC: {
		Base: "CPX",
		Mode: ADDR_ABS,
	},
	0xED: {
		Base: "SBC",
		Mode: ADDR_ABS,
	},
	0xEE: {
		Base: "INC",
		Mode: ADDR_ABS,
	},
	0xF0: {
		Base: "BEQ",
		Mode: ADDR_REL,
	},
	0xF1: {
		Base: "SBC",
		Mode: ADDR_INDY,
	},
	0xF5: {
		Base: "SBC",
		Mode: ADDR_ZPGX,
	},
	0xF6: {
		Base: "INC",
		Mode: ADDR_ZPGX,
	},
	0xF8: {
		Base: "SED",
		Mode: ADDR_IMPL,
	},
	0xF9: {
		Base: "SBC",
		Mode: ADDR_ABSY,
	},
	0xFD: {
		Base: "SBC",
		Mode: ADDR_ABSX,
	},
	0xFE: {
		Base: "INC",
		Mode: ADDR_ABSX,
	},
}

type OpCode struct {
	Base    string
	Mode    int
	Cycle   int
	Operand int
}

func (opCode *OpCode) FetchOperand(cpu *Cpu) {
	switch opCode.Mode {
	case ADDR_IMPL:
		return
	case ADDR_A:
		return
	case ADDR_IMMEDIATE:
		opCode.Operand = cpu.Fetch()
	case ADDR_ZPG:
		opCode.Operand = cpu.Fetch()
	case ADDR_ZPGX:
		opCode.Operand = (cpu.Fetch() + cpu.Register.X) & 0xFF
	case ADDR_ZPGY:
		opCode.Operand = (cpu.Fetch() + cpu.Register.Y) & 0xFF
	case ADDR_ABS:
		l := cpu.Fetch()
		h := cpu.Fetch()
		opCode.Operand = l + h*256
	case ADDR_ABSX:
		l := cpu.Fetch()
		h := cpu.Fetch()
		opCode.Operand = l + h*256 + cpu.Register.X
	case ADDR_ABSY:
		l := cpu.Fetch()
		h := cpu.Fetch()
		opCode.Operand = l + h*256 + cpu.Register.Y
	case ADDR_REL:
		rel := cpu.Fetch()
		if rel < 0x7F {
			opCode.Operand = cpu.Register.PC + rel
		} else {
			opCode.Operand = cpu.Register.PC - (rel ^ 0xFF) - 1
		}
	case ADDR_XIND:
		addr := (cpu.Read(cpu.Fetch()) + cpu.Register.X) & 0xFF
		opCode.Operand = cpu.Read(addr)+ cpu.Read(addr+1)*256
	case ADDR_INDY:
		addr := cpu.Fetch()
		opCode.Operand = (cpu.Read(addr) + cpu.Read((addr+1)&0xFF)<<7 + cpu.Register.Y)&0xFFFF
	case ADDR_IND:
		opCode.Operand = cpu.Read(cpu.Fetch()+cpu.Fetch()*256)
	}
}
