package main

import (
	"encoding/json"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Tileset interface {
  Img(id int) *ebiten.Image
}

type UniformTilesetJSON struct {
  Path  string `json:"image"`
}

type UniformTileset struct {
  img *ebiten.Image
  gid int
}

func (u *UniformTileset) Img(id int) *ebiten.Image{

  id -= u.gid

  srcX := (id - 1) % 22
  srcY := (id - 1) / 22

  srcX *= 16
  srcY *= 16
  return u.img.SubImage(
    image.Rect(
      srcX, srcY, srcX+16, srcY+16,
    ),
  ).(*ebiten.Image)
}

type TileJSON struct {
 Id     int   `json:"id"`
 Path   string `json:"image"`
 Width  int     `json:"imagewidth"`
 Height int      `json:"imageheight"`

}

type DynTilesetJSON struct {
 Tiles []*TileJSON `json:"tiles"` 
}

type DynTileset struct {
 imgs[]*ebiten.Image 
 gid   int
}


func (d *DynTileset) Img(id int) *ebiten.Image{
  id -= d.gid

  srcX := (id - 1) % 22
  srcY := (id - 1) / 22

  srcX *= 16
  srcY *= 16
  return d.imgs[id]
}

func NewTileset(path string, gid int) (Tileset, error) {
  contents, err := os.ReadFile(path)
  if err != nil {
    return nil, err
  }

  if strings.Contains(path, "buildings") {
    //return DynTileset
    var dynTilesetJSON DynTilesetJSON
    err = json.Unmarshal(contents, &dynTilesetJSON)
    if err != nil {
      return nil, err
    }
    
    dynTileset := DynTileset{}
    dynTileset.gid = gid
    dynTileset.imgs = make([]*ebiten.Image, 0)

    for _, tileJSON := range dynTilesetJSON.Tiles {
      
      tileJSONPath := tileJSON.Path
      tileJSONPath = filepath.Clean(tileJSONPath)
      tileJSONPath = strings.ReplaceAll(tileJSONPath, "\\", "/")
      tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
      tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
      tileJSONPath = filepath.Join("Assets/", tileJSONPath)
      

      img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
      if err != nil {
        return nil, err
      }
      dynTileset.imgs = append(dynTileset.imgs, img)
    }

    return &dynTileset, nil
  }
    //return UniformTileset 
    var uniformTilesetJSON UniformTilesetJSON
    err = json.Unmarshal(contents, &uniformTilesetJSON)
    if err != nil {
      return nil, err
    }

    uniformTileset := UniformTileset{}


      tileJSONPath := uniformTilesetJSON.Path
      tileJSONPath = filepath.Clean(tileJSONPath)
      tileJSONPath = strings.ReplaceAll(tileJSONPath, "\\", "/")
      tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
      tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
      tileJSONPath = filepath.Join("Assets/", tileJSONPath)

    
    img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
    if err != nil {
      return nil, err
    }
    uniformTileset.img = img
    uniformTileset.gid = gid

    return &uniformTileset, nil
  }
