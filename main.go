package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/k0kubun/pp"
	"image/color"
	"io/ioutil"
	"log"
	"math"
	"os"
)

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
	prgRomEnd := 0x10 + prgSize*0x4000
	prgRom := b[0x10:prgRomEnd]
	chrRom := b[prgRomEnd : prgRomEnd+chrSize*0x2000]
	fmt.Printf("PRG SIZE: %d => %d\n", prgSize, len(prgRom))
	fmt.Printf("CHR SIZE: %d => %d\n", chrSize, len(chrRom))

	cpu := &Cpu{
		RAM: make([]int, 0x0800),
		Register: &Register{
			P: &StatusRegister{},
		},
		PrgROM: prgRom,
		ChrROM: chrRom,
	}

	ppu := NewPPU()
	sprites := make([]*Sprite, 512)
	for i := 0; i < 512; i++ {
		index := i * 16
		sprites[i] = NewImage(chrRom[index : index+16])
	}
	ppu.sprites = sprites
	cpu.Reset()
	nes := &NES{
		cpu:   cpu,
		ppu:   ppu,
		cycle: 0,
	}
	if err := ebiten.Run(nes.update, 512, 480, 1, "NES"); err != nil {
		log.Fatal(err)
	}
}

type NES struct {
	cpu        *Cpu
	ppu        *PPU
	cycle      int
	background *BackGround
	pallet     *Pallet
}

func (nes *NES) update(screen *ebiten.Image) error {
	nes.cycle += nes.cpu.Run()
	background := nes.ppu.Run(nes.cycle * 3)
	if background != nil {
		nes.render(screen, background)
	} else if nes.background != nil {
		nes.render(screen, nes.background)
	}
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	return nil
}

func (nes *NES) render(screen *ebiten.Image, background *BackGround) {
	for i, line := range background.tiles {
		for j, tile := range line {
			nes.renderTile(screen, j, i, tile)
		}
	}
}

func (nes *NES) renderTile(screen *ebiten.Image, x, y int, tile *Tile) {
	for i, line := range tile.img.bitMap {
		for j, bit := range line {
			img, _ := ebiten.NewImage(1, 1, 0)
			c := colors[bit]
			err := img.Fill(color.RGBA{c.R, c.G, c.B, 0xff})
			if err != nil {
				panic(err)
			}
			options := &ebiten.DrawImageOptions{}
			options.GeoM.Translate(float64(x*8+j), float64(y*8+i))
			err = screen.DrawImage(img, options)
			if err != nil {
				panic(err)
			}
		}
	}
}

type Register struct {
	A  int
	X  int
	Y  int
	P  *StatusRegister
	SP int
	PC int
}

type PPU struct {
	RAM        []int
	cycle      int
	line       int
	background *BackGround
	sprites    []*Sprite
}

func NewPPU() *PPU {
	return &PPU{}
}

type BackGround struct {
	tiles [][]*Tile
}

func (b *BackGround) Add(x, y int, tile *Tile) {
	b.tiles[y][x] = tile
}

func NewBackGround() *BackGround {
	tiles := make([][]*Tile, 30)
	for i := 0; i < 30; i++ {
		tiles[i] = make([]*Tile, 32)
	}
	return &BackGround{
		tiles: tiles,
	}
}

type Pallet struct {
}

type Tile struct {
	img      *Sprite
	palletId int
}

func (ppu *PPU) Run(cycle int) *BackGround {
	ppu.cycle += cycle
	if ppu.cycle >= 341 {
		ppu.cycle -= 341
		ppu.line++
		if ppu.line < 240 && ppu.line%8 == 0 {
			ppu.BuildBackGround()
		}
		if ppu.line == 262 {
			ppu.line = 0
			return ppu.background
		}
	}
	return nil
}

func (ppu *PPU) BuildBackGround() {
	y := ppu.line / 8
	for x := 0; x < 32; x++ {
		tile := ppu.BuildTile(x, y)
		ppu.background.Add(x, y, tile)
	}
}

func (ppu *PPU) BuildTile(x, y int) *Tile {
	palletId := ppu.getPalletId(x, y)
	sprite := ppu.getSprite(x, y)
	return &Tile{
		img:      sprite,
		palletId: palletId,
	}
}

func (ppu *PPU) getPalletId(x, y int) int {
	tmpX := x / 2
	tmpY := y / 2
	palletBlock := ppu.RAM[tmpX+tmpY*8]

	var blockId uint
	cmpX := (tmpX / 2) % 2
	cmpY := (tmpY / 2) % 2
	if cmpX == 0 {
		if cmpY == 0 {
			blockId = 1
		} else {
			blockId = 3
		}
	} else {
		if cmpY == 0 {
			blockId = 2
		} else {
			blockId = 4
		}
	}
	return (palletBlock >> (blockId * 2)) & 0x11
}

func (ppu *PPU) getSprite(x, y int) *Sprite {
	return ppu.sprites[x+y*32]
}

type StatusRegister struct {
	Negative  bool
	Overflow  bool
	Reserved  bool
	Break     bool
	Decimal   bool
	Interrupt bool
	Zero      bool
	Carry     bool
}

func (r *StatusRegister) Int() int {
	return bool2int(r.Negative)*int(math.Pow(2, 7)) +
		bool2int(r.Overflow)*int(math.Pow(2, 6)) +
		bool2int(r.Reserved)*int(math.Pow(2, 5)) +
		bool2int(r.Break)*int(math.Pow(2, 4)) +
		bool2int(r.Decimal)*int(math.Pow(2, 3)) +
		bool2int(r.Interrupt)*int(math.Pow(2, 2)) +
		bool2int(r.Zero)*int(math.Pow(2, 1)) +
		bool2int(r.Carry)*int(math.Pow(2, 0))
}

func (r *StatusRegister) Set(v int) {
	r.Negative = v&int(math.Pow(2, 7)) != 0
	r.Overflow = v&int(math.Pow(2, 6)) != 0
	r.Reserved = v&int(math.Pow(2, 5)) != 0
	r.Break = v&int(math.Pow(2, 4)) != 0
	r.Decimal = v&int(math.Pow(2, 3)) != 0
	r.Interrupt = v&int(math.Pow(2, 2)) != 0
	r.Zero = v&int(math.Pow(2, 1)) != 0
	r.Carry = v&int(math.Pow(2, 0)) != 0
}

type Cpu struct {
	RAM      []int
	PPU      []int
	APU      []int
	Register *Register
	PrgROM   []byte
	ChrROM   []byte
}

var sprites []*Sprite

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
		opCode.Operand = cpu.Fetch() + cpu.Register.X
	case ADDR_ZPGY:
		opCode.Operand = cpu.Fetch() + cpu.Register.Y
	case ADDR_ABS:
		opCode.Operand = cpu.Fetch() + cpu.Fetch()*256
	case ADDR_ABSX:
		opCode.Operand = cpu.Fetch() + cpu.Fetch()*256 + cpu.Register.X
	case ADDR_ABSY:
		opCode.Operand = cpu.Fetch() + cpu.Fetch()*256 + cpu.Register.Y
	case ADDR_REL:
		opCode.Operand = cpu.Fetch() + cpu.Fetch()
	case ADDR_XIND:
		opCode.Operand = cpu.Read(cpu.Fetch()) + cpu.Register.X + cpu.Fetch()*256
	case ADDR_INDY:
		opCode.Operand = cpu.Read(cpu.Fetch()) + (cpu.Fetch()+cpu.Register.Y)*256
	case ADDR_IND:
		opCode.Operand = cpu.Read(cpu.Fetch()+cpu.Fetch()*256) + cpu.Fetch()*256
	}
}

type RGB struct {
	R uint8
	G uint8
	B uint8
}

func (cpu *Cpu) Write(index int, value int) {
	if index < 0x0800 {
		cpu.RAM[index] = value
	} else if index < 0x2000 {
		cpu.RAM[index-0x0800] = value
	} else if index < 0x2008 {

	} else if index < 0x4000 {

	} else if index < 0x4020 {

	} else if index < 0x6000 {

	} else if index < 0x8000 {

	} else {
		cpu.PrgROM[index-0x8000] = byte(value)
	}
}

func (cpu *Cpu) Read(index int) int {
	if index < 0x0800 {
		return cpu.RAM[index]
	}
	if index < 0x2000 {
		return cpu.RAM[index-0x800]
	}
	if index < 0x2008 {
		panic(index)
	}
	if index < 0x4000 {

	}
	if index < 0x4020 {

	}
	if index < 0x6000 {

	}
	if index < 0x8000 {

	}
	return int(cpu.PrgROM[index-0x8000])
}

func (cpu *Cpu) Reset() {
	f := cpu.Read(0xFFFC)
	s := cpu.Read(0xFFFD)
	cpu.Register.PC = int(s)*256 + int(f)
}

func (cpu *Cpu) Fetch() int {
	ret := int(cpu.Read(cpu.Register.PC))
	cpu.Register.PC++
	return ret
}

func (cpu *Cpu) Run() int {
	opCodeRaw := cpu.Fetch()
	opCode := opCodeList[opCodeRaw]
	opCode.FetchOperand(cpu)
	cpu.Execute(opCode)
	return cycles[opCodeRaw]
}

func (cpu *Cpu) Execute(opCode *OpCode) {
	var data int
	switch opCode.Base {
	case "ADC":
		if opCode.Mode == ADDR_IMMEDIATE {
			data = opCode.Operand
		} else {
			data = cpu.Read(opCode.Operand)
		}
		cpu.Register.A = cpu.Register.A + int(data) + bool2int(cpu.Register.P.Carry)
	case "SBC":
		if opCode.Mode == ADDR_IMMEDIATE {
			data = opCode.Operand
		} else {
			data = cpu.Read(opCode.Operand)
		}
		cpu.Register.A = cpu.Register.A - int(data) + bool2int(!cpu.Register.P.Carry)
	case "AND":
		if opCode.Mode == ADDR_IMMEDIATE {
			data = opCode.Operand
		} else {
			data = cpu.Read(opCode.Operand)
		}
		cpu.Register.A = cpu.Register.A & int(data)
	case "ORA":
		if opCode.Mode == ADDR_IMMEDIATE {
			data = opCode.Operand
		} else {
			data = cpu.Read(opCode.Operand)
		}
		cpu.Register.A = cpu.Register.A | int(data)
	case "EOR":
	case "ASL":
		cpu.Register.A <<= 1
		cpu.Register.P.Carry = (cpu.Register.A & int(math.Pow(2, 7))) != 0
	case "LSR":
		cpu.Register.A >>= 1
		cpu.Register.P.Carry = (cpu.Register.A & int(math.Pow(2, 0))) != 0
	case "ROL":
		cpu.Register.A = cpu.Register.A<<1 + bool2int(cpu.Register.P.Carry)
		cpu.Register.P.Carry = (cpu.Register.A & int(math.Pow(2, 7))) != 0
	case "ROR":
		cpu.Register.A = cpu.Register.A>>1 + bool2int(cpu.Register.P.Carry)*int(math.Pow(2, 7))
		cpu.Register.P.Carry = (cpu.Register.A & int(math.Pow(2, 0))) != 0
	case "BCC":
	case "BCS":
	case "BEQ":
	case "BNE":
	case "BVC":
	case "BVS":
	case "BPL":
	case "BMI":
	case "BIT":
	case "JMP":
		cpu.Register.PC = cpu.Read(opCode.Operand)
	case "JSR":
		// push PC to Stack
		cpu.Register.PC = cpu.Read(opCode.Operand)
	case "RTS":
		// pop Stack to PC
	case "BRK":
	case "RTI":
		// pop Stack
	case "CMP":
		if opCode.Mode == ADDR_IMMEDIATE {
			data = opCode.Operand
		} else {
			data = cpu.Read(opCode.Operand)
		}
		if cpu.Register.A-data > 0 {
			cpu.Register.P.Carry = true
		} else {
			cpu.Register.P.Carry = false
		}
	case "CPX":
		if cpu.Register.X-data > 0 {
			cpu.Register.P.Carry = true
		} else {
			cpu.Register.P.Carry = false
		}
	case "CPY":
		if cpu.Register.Y-data > 0 {
			cpu.Register.P.Carry = true
		} else {
			cpu.Register.P.Carry = false
		}
	case "INC":
		data = cpu.Read(opCode.Operand)
		cpu.Write(opCode.Operand, data+1)
	case "DEC":
		data = cpu.Read(opCode.Operand)
		cpu.Write(opCode.Operand, data-1)
	case "INX":
		cpu.Register.X++
	case "DEX":
		cpu.Register.X--
	case "INY":
		cpu.Register.Y++
	case "DEY":
		cpu.Register.Y--
	case "CLC":
		cpu.Register.P.Carry = false
	case "SEC":
		cpu.Register.P.Carry = true
	case "CLI":
		cpu.Register.P.Interrupt = false
	case "SEI":
		cpu.Register.P.Interrupt = true
	case "CLD":
		cpu.Register.P.Decimal = false
	case "SED":
		cpu.Register.P.Decimal = true
	case "CLV":
		cpu.Register.P.Overflow = false
	case "LDA":
		if opCode.Mode == ADDR_IMMEDIATE {
			data = opCode.Operand
		} else {
			data = cpu.Read(opCode.Operand)
		}
		cpu.Register.A = data
	case "LDX":
		if opCode.Mode == ADDR_IMMEDIATE {
			data = opCode.Operand
		} else {
			data = cpu.Read(opCode.Operand)
		}
		cpu.Register.X = data
	case "LDY":
		if opCode.Mode == ADDR_IMMEDIATE {
			data = opCode.Operand
		} else {
			data = cpu.Read(opCode.Operand)
		}
		cpu.Register.Y = data
	case "STA":
		cpu.Write(opCode.Operand, cpu.Register.A)
	case "STX":
		cpu.Write(opCode.Operand, cpu.Register.X)
	case "STY":
		cpu.Write(opCode.Operand, cpu.Register.Y)
	case "TAX":
		cpu.Register.X = cpu.Register.A
	case "TXA":
		cpu.Register.A = cpu.Register.X
	case "TAY":
		cpu.Register.Y = cpu.Register.A
	case "TYA":
		cpu.Register.A = cpu.Register.Y
	case "TSX":
		cpu.Register.X = cpu.Register.SP
	case "TXS":
		cpu.Register.SP = cpu.Register.A
	case "PHA":
		cpu.PushStack(cpu.Register.A)
	case "PLA":
		cpu.Register.A = cpu.PopStack()
	case "PHP":
		cpu.PushStack(cpu.Register.P.Int())
	case "PLP":
		cpu.Register.P.Set(cpu.PopStack())
	case "NOP":
		return
	}
}

func (cpu *Cpu) PushStack(value int) {
	cpu.RAM[cpu.Register.SP] = value
	cpu.Register.SP++
}

func (cpu *Cpu) PopStack() int {
	cpu.Register.SP--
	return cpu.RAM[cpu.Register.SP]
}

type Sprite struct {
	bitMap [][]int
}

func (sprite *Sprite) Render() {
	for _, bits := range sprite.bitMap {
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

func NewImage(src []byte) *Sprite {
	bitMap := make([][]int, ImageSize)
	for i := 0; i < ImageSize; i++ {
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
	return &Sprite{
		bitMap: bitMap,
	}
}

func debug(args ...interface{}) {
	pp.Println(args...)
}

func bool2int(v bool) int {
	if v {
		return 1
	}
	return 0
}
