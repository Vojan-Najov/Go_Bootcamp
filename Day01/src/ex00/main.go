package main

import (
  "os"
  "fmt"
  "flag"
  "strings"
  "errors"
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

func main() {
  flag.Parse()
  if flag.NArg() != 0 {
    fmt.Fprintln(os.Stderr,
                 "No arguments are expected except for the -f option")
    flag.PrintDefaults()
    return 
  } else if len(filenameFlag) == 0 {
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
