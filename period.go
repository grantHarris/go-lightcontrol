package goperiodicity

const (
	SIN WAVEFORM = iota
	COS
	SAW
	TRI
	SQU
)


type Period struct{
	index float64
	period float64
	width float64
	fixture_offset float64
	waveform WAVEFORM
}


func (p *Period) SetPeriod(period float64){
	p.period = period
}

func (p *Period) Increment(){
	p.index = p.index + p.period
}


func (p *Period) Value(fixture_index int) float64{
	switch p.waveform{
		case SIN:
			return scale(-1, 1, 0, 360 * p.width, math.Sin(p.index + p.fixture_offset*float64(fixture_index)))
		case COS:
			return scale(-1, 1, 0, 360 * p.width, math.Cos(p.index + p.fixture_offset*float64(fixture_index)))
		// case SAW:
		// case TRI
		// case SQU:
	}
	return 0.0
}