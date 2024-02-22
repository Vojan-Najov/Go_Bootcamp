package main

import (
	"bufio"
	"fmt"
	"myWc/settings"
	"os"
	"sync"
)

func main() {
	settings := settings.New(os.Args[1:])
	if settings == nil {
		os.Exit(1)
	}

	split := bufio.ScanLines
	if settings.CountCharacters {
		split = bufio.ScanRunes
	} else if settings.CountWords {
		split = bufio.ScanWords
	}

	var wg sync.WaitGroup
	for _, filename := range settings.Filenames {
		file := os.Stdin
		if filename != "-" {
			var err error
			file, err = os.Open(filename)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			defer file.Close()
		} else {
			filename = "stdin"
		}

		fileScanner := bufio.NewScanner(file)
		fileScanner.Split(split)

		wg.Add(1)
		go wc(&wg, filename, fileScanner)
	}
	wg.Wait()
}

func wc(wg *sync.WaitGroup, filename string, scanner *bufio.Scanner) {
	var count uint
	for scanner.Scan() {
		count++
	}
	fmt.Printf("%d\t%s\n", count, filename)
	wg.Done()
}
