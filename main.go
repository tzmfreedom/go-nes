package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
)

type Register struct{
	A int
	X int
	Y int
	P *StatusRegister
	SP int
	PC int
}

type StatusRegister struct{
	Negative bool
	Overflow bool
	Reserved bool
	Break bool
	Decimal bool
	Interrupt bool
	Zero bool
	Carry bool
}

type Cpu struct {
	Register *Register
	PrgROM []byte
	ChrROM []byte
}

var images []*Image

type OpCode struct {
	Base string
	Mode int
	Cycle int
}

const (
	ADDR_IMPL = iota
	ADDR_A
	ADDR_SHARP
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

var opCodeList = map[byte]*OpCode{
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
		Mode: ADDR_SHARP,
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
		Mode: ADDR_SHARP,
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
		Mode: ADDR_SHARP,
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
		Mode: ADDR_SHARP,
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
		Mode: ADDR_SHARP,
	},
	0xA1: {
		Base: "LDA",
		Mode: ADDR_XIND,
	},
	0xA2: {
		Base: "LDX",
		Mode: ADDR_SHARP,
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
		Mode: ADDR_SHARP,
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
		Mode: ADDR_SHARP,
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
		Mode: ADDR_SHARP,
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
		Mode: ADDR_SHARP,
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
		Mode: ADDR_SHARP,
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

func main() {
	filename := os.Args[1]
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(b))
	fmt.Println(string(b[0]) + string(b[1]) + string(b[2]))
	prgSize := int(b[4])
	chrSize := int(b[5])
	prgRomEnd := 0x10 + prgSize * 0x4000
	prgRom := b[0x10:prgRomEnd]
	chrRom := b[prgRomEnd:prgRomEnd+chrSize*0x2000]
	fmt.Printf("PRG SIZE: %d => %d\n", prgSize, len(prgRom))
	fmt.Printf("CHR SIZE: %d => %d\n", chrSize, len(chrRom))

	cpu := &Cpu{
		Register: &Register{
			P: &StatusRegister{},
		},
		PrgROM: prgRom,
		ChrROM: chrRom,
	}

	for i := 0; i < 512; i++ {
		index := i*16
		NewImage(chrRom[index:index+16])
	}
	cpu.Reset()
}

func (cpu *Cpu) Reset() {
	f := cpu.PrgROM[0xFFFC - 0x8000]
	s := cpu.PrgROM[0xFFFD - 0x8000]
	cpu.Register.PC = int(s)*256+int(f) - 0x8000
}

func (cpu *Cpu) Fetch() byte {
	ret := cpu.PrgROM[cpu.Register.PC]
	cpu.Register.PC++
	return ret
}

func (cpu *Cpu) FetchOperand(opCode *OpCode) byte {
	return 0x00
}

func (cpu *Cpu) Run() {
	opCodeRaw := cpu.Fetch()
	opCode := opCodeList[opCodeRaw]
	opRand := cpu.FetchOperand(opCode)
	cpu.Execute(opCode, opRand)
}

func (cpu *Cpu) Execute(opCode *OpCode, opRand byte) {}

type Image struct {
	bitMap [][]int
}

func (img *Image) Render() {
	for _, bits := range img.bitMap {
		for _, b := range bits {
			if b == 0 {
				fmt.Print(" ")
			} else {
				fmt.Print("*")
			}
		}
		fmt.Println()
	}
}

const ImageSize = 8

func NewImage(src []byte) *Image {
	bitMap := make([][]int, ImageSize)
	for i := 0; i < ImageSize; i ++ {
		bit := make([]int, ImageSize)
		for j := 0; j < ImageSize; j++ {
			b := 0
			mul := int(math.Pow(2, float64(j)))
			if (int(src[i]) & mul) != 0 {
				b += 1
			}
			if (int(src[i+ImageSize]) & mul) != 0 {
				b += 2
			}
			bit[j] = b
		}
		bitMap[i] = bit
	}
	return &Image{
		bitMap: bitMap,
	}
}
