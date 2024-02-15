package main

import (
  "os"
  "fmt"
  "flag"
  "strings"
  "errors"
  "encoding/json"
)

type filename []string

func (f *filename) String() string {
  return fmt.Sprint(*f)
}

func (f *filename) Set(value string) error {
  if !strings.HasSuffix(value, ".xml") &&
     !strings.HasSuffix(value, ".json") {
    return errors.New("Expected xml or json extension")
  }
  *f = append(*f, value)
  return nil
}

var filenameFlag filename

func init() {
  flag.Var(&filenameFlag, "f",
           "A string. Set filename json or xml extension for DB")
}

type Ingredient struct {
  Name string  `json:"ingredient_name"`
  Count string `json:"ingredient_count"`
  Unit string  `json:"ingredient_unit,omitempty"`
}

type CakeRecipe struct {
  Name string              `json:"name"`
  Time string              `json:"time"`
  Ingredients []Ingredient `json:"ingredients"`
}

type CookBook struct {
  Cakes []CakeRecipe `json:"cake"`
}

type DBReader interface {
  Read() (*CookBook, error)
}

type JSONReader struct {
  Filename string
}

func (reader JSONReader) Read() (*CookBook, error) {
  data, err := os.ReadFile(reader.Filename)
  if err != nil {
    return nil, err
  }

  var cookbook CookBook
  if err = json.Unmarshal(data, &cookbook); err != nil {
    return nil, err
  }
  
  return &cookbook, nil  
}

type XMLReader struct {
  Filename string
}

func (reader XMLReader) Read() (*CookBook, error) {
  return nil, nil
}

type DBWriter interface {
  Write(cookbook CookBook) error
}

type JSONWriter struct {}

func (writer JSONWriter) Write(cookbook CookBook) error {
  data, err := json.MarshalIndent(cookbook, "", "    ")
  if err != nil {
    return err
  }
  fmt.Println(string(data))
  return nil
}

func main() {
  flag.Parse()
  if flag.NArg() != 0 {
    fmt.Fprintln(os.Stderr,
                 "No arguments are expected except for the -f option")
    flag.PrintDefaults()
    return 
  }

  for _, f := range filenameFlag {
    var reader DBReader
    if strings.HasSuffix(f, ".json") {
      reader = JSONReader{f}
    } else {
      reader = XMLReader{f}
    }
    cookbook, err := reader.Read()
    if err != nil {
      fmt.Fprintf(os.Stderr, "%s: %s\n", f, err)
      return 
    }

    if cookbook != nil {
      fmt.Println(cookbook)
    }

    writer := JSONWriter{}

    writer.Write(*cookbook)

    //var writer DBWriter
    //if strings.HasSuffix(f, ".json") {
    //  writer = XMLWriter{cookbook}
    //} else {
    //  writer = JSONWriter{cookbook}
    //}
    //writer.Write()
  }
}
