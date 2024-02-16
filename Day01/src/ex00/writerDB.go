package main

import (
  "fmt"
  "encoding/xml"
  "encoding/json"
)

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
