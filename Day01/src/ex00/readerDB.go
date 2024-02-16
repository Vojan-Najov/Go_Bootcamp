package main

import (
  "os"
  "encoding/xml"
  "encoding/json"
)

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
