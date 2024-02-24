package main

import (
	"fmt"
	"io/fs"
	"myFind/settings"
	"os"
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
				if path == settings.Dirname {
					return err
				}
				return nil
			}

			if entry.IsDir() {
				printDirname(path, settings)
			} else if entry.Type()&fs.ModeSymlink != 0 {
				printSymlink(path, settings)
			} else {
				printFilename(path, settings)
			}
			return nil
		})

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func printDirname(path string, stngs *settings.Settings) {
	if stngs.PrintDirectories && path != stngs.Dirname {
		fmt.Println(path)
	}
}

func printSymlink(path string, stngs *settings.Settings) {
	if stngs.PrintSymlinks {
		pathlink, err := filepath.EvalSymlinks(path)
		if err != nil {
			pathlink = "[broken]"
		}
		fmt.Printf("%s -> %s\n", path, pathlink)
	}
}

func printFilename(path string, stngs *settings.Settings) {
	if stngs.PrintFilenames {
		if !stngs.OnlyExt || stngs.OnlyExt && filepath.Ext(path) == stngs.Ext {
			fmt.Println(path)
		}
	}
}
