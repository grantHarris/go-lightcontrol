package main

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

type MODE uint8;
type WAVEFORM uint8;

type Fixture struct {
	buffer []byte    
    R_addr uint8
    G_addr uint8
    B_addr uint8
    W_addr uint8
    BRIGHTNESS_addr uint8
}

const (
	ADD MODE = iota
	SUBTRACT
	MULTIPLY
	DIVIDE
	SCREEN
	OVERLAY
)

const (
	SIN WAVEFORM = iota
	COS
	SAWTOOTH
	TRIANGLE
	SQUARE
)

type Layer struct{
	R float64;
	G float64;
	B float64;
	A float64;
	mode MODE
	enabled bool
}

func (f *Fixture) Set(r float64, g float64, b float64){
	
	if f.BRIGHTNESS_addr != 0{
		f.buffer[f.BRIGHTNESS_addr] = 255
	}

	/*
		Case where the fixture has a white channel. Subtract the common values
		from each r g b and push to white channel
	*/
	if f.W_addr != 0 && r > 0 && g > 0 && b > 0 {
		common := math.Min(math.Min(r, g), b)
		r = r - common
		g = g - common
		b = b - common
		f.buffer[f.W_addr] = byte(math.Max(math.Min(common * 255, 255), 0))
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

func (f *Fixture) Get()(r, g, b byte){
	return f.buffer[f.R_addr], f.buffer[f.G_addr], f.buffer[f.B_addr]
}

func (f *Fixture) Print(){
	fmt.Println()
}


type FixtureSet struct{
	set []Fixture
}

type Period struct{
	index float64
	period float64
	width float64
	waveform WAVEFORM
	nano_channel uint8
}

type Mode struct{
	hue Period
	sat Period
	val Period
	alpha Period
	fixture_set FixtureSet
}

// func NewMode(){

// 	//m = Mode{}
// 	// m.hue := 
// 	// m.sat := 
// 	// m.val := 
// 	// m.alpha := 
// 	// m.fixture_set

// 	return m
// }

func scale(old_min, old_max, new_min, new_max, value float64) float64{
    return ((value - old_min) / (old_max - old_min) ) * (new_max - new_min) + new_min
}

func main() {
    
    ip := "127.0.0.1"
    port := "7770"

    nanokontrol2 := nanokontrol2.Initialize()    
    osc_buffer := make([]byte, 512)

    address := osc.NewAddress(&ip, &port)

    var hue_index float64
    hue_index = 0

    var value_index float64
    value_index = 0

	par1 := Fixture{osc_buffer, 4, 5, 6, 7, 3}
	par2 := Fixture{osc_buffer, 14, 15, 16, 17, 13}
	par3 := Fixture{osc_buffer, 24, 25, 26, 27, 23}
	par4 := Fixture{osc_buffer, 34, 35, 36, 37, 33}

	// kitchen_1 := Fixture{osc_buffer, 65, 64, 66, 0, 0}
	// kitchen_2 := Fixture{osc_buffer, 75, 74, 76, 0, 0}

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
        fmt.Println(uint8(color.R * 255), uint8(color.G * 255), uint8(color.B * 255))

        par1.Set(color.R, color.G, color.B)
        par2.Set(color.R, color.G, color.B)
        par3.Set(color.R, color.G, color.B)
        par4.Set(color.R, color.G, color.B)

        message := make(osc.Message, 0)
        message = append(message, osc.Blob(osc_buffer))
        message.Send(address, "/dmx/universe/1")
     	
        //DMX has a 44Hz max refresh rate
        time.Sleep(time.Second / 44)
    }
}