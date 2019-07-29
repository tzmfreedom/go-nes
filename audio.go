package main

// typedef unsigned char Uint8;
// void MyCallback(void *userdata, Uint8 *stream, int len);
// void SineWave(void *userdata, Uint8 *stream, int len);
import "C"
import (
	"github.com/veandco/go-sdl2/sdl"
	"reflect"
	"unsafe"
)

type Audio struct {
	spec sdl.AudioSpec
}

func NewAudio() *Audio {
	return &Audio{
		spec: sdl.AudioSpec{
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

//export MyCallback
func MyCallback(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length)
	len := n/2
	step := 0
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Uint8)(unsafe.Pointer(&hdr))
	for i := 0; i < len; i++ {
		if (step/10000)%2 == 0 {
			buf[i] = 1
		} else {
			buf[i] = 0
		}
	}
}