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
