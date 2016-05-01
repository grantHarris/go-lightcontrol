package goperiodicity

import (
 //"math/rand"
 "time"
 "math"
 "fmt"
 //"log"
)
import "github.com/grantHarris/go-nanokontrol2"
import "github.com/lucasb-eyer/go-colorful"
import "github.com/mkb218/go-osc/lib"
import "flag"

type MODE uint8;
type WAVEFORM uint8;


const (
	ADD MODE = iota
	SUBTRACT
	MULTIPLY
	DIVIDE
	SCREEN
	OVERLAY
)

type Layer struct{
	R float64;
	G float64;
	B float64;
	A float64;
	mode MODE
	enabled bool
}

type Fixture struct {
	buffer []byte
	name string
    R_addr uint8
    G_addr uint8
    B_addr uint8
    W_addr uint8
    BRIGHTNESS_addr uint8
}

func (f *Fixture) Set(r float64, g float64, b float64){
	
	if f.BRIGHTNESS_addr != 0{
		f.buffer[f.BRIGHTNESS_addr] = 255
	}

	/*
		Case where the fixture has a white channel. Subtract the common values
		from each r g b and push to white channel
	*/
	if f.W_addr != 0{
		if r > 0 && g > 0 && b > 0 {
			common := math.Min(math.Min(r, g), b)
			r = r - common
			g = g - common
			b = b - common
			f.buffer[f.W_addr] = byte(math.Max(math.Min(common * 255, 255), 0))
		}
	}else{
		f.buffer[f.W_addr] = 0;
	}


	f.buffer[f.R_addr] = byte(math.Max(math.Min(r * 255, 255), 0))
	f.buffer[f.G_addr] = byte(math.Max(math.Min(g * 255, 255), 0))
	f.buffer[f.B_addr] = byte(math.Max(math.Min(b * 255, 255), 0))
}

func (f *Fixture) Render(layers []Layer){
	var r, g, b float64

	for _, layer := range layers {
		if layer.enabled == true{
			layer_r := layer.R * layer.A
			layer_g := layer.G * layer.A
			layer_b := layer.B * layer.A

		    switch layer.mode{
				case ADD:
					r = r + layer_r
					g = g + layer_g
					b = b + layer_b
				case SUBTRACT:
					r = r - layer_r
					g = g - layer_g
					b = b - layer_b
				case MULTIPLY:
					r = r * layer_r
					g = g * layer_g
					b = b * layer_b
				case DIVIDE:
					r = r / layer_r
					g = g / layer_g
					b = b / layer_b
				case SCREEN:
					r = 1 - (1 - r)*(1 - layer_r)
					g = 1 - (1 - g)*(1 - layer_g)
					b = 1 - (1 - b)*(1 - layer_b)
				case OVERLAY:
					if layer.A < 0.5{
						r = 2*r*layer_r
						g = 2*g*layer_g
						b = 2*b*layer_b
					}else{
						r = 1 - 2*(1 - r)*(1 - layer_r)
						g = 1 - 2*(1 - g)*(1 - layer_g)
						b = 1 - 2*(1 - b)*(1 - layer_b)
					}
		    }
			r = math.Min(math.Max(r, 0.0), 1.0) 
			g = math.Min(math.Max(g, 0.0), 1.0) 
			b = math.Min(math.Max(b, 0.0), 1.0)
		}
	}
	f.Set(r, g, b)
}

// func (f *Fixture) Get()(r, g, b byte){
// 	return f.buffer[f.R_addr], f.buffer[f.G_addr], f.buffer[f.B_addr]
// }

func (f *Fixture) Print(){
	fmt.Print("(", f.name ,": ",f.buffer[f.R_addr], ", ", f.buffer[f.G_addr], ", ", f.buffer[f.B_addr])
	if f.W_addr != 0{
		fmt.Print(", ", f.buffer[f.W_addr])
	}
	fmt.Print(") ")
}


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

// type Mode struct{
// 	hue Period
// 	sat Period
// 	val Period
// 	alpha Period
// 	layers []Layers
// 	fixtures []Fixture
// }

// func (m *Mode) Iterate(){
// 	for i, fixture := range m.fixtures{
// 		for _, layer := range m.layers{
// 			fixture.Set(m.hue.Value(i), m.sat.Value(i), m.val.Value(i))
// 		}
// 	}
// }

// func NewMode() Mode{
// 	hue := Period{0.0, 0.0, 1.0, SIN}
// 	sat := Period{0.0, 0.0, 1.0, SIN}
// 	val := Period{0.0, 0.0, 1.0, SIN}
// 	alpha := 1.0
// 	m := Mode {hue, sat, val, alpha}
// 	return m
// }

func scale(old_min, old_max, new_min, new_max, value float64) float64{
    return ((value - old_min) / (old_max - old_min) ) * (new_max - new_min) + new_min
}

func main() {

	var ip = flag.String("a", "127.0.0.1", "OSC output Address")
	var port = flag.String("p", "1235", "OSC Output Port")
	var universe = flag.String("u", "/dmx/universe/1", "OSC Output Universe")

	var verbose = flag.Bool("v", false, "Verbose")
	//var http = flag.String("h", "6969, "HTTP Server")

	flag.Parse()

    nanokontrol2 := nanokontrol2.Initialize()    
    osc_buffer := make([]byte, 512)

    address := osc.NewAddress(ip, port)

    var hue_index float64
    hue_index = 0

    var value_index float64
    value_index = 0

	par1 := Fixture{osc_buffer, "Par 1", 4, 5, 6, 7, 3}
	par2 := Fixture{osc_buffer, "Par 2", 14, 15, 16, 17, 13}
	par3 := Fixture{osc_buffer, "Par 3", 24, 25, 26, 27, 23}
	par4 := Fixture{osc_buffer, "Par 4", 34, 35, 36, 37, 33}
	kitchen_1 := Fixture{osc_buffer, "Kitchen 1", 65, 64, 66, 0, 0}
	kitchen_2 := Fixture{osc_buffer, "Kitchen 2", 75, 74, 76, 0, 0}

	fixtures := []Fixture{}

	fixtures = append(fixtures, par1)
	fixtures = append(fixtures, par2)
	fixtures = append(fixtures, par3)
	fixtures = append(fixtures, par4)
	fixtures = append(fixtures, kitchen_1)
	fixtures = append(fixtures, kitchen_2)

     for{
     	hue_period := nanokontrol2.Get(16)
     	hue_width := nanokontrol2.Get(0)

     	hue_period = float64(1)
     	hue_width = float64(1)

		hue_index = hue_index + (hue_period*math.Pi/220)
        hue := scale(-1, 1, 0, 360*hue_width, math.Sin(hue_index) + math.Sin(hue_index/2))
     	
     	alpha_period := nanokontrol2.Get(17)
     	alpha_low := nanokontrol2.Get(1)
     	alpha_high := nanokontrol2.Get(2)

     	alpha_period = float64(1)
     	alpha_low = float64(1)
     	alpha_high = float64(1)

       	value_index = value_index + (alpha_period*math.Pi/220)
        value := scale(-1, 1, alpha_low, alpha_high, math.Sin(value_index))
        
        color := colorful.Hsv(hue, 1, value)

        for _,fixture := range fixtures{
        	fixture.Set(color.R, color.G, color.B)
        	if *verbose == true{
        		fixture.Print()
        	}
        }

        fmt.Println()

        message := make(osc.Message, 0)
        message = append(message, osc.Blob(osc_buffer))
        message.Send(address, *universe)
     	
        //DMX has a 44Hz max refresh rate
        time.Sleep(time.Second / 44)
    }
}