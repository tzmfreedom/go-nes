package main

import "math"

const SpriteSize = 8

type Sprite struct {
	bitMap [][]int
}

func NewSprite(src []byte) *Sprite {
	bitMap := make([][]int, SpriteSize)
	for i := 0; i < SpriteSize; i++ {
		bit := make([]int, SpriteSize)
		for j := 0; j < SpriteSize; j++ {
			b := 0
			mul := int(math.Pow(2, float64(j)))
			if (int(src[i]) & mul) != 0 {
				b += 1
			}
			if (int(src[i+SpriteSize]) & mul) != 0 {
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
