package main

import (
	"GoRPG/entities"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)


type Game struct{
  player      *entities.Player
  enemies     []*entities.Enemy
  potions     []*entities.Potion
  tilemapJSON *TilemapJSON
  tilemapImg  *ebiten.Image
  cam         *Camera
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

  g.cam.FollowTarget(g.player.X+8, g.player.Y+8, 320, 240)
  g.cam.Constrain(
    // float64(g.tilemapJSON.Layers[0].Width)*16.0,
    // float64(g.tilemapJSON.Layers[0].Height)*16.0,
    25.0*16.0,
    25.0*16.0,
    320,
    240,
  )

	return nil
}

  

func (g *Game) Draw(screen *ebiten.Image) {

  screen.Fill(color.RGBA{120, 180, 255, 255})
  opts := ebiten.DrawImageOptions{}

  //loop over the layers
  for _, layer := range g.tilemapJSON.Layers {
    for index, id := range layer.Data {
      x := index % layer.Width
      y := index / layer.Width

      x *= 16
      y *= 16

      srcX := (id - 1) % 22
      srcY := (id - 1) / 22

      srcX *= 16
      srcY *= 16

      opts.GeoM.Translate(float64(x), float64(y))
      
      opts.GeoM.Translate(g.cam.X, g.cam.Y)

      screen.DrawImage(
        g.tilemapImg.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
        &opts,
        )

      opts.GeoM.Reset()
    }
  }

  // set the translation of our DrawImageOptions to the players position
  opts.GeoM.Translate(g.player.X, g.player.Y)
  opts.GeoM.Translate(g.cam.X, g.cam.Y)
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
    opts.GeoM.Translate(g.cam.X, g.cam.Y)

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
    opts.GeoM.Translate(g.cam.X, g.cam.Y)

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

  PlayerImg, _, err := ebitenutil.NewImageFromFile("Assets/Images/ninja.png")
  if err != nil {
    //handle the error
    log.Fatal(err)
  }
  SkeletonImg, _, err := ebitenutil.NewImageFromFile("Assets/Images/skeleton.png")
  if err != nil {
    //handle the error
    log.Fatal(err)
  }
  potionImg, _, err := ebitenutil.NewImageFromFile("Assets/Images/potion.png")
  if err != nil {
    //handle the error
    log.Fatal(err)
  }

  tilemapImg, _, err := ebitenutil.NewImageFromFile("Assets/Tilesets/TilesetFloor.png")
  if err != nil {
    //handle the error
    log.Fatal(err)
  }

  
  tilemapJSON, err := NewTilemapJSON("Assets/map/spawn.json")
  if err != nil {
   log.Fatal(err) 
  }
  
  Game := Game {
    player: &entities.Player{
      Sprite: &entities.Sprite{
        Img: PlayerImg,
        X: 50.0,
        Y: 50.0,
      }, 
      Health: 3,
    },
    enemies: []*entities.Enemy {
      {
        Sprite: &entities.Sprite{
         Img: SkeletonImg,
          X: 100.0,
          Y: 100.0,
        },      
        FollowsPlayer: true,
      },
      {
        Sprite: &entities.Sprite{
          Img: SkeletonImg,
          X: 50.0,
          Y: 150.0,
        },      
        FollowsPlayer: false,
      },
    },
    potions: []*entities.Potion {
      {
        Sprite: &entities.Sprite{
          Img: potionImg,
          X: 210.0,
          Y: 50.0,
        },
        AmtHeal: 1.0,
      },
    },
    tilemapJSON: tilemapJSON,
    tilemapImg: tilemapImg,
    cam: NewCamera(0.0, 0.0),
  }

  if err := ebiten.RunGame(&Game); err != nil {
		log.Fatal(err)
	}
}
