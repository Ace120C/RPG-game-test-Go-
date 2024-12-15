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

//Game Struct

type Game struct{
  player      *entities.Player
  enemies     []*entities.Enemy
  potions     []*entities.Potion
  tilemapJSON *TilemapJSON
  tilesets    []Tileset
  tilemapImg  *ebiten.Image
  cam         *Camera
}

//Normalizing Vectors function, which means you character won't go faster when going diagonally than going straight in any direction

func norm(a, b float64) (float64, float64)  {
  h := math.Hypot(a, b)
  if h == 0{
    return 0, 0
  }
  return a/h, b/h
  
}

//Game update function, this executes every frame

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

  //setting up the max speed of the player, editing this will judge how fast or slow they can go
  const speed = 2
  
  g.player.X += dX * speed
  g.player.Y += dY * speed


  //this is loop is for the enemy sprites and it's for following the player while also applying normalization
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

  //this is for the camera to follow the player with an offset of 8px so it doesn't go out of bound
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

  
//background color for the "skybox"
func (g *Game) Draw(screen *ebiten.Image) {
  //background color (default: sky blue)
  screen.Fill(color.RGBA{120, 180, 255, 255})
  opts := ebiten.DrawImageOptions{}

  //loop over tilemap layers, also id is just the position on where the slice of that tileset is at, it starts from left to right with the position being id: 0
  for layerIndex, layer := range g.tilemapJSON.Layers {
    for index, id := range layer.Data {

      if id == 0 {
        continue
      }


      x := index % layer.Width
      y := index / layer.Width
      //tilemaps are 16x16 pixels so what this does is basically cropping the whole image into a 16x16 picture
      x *= 16
      y *= 16
      fmt.Println(layer.Data)
      img := g.tilesets[layerIndex].Img(id)

      if img == nil {
        fmt.Printf("Skipping invalid tile with ID: %d at layerIndex: %d\n", id, layerIndex)
        continue
      }
      opts.GeoM.Translate(float64(x), float64(y))
      opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy())+16))
      opts.GeoM.Translate(g.cam.X, g.cam.Y)
      screen.DrawImage(img, &opts)
      // srcX := (id - 1) % 22
      // srcY := (id - 1) / 22
      //
      // srcX *= 16
      // srcY *= 16
      //
      //
      // screen.DrawImage(
      //   g.tilemapImg.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
      //   &opts,
      //   )

      opts.GeoM.Reset()
    }
  }

  // set the translation of our DrawImageOptions to the players position
  opts.GeoM.Translate(g.player.X, g.player.Y)
  opts.GeoM.Translate(g.cam.X, g.cam.Y)
  // Drawing the Player
  screen.DrawImage(
    g.player.Img.SubImage(
      //defining the resolution of the sprite, in this case it's 16x16
      image.Rect(0, 0, 16, 16),
      ).(*ebiten.Image),
      &opts,
    )

  opts.GeoM.Reset()
  //Drawing the enemies
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
  //Drawing the potion
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
//display and resolution stuff
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
  //the in-game resolution (default: 320x240)
	return 320, 240//ebiten.WindowSize() 
}
//the actual window resolution
func main() {
  //default is 640x480
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("RPG Game Test")
  ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
  
  //importing the Assets from the declared folders, straightforward stuff
  PlayerImg, _, err := ebitenutil.NewImageFromFile("Assets/Images/ninja.png")
  //we expect an error here, maybe the file is missing or smth so we handle the errors like this
  if err != nil {
    //we used a fatal error because not importing the assest is pretty bad, so its a serious error
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
  //import the tileset
  tilemapImg, _, err := ebitenutil.NewImageFromFile("./Assets/Tilesets/TilesetFloor.png")
  if err != nil {
    //handle the error
    log.Fatal(err)
  }

 //importing the tilemap, this JSON file tells how the program should place each tile on the screen 
  tilemapJSON, err := NewTilemapJSON("Assets/map/spawn.json")
  if err != nil {
    //handle the error
   log.Fatal(err) 
  }
  
  
  tilesets, err := tilemapJSON.GenTilesets()
  if err != nil {
    log.Fatal(err)
  }

  //this is for the player and where it should be spawned, in this case 50 on the X axis and 50 on the Y axis
  Game := Game {
    player: &entities.Player{
      Sprite: &entities.Sprite{
        Img: PlayerImg,
        X: 50.0,
        Y: 50.0,
      }, 
      Health: 3,
    },

    //same thing for the enemy sprite
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
    tilesets: tilesets,
    //centering the camera on the screen
    cam: NewCamera(0.0, 0.0),
  }

  if err := ebiten.RunGame(&Game); err != nil {
		log.Fatal(err)
	}
}
