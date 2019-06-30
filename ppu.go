package main

type PPU struct {
	RAM        []int
	cycle      int
	line       int
	background *BackGround
	sprites    []*Sprite
	addr int
	isWriteHigher bool
}

func NewPPU() *PPU {
	return &PPU{
		RAM: make([]int, 0x4000),
		background: NewBackGround(),
	}
}

func (ppu *PPU) Read(index int) int {
	switch index {
	case 0x0000:
	case 0x0001:
	case 0x0002:
	case 0x0003:
	case 0x0004:
	case 0x0005:
	case 0x0006:
	case 0x0007:
		data := ppu.RAM[ppu.addr]
		ppu.addr += 0x01
		return data
	}
	return 0
}

func (ppu *PPU) Write(index, data int) {
	switch index {
	case 0x0000:
	case 0x0001:
	case 0x0002:
	case 0x0003:
	case 0x0004:
	case 0x0005:
	case 0x0006:
		if ppu.isWriteHigher {
			ppu.addr += data
			ppu.isWriteHigher = false
		} else {
			ppu.addr = data * 256
			ppu.isWriteHigher = true
		}
	case 0x0007:
		debug("address", ppu.addr)
		ppu.RAM[ppu.addr] = data
		ppu.addr += 0x01 // TODO: impl
	}
}

func (ppu *PPU) Run(cycle int) *BackGround {
	ppu.cycle += cycle
	if ppu.cycle >= 341 {
		ppu.cycle -= 341
		ppu.line++
		if ppu.line < 240 && (ppu.line-1)%8 == 0 {
			ppu.BuildBackGround()
		}
		if ppu.line == 262 {
			ppu.line = 0
			//debug(ppu.background)
			return ppu.background
		}
	}
	return nil
}

func (ppu *PPU) BuildBackGround() {
	y := (ppu.line-1)/ 8
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
	palletBlock := ppu.RAM[tmpX+tmpY*8+0x23C0]

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
	spriteId := ppu.RAM[x+y*32+0x2000]
	return ppu.sprites[spriteId]
}

type RGB struct {
	R uint8
	G uint8
	B uint8
}

