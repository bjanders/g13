package main

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"log"
	"os"

	"github.com/bjanders/g13"
)

func main() {

	g, err := g13.NewG13()
	if err != nil {
		log.Fatalln(err)
	}
	var color color.RGBA
	// g.SetMLEDs(g13.LEDM1)
	// time.Sleep(500 * time.Millisecond)
	// g.SetMLEDs(g13.LEDM1 | g13.LEDM2)
	// time.Sleep(500 * time.Millisecond)
	// g.SetMLEDs(g13.LEDM2)
	// time.Sleep(500 * time.Millisecond)
	// g.SetMLEDs(g13.LEDM3)
	// time.Sleep(500 * time.Millisecond)
	// g.SetMLEDs(g13.LEDMR)

	// for x := 0; x < g13.LCDSizeX; x++ {
	// 	for y := 0; y < g13.LCDSizeY; y++ {
	// 		g.SetPixel(x, y)
	// 	}
	// }
	g.ClearLCD()
	g.DrawLCD()
	file, err := os.Open("pkg.png")
	if err != nil {
		log.Fatal(err)
	}
	img, fmt, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	log.Printf("Format is %s", fmt)
	//for y := 0; y < 100; y++ {
	draw.Draw(g.LCD, g.LCD.Bounds(), img, image.Pt(0, 50), draw.Src)
	g.AddString("Hello", 30, 20)
	g.DrawLCD()

	// time.Sleep(10 * time.Millisecond)
	//}
	for {
		select {
		case stickState := <-g.StickCh:
			log.Printf("Stick: %d:%d", stickState.X, stickState.Y)
			// color.R = uint8(stickState.X)
			// color.A = uint8(stickState.Y)
			// g.SetColor(color)
			//g.SetLCDPixel(stickState.X, stickState.Y)
			//g.DrawLCD2()
		case keyState := <-g.KeyCh:
			if keyState.Key == g13.KeyLeft && color.B > 0 {
				color.B--
				g.SetColor(color)

			} else if keyState.Key == g13.KeyRight && color.B < 255 {
				color.B++
				g.SetColor(color)

			}
			log.Printf("%s = %v", g13.Keys[keyState.Key], keyState.Down)
		case backLight := <-g.BacklightCh:
			log.Printf("Backlight: %v", backLight)

		}
	}
}
