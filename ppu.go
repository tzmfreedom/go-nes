package main

type PPU struct {
	RAM           []int
	cycle         int
	line          int
	spriteRAM     []int
	sprites       []*SpriteData
	addr          int
	buffer        int
	controlRegister int
	controlRegister2 int
	statusRegister int
	spriteMemAddr int
	writeToggle bool
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

func (ppu *PPU) IsBackgroundEnabled() bool {
	return ppu.controlRegister2 & 0x08 != 0
}

func (ppu *PPU) IsSpriteEnabled() bool {
	return ppu.controlRegister2 & 0x10 != 0
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
		ppu.writeToggle = false
		return r
	case 0x0003:
		// no action
	case 0x0004:
		if ppu.spriteMemAddr%4 == 2 {
			return ppu.spriteRAM[ppu.spriteMemAddr]&0xE3
		}
		return ppu.spriteRAM[ppu.spriteMemAddr]
	case 0x0005:
		// no action
	case 0x0006:
		// no action
	case 0x0007:
		data := ppu.buffer
		ppu.buffer = ppu.RAM[ppu.addr]
		if ppu.controlRegister&0x04 == 0 {
			ppu.addr += 0x01
		} else {
			ppu.addr += 0x20
		}
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
		ppu.spriteMemAddr = data&0xFF
	case 0x0004:
		ppu.spriteRAM[ppu.spriteMemAddr] = data
		ppu.spriteMemAddr = (ppu.spriteMemAddr+1)&0xFF
	case 0x0005:
		if ppu.writeToggle {
			if data < 240 {
				ppu.scrollY = data
			} else {
				ppu.scrollY = 0
			}
			ppu.writeToggle = false
		} else {
			ppu.scrollX = data
			ppu.writeToggle = true
		}
	case 0x0006:
		if ppu.writeToggle {
			ppu.addr += data
			ppu.writeToggle = false
		} else {
			ppu.addr = data << 8
			ppu.writeToggle = true
		}
	case 0x0007:
		addr := ppu.addr
		if ppu.addr == 0x3F10 ||
			ppu.addr == 0x3F14 ||
			ppu.addr == 0x3F18 ||
			ppu.addr == 0x3F1C {
			addr -= 0x3F10
		}
		//if addr == 8683 && data == 199 {
		//	debug(data)
		//	debug(current)
		//}
		//debug(addr)
		ppu.RAM[addr] = data
		if ppu.controlRegister&0x04 == 0 {
			ppu.addr += 0x01
		} else {
			ppu.addr += 0x20
		}
		if ppu.addr == 8685 {
			//debug(current)
			//debug(newCpu.Register.X)
			//debug(ppu.addr)
			//debug(ppu.RAM[ppu.addr])
		}
		ppu.addr &= 0x3FFF
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
		if ppu.line == 1 && ppu.IsSpriteEnabled() && ppu.IsBackgroundEnabled() {
			ppu.controlRegister2 |= 0x40
		}
		if ppu.line == 241 {
			ppu.statusRegister |= 0x80
			if ppu.controlRegister & 0x80 != 0 {
				ppu.interrupts.Nmi = true
			}
		}
		if ppu.line == 262 {
			ppu.line = 0
			ppu.statusRegister &= 0x7F
			ppu.controlRegister2 &= 0xBF
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

func (ppu *PPU) getPallet() *Pallet {
	return NewPallet(ppu.RAM[0x3F00:0x3F20])
}

func (ppu *PPU) getPalletId(x, y, offset int) int {
	tmpX := x / 4
	tmpY := y / 4
	palletBlock := ppu.RAM[tmpX+tmpY*8+offset+0x03C0]

	cmpX := (x/2) % 2
	cmpY := (y/2) % 2
	var operand uint
	if cmpX == 0 {
		if cmpY == 0 {
			operand = 0
		} else {
			operand = 4
		}
	} else {
		if cmpY == 0 {
			operand = 2
		} else {
			operand = 6
		}
	}
	return (palletBlock >> operand) & 0x03
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
	return colors[p.src[palletId<<2+bit]]
}

func (p *Pallet) getSpriteColor(palletId int, bit int) *RGB {
	return colors[p.src[0x10+palletId<<2+bit]]
}

type SpriteData struct {
	x, y int
	spriteId int
	attr int
}