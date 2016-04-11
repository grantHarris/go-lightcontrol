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

type Fixture struct {    
    R_addr uint8
    G_addr uint8
    B_addr uint8
    W_addr uint8
    brightness_addr uint8
}

func (f *Fixture) Set(buffer []byte, r float64, g float64, b float64){
	
	if f.brightness_addr != 0{
		buffer[f.brightness_addr] = 255
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
		buffer[f.W_addr] = byte(common * 255)
	}

	buffer[f.R_addr] = byte(r * 255)
	buffer[f.G_addr] = byte(g * 255)
	buffer[f.B_addr] = byte(b * 255)
}


func scale(old_min, old_max, new_min, new_max, value float64) float64{
    return ((value - old_min) / (old_max - old_min) ) * (new_max - new_min) + new_min
}

func main() {
	kitchen_right := Fixture{64, 63, 65, 0, 0}
	kitchen_left := Fixture{74, 73, 75, 0, 0}

    n := nanokontrol2.Initialize()    
    b := make([]byte, 512)
    
    ip := "127.0.0.1"
    port := "7770"

    address := osc.NewAddress(&ip, &port)

    var counter float64
    counter = 0

    var ac float64
    ac = 0

     for{
     	//Hue period
     	per := n.Get(16)
		counter = counter + (per*math.Pi/220)

     	//Hue size
     	width := n.Get(0)
     	width = 360 * width

     	//Alpha period
     	alpha := n.Get(17)
		ac = ac + (alpha*math.Pi/220)
     	
     	//Alpha low
     	alpha_low := n.Get(1)
     	//Alpha high
     	alpha_high := n.Get(2)
	     	
        wave := math.Sin(counter) + math.Sin(counter/2)
        hue := scale(-1, 1, 0, width, wave)
        
        wave_value := math.Sin(ac)
        value := scale(-1, 1, alpha_low, alpha_high, wave_value)

        color := colorful.Hsv(hue, 1, value)
        fmt.Println(uint8(color.R * 255), uint8(color.G * 255), uint8(color.B * 255))

        kitchen_right.Set(b, color.R, color.G, color.B)
        kitchen_left.Set(b, color.R, color.G, color.B)

        message := make(osc.Message, 0)
        message = append(message, osc.Blob(b))
        message.Send(address, "/dmx/universe/1")
     	
        //DMX has a 44Hz max refresh rate
        time.Sleep(time.Second / 44)
    }
}