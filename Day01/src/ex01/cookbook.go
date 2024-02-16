package main

import (
  "fmt"
  "encoding/xml"
)

type Ingredient struct {
  Name  string `xml:"itemname" json:"ingredient_name"`
  Count string `xml:"itemcount" json:"ingredient_count"`
  Unit  string `xml:"itemunit" json:"ingredient_unit,omitempty"`
}

type CakeRecipe struct {
  Name        string       `xml:"name" json:"name"`
  Time        string       `xml:"stovetime" json:"time"`
  Ingredients []Ingredient `xml:"ingredients>item" json:"ingredients"`
}

type CookBook struct {
  XMLName xml.Name     `xml:"recipes" json:"-"`
  Cakes   []CakeRecipe `xml:"cake" json:"cake"`
}

type difference struct {
  count int
  oldTime string
  newTime string
  oldIngredients []Ingredient
  newIngredients []Ingredient
}

type ingredientDifference struct {
  count int
  newCount string
  newUnit string
  oldCount string
  oldUnit string
}

const (
  CakeRemovedFmt = "REMOVED cake \"%s\"\n"
  CakeAddedFmt = "ADDED cake \"%s\"\n"
  TimeChangedFmt = 
    "CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n"
  IngredientRemovedFmt = 
    "REMOVED ingredient \"%s\" for cake  \"%s\"\n"
  IngredientAddedFmt = 
    "ADDED ingredient \"%s\" for cake  \"%s\"\n"
  UnitRemovedFmt =
    "REMOVED unit \"%s\" for ingredient \"%s\" for cake  \"%s\"\n"
  UnitChangedFmt =
    "CHANGED unit for ingredient \"%s\" for cake  \"%s\" - \"%s\" instead of \"%s\"\n"
  CountChangedFmt =
    "CHANGED unit count for ingredient \"%s\" for cake  \"%s\" - \"%s\" instead of \"%s\"\n"
)

func PrintCakeDifference(old, new *CookBook) {
  cakes := make(map[string]*difference)
  for _, cake := range new.Cakes {
    if cakes[cake.Name] == nil {
      dif := difference{
        count: 1, newTime: cake.Time, newIngredients: cake.Ingredients,
      }
      cakes[cake.Name] = &dif
    } else {
      cakes[cake.Name].count++
      cakes[cake.Name].newTime = cake.Time
      cakes[cake.Name].newIngredients = cake.Ingredients
    }
  }
  for _, cake := range old.Cakes {
    if cakes[cake.Name] == nil {
      dif := difference{
        count: -1, oldTime: cake.Time, oldIngredients: cake.Ingredients,
      }
      cakes[cake.Name] = &dif
    } else {
      cakes[cake.Name].count--
      cakes[cake.Name].oldTime = cake.Time
      cakes[cake.Name].oldIngredients = cake.Ingredients

    }
  }
  for k, v := range cakes {
    if v.count < 0 {
      fmt.Printf(CakeRemovedFmt, k)
    } else if v.count > 0 {
      fmt.Printf(CakeAddedFmt, k)
    }
    if v.count != 0 {
      delete(cakes, k)
    }
  }
  for k, v := range cakes {
    if v.oldTime != v.newTime {
      fmt.Printf(TimeChangedFmt, k, v.newTime, v.oldTime)
    }
    PrintIngredientDifference(v.oldIngredients, v.newIngredients, k)
  }
}

func PrintIngredientDifference(old, new []Ingredient, cake string) {
  ingredients := make(map[string]*ingredientDifference)
  for _, ing := range new {
    if ingredients[ing.Name] == nil {
      dif := ingredientDifference{
        count: 1, newCount: ing.Count, newUnit: ing.Unit,
      }
      ingredients[ing.Name] = &dif
    }
  }
  for _, ing := range old {
    if ingredients[ing.Name] == nil {
      dif := ingredientDifference{
        count: -1, oldCount: ing.Count, oldUnit: ing.Unit,
      }
      ingredients[ing.Name] = &dif
    } else {
      ingredients[ing.Name].count--
      ingredients[ing.Name].oldCount = ing.Count
      ingredients[ing.Name].oldUnit = ing.Unit
    }
  }
  for k, v := range ingredients {
    if v.count < 0 {
      fmt.Printf(IngredientRemovedFmt, k, cake)
    } else if v.count > 0 {
      fmt.Printf(IngredientAddedFmt, k, cake)
    }
    if v.count != 0 {
      delete(ingredients, k)
    }
  }
  for k, v := range ingredients {
    if v.newUnit == "" && v.oldUnit != "" {
      fmt.Printf(UnitRemovedFmt, v.oldUnit, k, cake)
    } else if v.newUnit != v.oldUnit {
      fmt.Printf(UnitChangedFmt, k, cake, v.newUnit,  v.oldUnit)
    } else if v.newCount != v.oldCount {
      fmt.Printf(UnitChangedFmt, k, cake, v.newCount,  v.oldCount)
    }
  }
}
