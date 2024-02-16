package main

import (
  "os"
  "fmt"
  "flag"
  "bufio"
  "errors"
  "strings"
)

type OldSnapshot struct {
  Name string
}

func (s *OldSnapshot) String() string {
  return s.Name
}

func (s *OldSnapshot) Set(value string) error {
  if s.Name != "" {
    return errors.New("Only one old snapshot is expected")
  }
  s.Name = value
  return nil
}

type NewSnapshot struct {
  Name string
}

func (s *NewSnapshot) String() string {
  return s.Name
}

func (s *NewSnapshot) Set(value string) error {
  if s.Name != "" {
    return errors.New("Only one new snapshot is expected")
  }
  s.Name = value
  return nil
}

var oldSnapshot OldSnapshot
var newSnapshot NewSnapshot

func init() {
  flag.Var(&oldSnapshot, "old", "A string. Set old snapshot filename")
  flag.Var(&newSnapshot, "new", "A string. Set new snapshot filename")
}

func printDifference(oldSnap, newSnap string) error {
  filepathes := make(map[string]int8)

  file, err := os.Open(oldSnap)
  if err != nil {
    return err
  }
  defer file.Close()

  fileScanner := bufio.NewScanner(file)
  for fileScanner.Scan() {
    filepathes[strings.TrimSpace(fileScanner.Text())] = -1
  }

  file, err = os.Open(newSnap)
  if err != nil {
    return err
  }
  defer file.Close()

  fileScanner = bufio.NewScanner(file)
  for fileScanner.Scan() {
    str := strings.TrimSpace(fileScanner.Text())
    if filepathes[str] == 0 {
      filepathes[str] = 1
    } else {
      delete(filepathes, str)
    }
  }

  for k, v := range filepathes {
    if v > 0 {
      fmt.Printf("ADDED %s\n", k)
    } else if v < 0 {
      fmt.Printf("REMOVED %s\n", k)
    }
  }

  return nil
}

func main() {
  flag.Parse()
  if flag.NArg() != 0 {
    fmt.Fprintln(os.Stderr,
                 "No argumets are expected except old and new snapshots")
    return
  } else if oldSnapshot.Name == "" {
    fmt.Fprintln(os.Stderr, "Expected old snapshot")
    flag.PrintDefaults()
    return
  } else if newSnapshot.Name == "" {
    fmt.Fprintln(os.Stderr, "Expected new snapshot")
    flag.PrintDefaults()
    return
  }

  err := printDifference(oldSnapshot.Name, newSnapshot.Name)
  if err != nil {
    fmt.Fprintln(os.Stderr, err)
  }
}
