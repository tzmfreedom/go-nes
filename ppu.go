package main

type PPU struct {
	RAM           []int
	cycle         int
	line          int
	background    *BackGround
	sprites       []*Sprite
	addr          int
	isWriteHigher bool
	controlRegister int
	controlRegister2 int
	statusRegister int
	spriteMemAddr int
	isWriteScrollV bool
	scrollH int
	scrollV int
}

func NewPPU() *PPU {
	return &PPU{
		RAM:        make([]int, 0x4000),
		background: NewBackGround(),
	}
}

func (ppu *PPU) Read(index int) int {
	switch index {
	case 0x0000:
		// no action
	case 0x0001:
		// no action
	case 0x0002:
		debug(ppu.statusRegister >> 7)
		return ppu.statusRegister
	case 0x0003:
		// no action
	case 0x0004:
		// no action
	case 0x0005:
		// no action
	case 0x0006:
		// no action
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
		ppu.controlRegister = data
	case 0x0001:
		ppu.controlRegister2 = data
	case 0x0002:
		// no action
	case 0x0003:
		ppu.spriteMemAddr = index
	case 0x0004:
		ppu.RAM[ppu.spriteMemAddr] = index
		ppu.spriteMemAddr += 0x01
	case 0x0005:
		if ppu.isWriteScrollV {
			ppu.scrollH = data
			ppu.isWriteScrollV = false
		} else {
			ppu.scrollV = data
			ppu.isWriteScrollV = true
		}
	case 0x0006:
		if ppu.isWriteHigher {
			ppu.addr += data
			ppu.isWriteHigher = false
		} else {
			ppu.addr = data * 256
			ppu.isWriteHigher = true
		}
	case 0x0007:
		ppu.RAM[ppu.addr] = data
		ppu.addr += 0x01 // TODO: impl
	}
}

func (ppu *PPU) Run(cycle int) (*BackGround, *Pallet) {
	ppu.cycle += cycle
	if ppu.cycle >= 341 {
		ppu.cycle -= 341
		ppu.line++
		if ppu.line < 240 && (ppu.line-1)%8 == 0 {
			ppu.BuildBackGround()
		}
		if ppu.line == 241 {
			ppu.statusRegister |= 0x80
		}
		if ppu.line == 262 {
			ppu.line = 0
			ppu.statusRegister &= 0x7F
			background := ppu.background
			ppu.background = NewBackGround()
			return background, ppu.getPallet()
		}
	}
	return nil, nil
}

func (ppu *PPU) BuildBackGround() {
	y := (ppu.line - 1) / 8
	for x := 0; x < 32; x++ {
		tile := ppu.BuildTile(x, y)
		ppu.background.Add(x, y, tile)
	}
}

func (ppu *PPU) getPallet() *Pallet {
	return NewPallet(ppu.RAM[0x3F00:0x3F10])
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
	return (palletBlock >> (blockId * 2)) & 3
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

type Pallet struct {
	src []int
}

func NewPallet(src []int) *Pallet {
	return &Pallet{
		src: src,
	}
}

func (p *Pallet) getColor(palletId int, bit int) *RGB {
	return colors[p.src[palletId*4+bit]]
}
