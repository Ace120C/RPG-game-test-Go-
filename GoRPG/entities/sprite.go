package entities

import "github.com/hajimehoshi/ebiten/v2"

// base struct for all moving drawn entities
type Sprite struct{
  Img  *ebiten.Image
  X, Y float64
}
