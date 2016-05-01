package goperiodicity

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
