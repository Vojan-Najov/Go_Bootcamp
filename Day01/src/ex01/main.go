package main

import (
  "os"
  "fmt"
  "flag"
  "strings"
  "errors"
)

type JSONFile struct {
  Name string
}

func (f *JSONFile) String() string {
  return f.Name
}

func (f *JSONFile) Set(value string) error {
  if !strings.HasSuffix(value, ".json") {
    return errors.New("Expected json extension")
  }
  if f.Name != "" {
    return errors.New("Only one file is expected")
  }
  f.Name = value
  return nil
}

type XMLFile struct {
  Name string
}

func (f *XMLFile) String() string {
  return f.Name
}

func (f *XMLFile) Set(value string) error {
  if !strings.HasSuffix(value, ".xml") {
    return errors.New("Expected xml extension")
  }
  if f.Name != "" {
    return errors.New("Only one file is expected")
  }
  f.Name = value
  return nil
}

var xmlFile XMLFile
var jsonFile JSONFile

func init() {
  flag.Var(&xmlFile, "old", "A string. Set xml filename")
  flag.Var(&jsonFile, "new", "A string. Set json filename")
}

func main() {
  flag.Parse()
  if flag.NArg() != 0 {
    fmt.Fprintln(os.Stderr,
                 "No argumets are expected except old and new databases")
    return
  } else if xmlFile.Name == "" {
    fmt.Fprintln(os.Stderr, "Expected xml database")
    flag.PrintDefaults()
    return
  } else if jsonFile.Name == "" {
    fmt.Fprintln(os.Stderr, "Expected json database")
    flag.PrintDefaults()
    return
  }

  var reader DBReader
  reader = XMLReader{xmlFile.Name}
  oldCookbook, err := reader.Read()
  if err != nil {
    fmt.Fprintf(os.Stderr, "%s: %s\n", xmlFile.Name, err)
    return
  }
  reader = JSONReader{jsonFile.Name}
  newCookbook, err := reader.Read()
  if err != nil {
    fmt.Fprintf(os.Stderr, "%s: %s\n", jsonFile.Name, err)
    return
  }

  PrintCakeDifference(oldCookbook, newCookbook)
}
