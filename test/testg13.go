package main

import (
	"image/color"
	"log"

	"github.com/bjanders/g13"
)

func main() {

	g, err := g13.NewG13()
	if err != nil {
		log.Fatalln(err)
	}
	keyCh, stickCh, _ := g.WatchKeys()
	var keyState g13.KeyState
	var stickState g13.StickState
	//var rgb [3]byte
	var color color.RGBA
	//g.SetMode(4, 255)
	// for i := 0; i <= 255; i++ {
	// 	log.Print(i)
	// 	g.SetMode(4, byte(i))
	// 	time.Sleep(50 * time.Millisecond)
	// }
	for {
		select {
		case stickState = <-stickCh:
			log.Printf("Stick: %d:%d", stickState.X, stickState.Y)
			color.R = uint8(stickState.X)
			color.A = uint8(stickState.Y)
			g.SetColor(color)
		case keyState = <-keyCh:
			if keyState.Key == g13.KeyLeft && color.B > 0 {
				color.B--
				g.SetColor(color)

			} else if keyState.Key == g13.KeyRight && color.B < 255 {
				color.B++
				g.SetColor(&color)

			}
			log.Printf("%s = %v", g13.Keys[keyState.Key], keyState.Down)

		}
	}
}
