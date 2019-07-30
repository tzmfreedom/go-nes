package main

// typedef unsigned char Uint8;
// void MyCallback(void *userdata, Uint8 *stream, int len);
import "C"
import (
	"github.com/veandco/go-sdl2/sdl"
	"reflect"
	"unsafe"
)

type Audio struct {
	spec *sdl.AudioSpec
}

func NewAudio() *Audio {
	return &Audio{
		spec: &sdl.AudioSpec{
			Freq: 44100,              // DSP frequency (samples per second)
			Format: sdl.AUDIO_S16SYS, // audio data format
			Channels: 1,              // number of separate sound channels
			Silence: 0,               // audio buffer silence value (calculated)
			Samples: 2048,            // audio buffer size in samples (power of 2)
			Size: 0,                  // audio buffer size in bytes (calculated)
			Callback: sdl.AudioCallback(C.MyCallback),       // the function to call when the audio device needs more data
			UserData: nil,            // a pointer that is passed to callback (otherwise ignored by SDL)
		},
	}
}

func (audio *Audio) Run() {
}

var pre = 0

//export MyCallback
func MyCallback(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Uint8)(unsafe.Pointer(&hdr))
	freq := (apu.channel1Register[3] & 0x7)<<8 + apu.channel1Register[2]
	if freq == 0 {
		return
	}
	th := 44100/freq
	for i := 0; i < n; i++ {
		tmp := i + pre
		if (tmp/th)%2 == 0 {
			buf[i] = 128
		} else {
			buf[i] = 0
		}
	}
	pre = (pre+2048)%th
}