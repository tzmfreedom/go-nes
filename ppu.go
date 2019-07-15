package main

type PPU struct {
	RAM           []int
	cycle         int
	line          int
	background    *BackGround
	spriteRAM     []int
	sprites       []*SpriteData
	addr          int
	isWriteHigher bool
	controlRegister int
	controlRegister2 int
	statusRegister int
	spriteMemAddr int
	isWriteScrollX bool
	scrollX int
	scrollY int
	interrupts *Interrupts
}

func NewPPU(interrupts *Interrupts) *PPU {
	return &PPU{
		RAM:        make([]int, 0x4000),
		spriteRAM:  make([]int, 0x100),
		interrupts: interrupts,
	}
}

func (ppu *PPU) Read(index int) int {
	switch index {
	case 0x0000:
		// no action
	case 0x0001:
		// no action
	case 0x0002:
		r := ppu.statusRegister
		ppu.statusRegister &= 0x7F
		ppu.isWriteScrollX = false
		return r
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
		ppu.spriteMemAddr = data
	case 0x0004:
		ppu.spriteRAM[ppu.spriteMemAddr] = data
		ppu.spriteMemAddr += 0x01
	case 0x0005:
		if ppu.isWriteScrollX {
			ppu.scrollY = data
			ppu.isWriteScrollX = false
		} else {
			ppu.scrollX = data
			ppu.isWriteScrollX = true
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
		if ppu.controlRegister&0x04 == 0 {
			ppu.addr += 0x01
		} else {
			ppu.addr += 0x20
		}
	}
}

func (ppu *PPU) Run(cycle int) (bool, *Pallet, []*SpriteData) {
	ppu.cycle += cycle
	if ppu.line == 0 {
		ppu.buildSprites()
	}
	if ppu.cycle >= 341 {
		ppu.cycle -= 341
		ppu.line++
		//if ppu.line < 241 && ppu.line%8 == 0 {
		//	ppu.BuildBackGround()
		//}
		if ppu.line == 241 {
			ppu.statusRegister |= 0x80
			if ppu.controlRegister & 0x80 != 0 {
				ppu.interrupts.Nmi = true
			}
		}
		if ppu.line == 262 {
			ppu.line = 0
			ppu.statusRegister &= 0x7F
			//background := ppu.background
			//ppu.background = NewBackGround()
			return true, ppu.getPallet(), ppu.sprites
		}
	}
	return false, nil, nil
}

func (ppu *PPU) buildSprites() {
	ppu.sprites = []*SpriteData{}
	for i := 0; i < 0xff; i+=4 {
		y := ppu.spriteRAM[i]+1
		spriteId := ppu.spriteRAM[i+1]
		attr := ppu.spriteRAM[i+2]
		x := ppu.spriteRAM[i+3]
		ppu.sprites = append(ppu.sprites, &SpriteData{
			spriteId: spriteId,
			attr: attr,
			x: x,
			y: y,
		})
	}
}

//func (ppu *PPU) BuildBackGround() {
//	y := (ppu.line - 1) / 8
//	for x := 0; x < 32; x++ {
//		tile := ppu.BuildTile(x, y)
//		ppu.background.Add(x, y, tile)
//	}
//}

func (ppu *PPU) getPallet() *Pallet {
	return NewPallet(ppu.RAM[0x3F00:0x3F20])
}

//func (ppu *PPU) BuildTile(x, y int) *Tile {
//	palletId := ppu.getPalletId(x, y)
//	spriteId := ppu.getSpriteId(x, y)
//	return &Tile{
//		spriteId: spriteId,
//		palletId: palletId,
//	}
//}

func (ppu *PPU) getPalletId(x, y int) int {
	offset := 0x23C0
	if x + ppu.scrollX/8 >= 32 {
		offset += 0x0400
	}
	if y + ppu.scrollY/8 >= 32 {
		offset += 0x0800
	}

	tmpX := x / 2
	tmpY := y / 2
	palletBlock := ppu.RAM[tmpX+tmpY*8+offset]

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

func (ppu *PPU) getSpriteId(x, y int) int {
	offset := 0x2000
	if x + ppu.scrollX/8 >= 32 {
		offset += 0x0400
	}
	if y + ppu.scrollY/8 >= 32 {
		offset += 0x0800
	}

	return ppu.RAM[x+y*32+offset]
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
	src[0x04] = src[0x00]
	src[0x08] = src[0x00]
	src[0x0C] = src[0x00]
	src[0x10] = src[0x00]
	src[0x14] = src[0x04]
	src[0x18] = src[0x08]
	src[0x1C] = src[0x0C]
	return &Pallet{
		src: src,
	}
}

func (p *Pallet) getBackgroundColor(palletId int, bit int) *RGB {
	return colors[p.src[palletId*4+bit]]
}

func (p *Pallet) getSpriteColor(palletId int, bit int) *RGB {
	return colors[p.src[0x10+palletId*4+bit]]
}

type SpriteData struct {
	x, y int
	spriteId int
	attr int
}