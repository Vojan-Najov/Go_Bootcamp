package settings

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type Settings struct {
	Filenames       []string
	CountLines      bool
	CountCharacters bool
	CountWords      bool
}

func (s *Settings) handleL(str string) error {
	if str != "true" {
		return errors.New("incorrect value")
	}
	s.CountLines = true
	return nil
}

func (s *Settings) handleW(str string) error {
	if str != "true" {
		return errors.New("incorrect value")
	}
	s.CountWords = true
	return nil
}

func (s *Settings) handleM(str string) error {
	if str != "true" {
		return errors.New("incorrect value")
	}
	s.CountCharacters = true
	return nil
}

func New(args []string) *Settings {
	var stngs Settings

	flagset := flag.NewFlagSet("myWc", flag.ContinueOnError)
	flagset.SetOutput(os.Stderr)

	flagset.BoolFunc("l", "Count lines", stngs.handleL)
	flagset.BoolFunc("m", "Count characters", stngs.handleM)
	flagset.BoolFunc("w", "Count words", stngs.handleW)

	if err := flagset.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "%v\n", err)
		return nil
	}

	if stngs.CountLines && (stngs.CountWords || stngs.CountCharacters) ||
		stngs.CountWords && stngs.CountCharacters {
		fmt.Fprintln(os.Stderr, "A single flag is expected: -l, -w or -m")
		flagset.PrintDefaults()
		return nil
	}

	if !(stngs.CountLines || stngs.CountWords || stngs.CountCharacters) {
		stngs.CountWords = true
	}

	if flagset.NArg() == 0 {
		stngs.Filenames = []string{"-"}
	} else {
		stngs.Filenames = flagset.Args()
	}

	return &stngs
}
