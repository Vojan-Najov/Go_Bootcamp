package main

import (
  "encoding/xml"
)

type Ingredient struct {
  Name  string `xml:"itemname" json:"ingredient_name"`
  Count string `xml:"itemcount" json:"ingredient_count"`
  Unit  string `xml:"itemunit" json:"ingredient_unit,omitempty"`
}

type CakeRecipe struct {
  Name         string      `xml:"name" json:"name"`
  Time         string      `xml:"stovetime" json:"time"`
  Ingredients []Ingredient `xml:"ingredients>item" json:"ingredients"`
}

type CookBook struct {
  XMLName xml.Name     `xml:"recipes" json:"-"`
  Cakes   []CakeRecipe `xml:"cake" json:"cake"`
}
