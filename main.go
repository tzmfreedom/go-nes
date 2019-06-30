package main

import (
	"fmt"
	"github.com/k0kubun/pp"
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
	}

	ppu := NewPPU()
	sprites := make([]*Sprite, 512)
	for i := 0; i < 512; i++ {
		index := i * 16
		sprites[i] = NewSprite(chrRom[index : index+16])
	}
	ppu.sprites = sprites
	cpu.PPU = ppu
	nes := &NES{
		cpu:   cpu,
		ppu:   ppu,
		cycle: 0,
	}
	nes.Run()
	log.Fatal(err)
}

type NES struct {
	cpu        *Cpu
	ppu        *PPU
	cycle      int
	background *BackGround
}

func (nes *NES) Run() {
	nes.cpu.Reset()
	for {
		nes.cycle += nes.cpu.Run()
		background := nes.ppu.Run(nes.cycle * 3)
		if background != nil {
			nes.render(background)
		}
	}
}

func (nes *NES) render_(background *BackGround) {
	//fmt.Print("\033[2J")
	//fmt.Print("\r")
	//fmt.Print("\033[;H")
	for i, line := range background.tiles {
		for j, tile := range line {
			nes.renderTile(j, i, tile)
		}
	}
	os.Exit(0)
}

func (nes *NES) render(background *BackGround) {
	fmt.Print("\033[2J")
	fmt.Print("\r")
	fmt.Print("\033[;H")
	for _, line := range background.tiles {
		for k := 0; k < 8; k++ {
			for _, tile := range line {
				for l := 0; l < 8; l++ {
					if tile.img.bitMap[k][l] == 0 {
						fmt.Print(" ")
					} else {
						fmt.Print("*")
					}
				}
				fmt.Print("  ")
			}
			fmt.Println()
		}
	}
	os.Exit(0)
}

func (nes *NES) renderTile(x, y int, tile *Tile) {
	for _, line := range tile.img.bitMap {
		for _, bit := range line {
			if bit == 0 {
				fmt.Print(" ")
			} else {
				fmt.Print("*")
			}
			//img, _ := ebiten.NewImage(1, 1, 0)
			//c := colors[bit]
			//err := img.Fill(color.RGBA{c.R, c.G, c.B, 0xff})
			//if err != nil {
			//	panic(err)
			//}
			//options := &ebiten.DrawImageOptions{}
			//options.GeoM.Translate(float64(x*SpriteSize+j), float64(y*SpriteSize+i))
			//err = screen.DrawImage(img, options)
			//if err != nil {
			//	panic(err)
			//}
		}
		fmt.Println()
	}
	fmt.Println()
}

type Register struct {
	A  int
	X  int
	Y  int
	P  *StatusRegister
	SP int
	PC int
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

type Tile struct {
	img      *Sprite
	palletId int
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
	r.Negative  = int(math.Pow(2, 7)) != 0
	r.Overflow  = int(math.Pow(2, 6)) != 0
	r.Reserved  = int(math.Pow(2, 5)) != 0
	r.Break     = int(math.Pow(2, 4)) != 0
	r.Decimal   = int(math.Pow(2, 3)) != 0
	r.Interrupt = int(math.Pow(2, 2)) != 0
	r.Zero      = int(math.Pow(2, 1)) != 0
	r.Carry     = int(math.Pow(2, 0)) != 0
}

func debug(args ...interface{}) {
	if false {
		pp.Println(args...)
	}
}

func bool2int(v bool) int {
	if v {
		return 1
	}
	return 0
}
