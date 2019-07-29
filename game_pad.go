package main


type GamePad struct {
	A bool
	B bool
	Select bool
	Start bool
	Up bool
	Down bool
	Left bool
	Right bool
}

func (p *GamePad) Reset() {
	p.A = false
	p.B = false
	p.Select = false
	p.Start = false
	p.Up = false
	p.Down = false
	p.Left = false
	p.Right = false
}
