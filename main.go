package main

import (
	"fmt"
	"github.com/k0kubun/pp"
	"github.com/veandco/go-sdl2/sdl"
	"io/ioutil"
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

	cpu := NewCpu(prgRom)
	nes := NewNES(cpu, chrRom)
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

func NewNES(cpu *Cpu, chrRom []byte) *NES {
	sprites := make([]*Sprite, 512)
	for i := 0; i < 512; i++ {
		index := i * 16
		sprites[i] = NewSprite(chrRom[index : index+16])
	}

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
		ppu: cpu.PPU,
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
	spIndex := 0
	bgIndex := 0
	if nes.ppu.controlRegister & 0x08 != 0 {
		spIndex = 0x1000/16
	}
	if nes.ppu.controlRegister & 0x10 != 0 {
		bgIndex = 0x1000/16
	}
	start := time.Now().UnixNano()
	for i, line := range background.tiles {
		for j, tile := range line {
			sprite := nes.sprites[bgIndex+tile.spriteId]
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
		s := nes.sprites[spIndex+sprite.spriteId]
		isVerticalReverse := sprite.attr & 0x80 != 0
		isHoriozntalReverse := sprite.attr & 0x40 != 0
		//isPriority := sprite.attr & 0x20
		palletId := sprite.attr & 0x03
		for y, line := range s.bitMap {
			for x, bit := range line {
				if isVerticalReverse {
					y = SpriteSize - y
				}
				if isHoriozntalReverse {
					x = SpriteSize - x
				}
				c := pallet.getSpriteColor(palletId, bit)
				nes.renderer.SetDrawColor(c.R, c.G, c.B, 0xff)
				nes.renderer.DrawPoint(int32(sprite.x+x), int32(sprite.y+y))
			}
		}
	}

	end := time.Now().UnixNano()
	nes.frame++
	if nes.frame == 60 {
		nes.frame -= 60
		if nes.time > 0 {
			debug(sprites[0:2])
			debug(end-start)
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
	return bool2int(r.Negative)<<7 +
		bool2int(r.Overflow)<<6 +
		bool2int(r.Reserved)<<5 +
		bool2int(r.Break)<<4 +
		bool2int(r.Decimal)<<3 +
		bool2int(r.Interrupt)<<2 +
		bool2int(r.Zero)<<1 +
		bool2int(r.Carry)
}

func (r *StatusRegister) Set(v int) {
	r.Negative = v & 0x80 != 0
	r.Overflow = v & 0x40 != 0
	r.Reserved = v & 0x20 != 0
	r.Break = v & 0x10 != 0
	r.Decimal = v & 0x08 != 0
	r.Interrupt = v & 0x04 != 0
	r.Zero = v & 0x02 != 0
	r.Carry = v & 0x01 != 0
}

func debug(args ...interface{}) {
	if false {
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
