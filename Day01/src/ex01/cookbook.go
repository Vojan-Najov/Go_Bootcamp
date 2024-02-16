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
  newItemCount string
  newItemUnit string
  oldItemCount string
  oldItemUnit string
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
  cakeNames := make(map[string]*difference)
  for _, cake := range new.Cakes {
    if cakeNames[cake.Name] == nil {
      dif := difference{count: 1, newTime: cake.Time, newIngredients: cake.Ingredients}
      cakeNames[cake.Name] = &dif
    } else {
      cakeNames[cake.Name].count++
      cakeNames[cake.Name].newTime = cake.Time
      cakeNames[cake.Name].newIngredients = cake.Ingredients
    }
  }
  for _, cake := range old.Cakes {
    if cakeNames[cake.Name] == nil {
      dif := difference{count: -1, oldTime: cake.Time, oldIngredients: cake.Ingredients}
      cakeNames[cake.Name] = &dif
    } else {
      cakeNames[cake.Name].count--
      cakeNames[cake.Name].oldTime = cake.Time
      cakeNames[cake.Name].oldIngredients = cake.Ingredients

    }
  }
  for k, v := range cakeNames {
    if v.count < 0 {
      fmt.Printf(CakeRemovedFmt, k)
    } else if v.count > 0 {
      fmt.Printf(CakeAddedFmt, k)
    }
    if v.count != 0 {
      delete(cakeNames, k)
    }
  }
  for k, v := range cakeNames {
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
        count: 1, newItemCount: ing.Count, newItemUnit: ing.Unit,
      }
      ingredients[ing.Name] = &dif
    }
  }
  for _, ing := range old {
    if ingredients[ing.Name] == nil {
      dif := ingredientDifference{
        count: -1, oldItemCount: ing.Count, oldItemUnit: ing.Unit,
      }
      ingredients[ing.Name] = &dif
    } else {
      ingredients[ing.Name].count--
      ingredients[ing.Name].oldItemCount = ing.Count
      ingredients[ing.Name].oldItemUnit = ing.Unit
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
    if v.newItemUnit == "" && v.oldItemUnit != "" {
      fmt.Printf(UnitRemovedFmt, v.oldItemUnit, k, cake)
    } else if v.newItemUnit != v.oldItemUnit {
      fmt.Printf(UnitChangedFmt, k, cake, v.newItemUnit,  v.oldItemUnit)
    } else if v.newItemCount != v.oldItemCount {
      fmt.Printf(UnitChangedFmt, k, cake, v.newItemCount,  v.oldItemCount)
    }
  }
}
