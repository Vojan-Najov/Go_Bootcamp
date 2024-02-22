package main

import (
  "os"
  "fmt"
  "io/fs"
  "myFind/settings"
  "path/filepath"
)

func main() {
  settings := settings.New(os.Args[1:])
  if settings == nil {
    os.Exit(1)
  }

  err := filepath.WalkDir(settings.Dirname,
           func(path string, entry fs.DirEntry, err error) error {
    if err != nil {
      return nil
    }

    fullpath := fmt.Sprintf("%s/%s", settings.Dirname, path)

    if entry.IsDir() {
      if settings.PrintDirectories && path != "." {
        fmt.Println(fullpath)
      }
    } else if entry.Type() & fs.ModeSymlink != 0 {
      if settings.PrintSymlinks {
        path, err := filepath.EvalSymlinks(fullpath)
        if err != nil {
          path = "[broken]"
        }
        fmt.Printf("%s -> %s\n", fullpath, path)
      }
    } else {
      if settings.PrintFilenames {
        fmt.Println(fullpath)
      }
    }
    return nil
  })

  if err != nil {
    fmt.Fprintln(os.Stderr, err)
  }
}

