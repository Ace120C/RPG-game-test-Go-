package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{
  PlayerImage *ebiten.Image
  X, Y float64
}

func (g *Game) Update() error {
  //react to key presses
  if ebiten.IsKeyPressed(ebiten.KeyRight) {
    g.X += 2 
  }

  if ebiten.IsKeyPressed(ebiten.KeyLeft) {
    g.X -= 2 
  }

  if ebiten.IsKeyPressed(ebiten.KeyUp) {
    g.Y -= 2 
  }

  if ebiten.IsKeyPressed(ebiten.KeyDown) {
    g.Y += 2 
  }


	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

  screen.Fill(color.RGBA{120, 180, 255, 255})

  opts := ebiten.DrawImageOptions{}
  opts.GeoM.Translate(g.X, g.Y)
  // draw the Player
  screen.DrawImage(
    g.PlayerImage.SubImage(
      image.Rect(0, 0, 16, 16),
      ).(*ebiten.Image),
      &opts,
    )

	//ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240//ebiten.WindowSize() 
}

func main() {
  //default is 640x480
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
  ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

  PlayerImg, _, err := ebitenutil.NewImageFromFile("/home/ace/Golang/GoRPG/Assets/Images/ninja.png")
  if err != nil {
    //handle the error
    log.Fatal(err)
  }
  if err := ebiten.RunGame(&Game{PlayerImage: PlayerImg, X: 100, Y: 100}); err != nil {
		log.Fatal(err)
	}
}
