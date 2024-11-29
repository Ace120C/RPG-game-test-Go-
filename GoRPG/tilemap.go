package main

import (
	"encoding/json"
	"os"
)

type TilemapplayerJSON struct {
  Data []int  `json:"data"`
  Width int   `json:"width"`
  Height int  `json:"height"`
}

type TilemapJSON struct {
  Layers []TilemapplayerJSON `json:"layers"`
}

func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
  contents, err := os.ReadFile(filepath)
  if err != nil {
    return nil, err
  }
  var tilemapJSON TilemapJSON
  err = json.Unmarshal(contents, &tilemapJSON)
  if err != nil {
    return nil, err
  }

  return &tilemapJSON, nil
}
