package main

import (
	"C"
	"github.com/veandco/go-sdl2/sdl"
)

type APU struct {
	spec sdl.AudioSpec
	channel1Register []int
	channel2Register []int
}

func NewAPU() *APU {
	return &APU{
		channel1Register: make([]int, 4),
		channel2Register: make([]int, 4),
	}
}

func (apu *APU) Write(index, value int) {
	if index == 0x0000 {
		apu.channel1Register[0] = value
	}
	if index == 0x0001 {
		apu.channel1Register[1] = value
	}
	if index == 0x0002 {
		apu.channel1Register[2] = value
	}
	if index == 0x0003 {
		apu.channel1Register[3] = value
	}
	if index == 0x0004 {
		apu.channel2Register[0] = value
	}
	if index == 0x0005 {
		apu.channel2Register[1] = value
	}
	if index == 0x0006 {
		apu.channel2Register[2] = value
	}
	if index == 0x0007 {
		apu.channel2Register[3] = value
	}
}

func (apu *APU) Read(index int) int {
	if index == 0x0000 {
		return apu.channel1Register[0]
	}
	if index == 0x0001 {
		return apu.channel1Register[1]
	}
	if index == 0x0002 {
		return apu.channel1Register[2]
	}
	if index == 0x0003 {
		return apu.channel1Register[3]
	}
	if index == 0x0004 {
		return apu.channel2Register[0]
	}
	if index == 0x0005 {
		return apu.channel2Register[1]
	}
	if index == 0x0006 {
		return apu.channel2Register[2]
	}
	if index == 0x0007 {
		return apu.channel2Register[3]
	}
	return 0
}
