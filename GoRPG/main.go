package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
  "math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)


type Sprite struct{
  Img  *ebiten.Image
  X, Y float64
}

type Player struct{
  *Sprite
  Health uint
}

type Enemy struct {
  *Sprite
  FollowsPlayer bool
}

type Potion struct {
  *Sprite
  AmtHeal uint
}

type Game struct{
  player  *Player
  enemies []*Enemy
  potions []*Potion
}


func norm(a, b float64) (float64, float64)  {
  h := math.Hypot(a, b)
  if h == 0{
    return 0, 0
  }
  return a/h, b/h
  
}

func (g *Game) Update() error {
  var dX, dY float64 = 0.0, 0.0
  //react to key presses
  if ebiten.IsKeyPressed(ebiten.KeyRight) {
    dX += 2 
  }

  if ebiten.IsKeyPressed(ebiten.KeyLeft) {
    dX -= 2 
  }

  if ebiten.IsKeyPressed(ebiten.KeyUp) {
    dY -= 2 
  }

  if ebiten.IsKeyPressed(ebiten.KeyDown) {
    dY += 2 
  }
  
  // Normalize the vector
  dX, dY = norm(dX, dY)

  const speed = 2
  
  g.player.X += dX * speed
  g.player.Y += dY * speed

  for _, sprite:= range g.enemies{
    if sprite.FollowsPlayer {
      dX, dY := g.player.X-sprite.X, g.player.Y-sprite.Y
      dX, dY = norm(dX, dY)
      const enemySpeed = 1.5
      sprite.X += dX * enemySpeed
      sprite.Y += dY * enemySpeed
      }
    }
for _, potion := range g.potions {
    if g.player.X > potion.X {
      g.player.Health += potion.AmtHeal
      fmt.Printf("Picked up potion! health: %d\n", g.player.Health)
    }
  } 
	return nil
}

  

func (g *Game) Draw(screen *ebiten.Image) {

  screen.Fill(color.RGBA{120, 180, 255, 255})

  opts := ebiten.DrawImageOptions{}
  opts.GeoM.Translate(g.player.X, g.player.Y)
  // draw the Player
  screen.DrawImage(
    g.player.Img.SubImage(
      image.Rect(0, 0, 16, 16),
      ).(*ebiten.Image),
      &opts,
    )

  opts.GeoM.Reset()

  for _, sprite := range g.enemies{
    opts.GeoM.Translate(sprite.X, sprite.Y)

    screen.DrawImage(
      sprite.Img.SubImage(
      image.Rect(0, 0, 16, 16),
      ).(*ebiten.Image),
      &opts,

    )

    opts.GeoM.Reset()
  }

  for _, sprite := range g.potions{
    opts.GeoM.Translate(sprite.X, sprite.Y)

    screen.DrawImage(
      sprite.Img.SubImage(
      image.Rect(0, 0, 16, 16),
      ).(*ebiten.Image),
      &opts,

    )

    opts.GeoM.Reset()
  }


	//ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240//ebiten.WindowSize() 
}

func main() {
  //default is 640x480
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("RPG Game Test")
  ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

  PlayerImg, _, err := ebitenutil.NewImageFromFile("/home/ace/Golang/GoRPG/Assets/Images/ninja.png")
  if err != nil {
    //handle the error
    log.Fatal(err)
  }
  SkeletonImg, _, err := ebitenutil.NewImageFromFile("/home/ace/Golang/GoRPG/Assets/Images/skeleton.png")
  if err != nil {
    //handle the error
    log.Fatal(err)
  }
  potionImg, _, err := ebitenutil.NewImageFromFile("/home/ace/Golang/GoRPG/Assets/Images/potion.png")
  if err != nil {
    //handle the error
    log.Fatal(err)
  }

  
  Game := Game {
    player: &Player{
      Sprite: &Sprite{
        Img: PlayerImg,
        X: 50.0,
        Y: 50.0,
      }, 
      Health: 3,
    },
    enemies: []*Enemy {
      {
        &Sprite{
         Img: SkeletonImg,
          X: 100.0,
          Y: 100.0,
        },      
        true,
      },
      {
        &Sprite{
          Img: SkeletonImg,
          X: 50.0,
          Y: 150.0,
        },      
        false,
      },
    },
    potions: []*Potion {
      {
        &Sprite{
          Img: potionImg,
          X: 210.0,
          Y: 50.0,
        },
        1.0,
      },
    },
  }

  if err := ebiten.RunGame(&Game); err != nil {
		log.Fatal(err)
	}
}
