package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type rand struct {
	x, y, z, w uint32
}

func (r *rand) next() uint32 {
	// math/rand is too slow to keep 60 FPS on web browsers.
	// Use Xorshift instead: http://en.wikipedia.org/wiki/Xorshift
	t := r.x ^ (r.x << 11)
	r.x, r.y, r.z = r.y, r.z, r.w
	r.w = (r.w ^ (r.w >> 19)) ^ (t ^ (t >> 8))
	return r.w
}

var theRand = &rand{12345678, 4185243, 776511, 45411}
var chip = &ChipContext{}

type Game struct {
	screen *image.RGBA
}

func (g Game) Update() error {
	return nil
}

func (g Game) Draw(screen *ebiten.Image) {
	gameScreen := ebiten.NewImage(ScreenWidth, ScreenHeight)
	gameScreen.WritePixels(g.screen.Pix)

	gameScale := float64(WindowWidth) / float64(ScreenWidth)
	gameScaleOp := &ebiten.DrawImageOptions{}
	gameScaleOp.GeoM.Scale(gameScale, gameScale)
	gameScaleOp.Filter = ebiten.FilterNearest
	screen.DrawImage(gameScreen, gameScaleOp)

	debugScreen := ebiten.NewImage(ScreenWidth, ScreenHeight)
	debugScreen.Fill(color.Transparent)
	ebitenutil.DebugPrint(debugScreen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
	debugScaleOp := &ebiten.DrawImageOptions{}
	debugScaleOp.GeoM.Scale(5, 5)
	debugScaleOp.Filter = ebiten.FilterNearest
	screen.DrawImage(debugScreen, debugScaleOp)
}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowWidth, WindowHeight
}

func main() {
	chip.FrameBuffer = [ScreenWidth * ScreenHeight]bool{}
	for i := range chip.FrameBuffer {
		chip.FrameBuffer[i] = theRand.next()%2 == 1
	}

	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle("Chip8 Emulator")
	g := &Game{
		screen: image.NewRGBA(image.Rect(0, 0, ScreenWidth, ScreenHeight)),
	}

	const l = ScreenWidth * ScreenHeight
	for i := range l {
		var c uint8
		if chip.FrameBuffer[i] {
			c = 0xff // white for 1
		} else {
			c = 0x00 // black for 0
		}
		g.screen.Pix[4*i] = c      // R
		g.screen.Pix[4*i+1] = c    // G
		g.screen.Pix[4*i+2] = c    // B
		g.screen.Pix[4*i+3] = 0xff // A
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
