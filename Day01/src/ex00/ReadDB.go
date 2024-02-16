package main

import (
  "os"
  "fmt"
  "flag"
  "strings"
  "errors"
  "encoding/xml"
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
  data, err := os.ReadFile(reader.Filename)
  if err != nil {
    return nil, err
  }

  var cookbook CookBook
  if err = xml.Unmarshal(data, &cookbook); err != nil {
    return nil, err
  }
  
  return &cookbook, nil  
}

type DBWriter interface {
  Write(cookbook CookBook) error
}

type JSONWriter struct {}

func (writer JSONWriter) Write(cookbook CookBook) error {
  data, err := json.MarshalIndent(cookbook, "", "  ")
  if err != nil {
    return err
  }
  fmt.Println(string(data))
  return nil
}

type XMLWriter struct {}

func (writer XMLWriter) Write(cookbook CookBook) error {
  data, err := xml.MarshalIndent(cookbook, "", "    ")
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
    var writer DBWriter
    if strings.HasSuffix(f, ".json") {
      reader = JSONReader{f}
      writer = XMLWriter{}
    } else {
      reader = XMLReader{f}
      writer = JSONWriter{}
    }

    cookbook, err := reader.Read()
    if err != nil {
      fmt.Fprintf(os.Stderr, "%s: %s\n", f, err)
      return 
    }

    if err = writer.Write(*cookbook); err != nil {
      fmt.Fprintf(os.Stderr, "%s: %s\n", f, err)
      return 
    }
  }
}
