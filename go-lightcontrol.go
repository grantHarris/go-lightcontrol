package main

import (
 //"math/rand"
 "time"
 "math"
 "fmt"
 //"log"
)
import "github.com/grantHarris/go-nanokontrol2"

//import "github.com/rakyll/portmidi"
import "github.com/lucasb-eyer/go-colorful"
import "github.com/mkb218/go-osc/lib"

// type Fixture struct {
//     b_addr uint8
//     r_addr uint8
//     g_addr uint8
//     b_addr uint8
//     w_addr uint8
// }


// func midiStream() *portmidi.Stream{
//     var stream *portmidi.Stream
//     portmidi.Initialize()
//     if portmidi.CountDevices() == 0{
//         log.Printf("No MIDI controller found")
//     }else{
//         midistream, err := portmidi.NewInputStream(portmidi.DefaultInputDeviceID(), 1024)
//         if err != nil {
//             log.Fatal(err)
//         }
//         stream = midistream
//     }
//     return stream
// }

func makeTimestamp() int64 {
    return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}

func scale(old_min, old_max, new_min, new_max, value float64) float64{
    return ((value - old_min) / (old_max - old_min) ) * (new_max - new_min) + new_min
}

func main() {

    //in := midiStream()
    n := nanokontrol2.Initialize()

    
    period := math.Cos(2 * math.Pi) / 2000
    b := make([]byte, 512)
    
    ip := "127.0.0.1"
    port := "7770"
    address := osc.NewAddress(&ip, &port)

     for{
        // if in != nil{
        //     result, err := in.Poll()
        //     if err != nil {
        //         log.Fatal(err)
        //     }

        //     if result {
        //         msg, err := in.Read(1024)
        //         if err != nil {
        //             log.Fatal(err)
        //         }

        //         for b := range msg {
        //             event := msg[b]
        //             fmt.Println(event.Data1, event.Data2)
        //         }
        //     }
        // }
        fmt.Println(n.Get(16))
        wave := math.Cos(period * float64(makeTimestamp()))
        hue := scale(-1, 1, 0, 360, wave)
        color := colorful.Hsv(hue, 1, 1)
     
        //fmt.Println(uint8(color.R * 255), uint8(color.G * 255), uint8(color.B * 255))
        
        b[0] = byte(color.R * 255)
        b[1] = byte(color.G * 255)
        b[2] = byte(color.B * 255)

        message := make(osc.Message, 0)
        message = append(message, osc.Blob(b))

        message.Send(address, "/dmx/universe/1")
        
        //DMX has a 44Hz max refresh rate
        time.Sleep(time.Second / 44)
    }
}