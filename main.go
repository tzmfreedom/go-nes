package main

import (
	"fmt"
	"github.com/k0kubun/pp"
	"github.com/veandco/go-sdl2/sdl"
	"io/ioutil"
	"math"
	"os"
	"time"
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
	cpu.PPU = ppu
	cpu.Reset()
	nes := NewNES(cpu, ppu, sprites)
	nes.run()
	defer nes.close()
}

type NES struct {
	cpu        *Cpu
	ppu        *PPU
	background *BackGround
	pallet     *Pallet
	sprites    []*Sprite
	spritesData []*SpriteData
	renderer   *sdl.Renderer
	window     *sdl.Window
	frame      int
	time       int64
}

func NewNES(cpu *Cpu, ppu *PPU, sprites []*Sprite) *NES {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		256, 240, sdl.WINDOW_SHOWN | sdl.WINDOW_RESIZABLE)
	if err != nil {
		panic(err)
	}
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	renderer.Clear()
	sdl.PollEvent()

	return &NES{
		cpu: cpu,
		ppu: ppu,
		sprites: sprites,
		renderer: renderer,
		window: window,
	}
}

func (nes *NES) close() {
	nes.renderer.Destroy()
	nes.window.Destroy()
	sdl.Quit()
}

func (nes *NES) run() error {
	for {
		cycle := nes.cpu.Run()
		background, pallet, sprites := nes.ppu.Run(cycle * 3)
		if background != nil {
			nes.render(background, pallet, sprites)
		}
	}
	return nil
}

func (nes *NES) render(background *BackGround, pallet *Pallet, sprites []*SpriteData) {
	for i, line := range background.tiles {
		for j, tile := range line {
			sprite := nes.sprites[tile.spriteId]
			for y, line := range sprite.bitMap {
				for x, bit := range line {
					c := pallet.getBackgroundColor(tile.palletId, bit)
					nes.renderer.SetDrawColor(c.R, c.G, c.B, 0xff)
					nes.renderer.DrawPoint(int32(j*SpriteSize+x), int32(i*SpriteSize+y))
				}
			}
		}
	}

	for _, sprite := range sprites	{
		s := nes.sprites[256+sprite.spriteId]
		//isVerticalReverse := sprite.attr & 0x80
		//isHoriozntalReverse := sprite.attr & 0x40
		//isPriority := sprite.attr & 0x20
		palletId := sprite.attr & 0x03
		for y, line := range s.bitMap {
			for x, bit := range line {
				c := pallet.getSpriteColor(palletId, bit)
				nes.renderer.SetDrawColor(c.R, c.G, c.B, 0xff)
				nes.renderer.DrawPoint(int32(sprite.x+x), int32(sprite.y+y))
			}
		}
	}
	nes.frame++
	if nes.frame == 60 {
		nes.frame -= 60
		if nes.time > 0 {
			debug((time.Now().UnixNano() - nes.time)/1000000000)
			debug(60*1000000000/(time.Now().UnixNano() - nes.time))
		}
		nes.time = time.Now().UnixNano()
	}
	nes.renderer.Present()
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
	spriteId int
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
	r.Negative = int(math.Pow(2, 7)) != 0
	r.Overflow = int(math.Pow(2, 6)) != 0
	r.Reserved = int(math.Pow(2, 5)) != 0
	r.Break = int(math.Pow(2, 4)) != 0
	r.Decimal = int(math.Pow(2, 3)) != 0
	r.Interrupt = int(math.Pow(2, 2)) != 0
	r.Zero = int(math.Pow(2, 1)) != 0
	r.Carry = int(math.Pow(2, 0)) != 0
}

func debug(args ...interface{}) {
	if true {
		pp.Println(args...)
	}
}

func exit() {
	os.Exit(1)
}

func bool2int(v bool) int {
	if v {
		return 1
	}
	return 0
}
