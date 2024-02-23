package main

import (
  "errors"
  "os"
  "fmt"
  "flag"
)

type destdir struct {
  path string
}

func (d *destdir) String() string {
  return d.path
}

func (d *destdir) Set(value string) error {
  if len(d.path) > 0 {
    return errors.New("Only one destinetion dir is expected")
  }
  
  fileInfo, err := os.Stat(value)
  if err != nil {
    return err
  }

  if !fileInfo.IsDir() {
    return errors.New("Expected path for the directory")
  }

  d.path = value
  return nil
}

var dir destdir

func init() {
  flag.Var(&dir, "a", "A string. Set path of destination directory")
}

func main() {
  flag.Parse()
  if flag.NArg() == 0 {
    fmt.Fprintln(os.Stderr, "Expected log filenames as argumets")
    flag.PrintDefaults()
    os.Exit(1)
  }

  filenames := flag.Args()
  if len(dir.path) == 0 {
    dir.path = "."
  }

  fmt.Printf("_%s_\n", dir.path)
  fmt.Println(filenames)
}
