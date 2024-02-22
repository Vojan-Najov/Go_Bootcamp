package settings

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Settings struct {
	Dirname          string
	PrintDirectories bool
	PrintSymlinks    bool
	PrintFilenames   bool
	OnlyExt          bool
	Ext              string
}

func (s *Settings) handleF(str string) error {
	if str != "true" {
		return errors.New("incorrect value")
	}
	s.PrintFilenames = true
	return nil
}

func (s *Settings) handleD(str string) error {
	if str != "true" {
		return errors.New("incorrect value")
	}
	s.PrintDirectories = true
	return nil
}

func (s *Settings) handleSl(str string) error {
	if str != "true" {
		return errors.New("incorrect value")
	}
	s.PrintSymlinks = true
	return nil
}

func (s *Settings) handleExt(str string) error {
	if !s.PrintFilenames {
		return errors.New("works ONLY when -f is specified")
	} else if len(strings.TrimSpace(str)) == 0 {
		return errors.New("incorrect value")
	}
	s.OnlyExt = true
	s.Ext = "." + str
	return nil
}

const (
	errmsg = "A single argument is expected for the source directory, but get: %v\n"
)

func New(args []string) *Settings {
	var stngs Settings

	flagset := flag.NewFlagSet("myFind", flag.ContinueOnError)
	flagset.SetOutput(os.Stderr)

	flagset.BoolFunc("f", "Print filenames", stngs.handleF)
	flagset.BoolFunc("d", "Print directories", stngs.handleD)
	flagset.BoolFunc("sl", "Print symlinks", stngs.handleSl)
	flagset.Func("ext",
		"Print Only files with a certain extension",
		stngs.handleExt)

	if err := flagset.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "%v\n", err)
		return nil
	}

	if !stngs.PrintFilenames && !stngs.PrintDirectories &&
		!stngs.PrintSymlinks {
		stngs.PrintFilenames, stngs.PrintDirectories = true, true
		stngs.PrintSymlinks, stngs.OnlyExt = true, false
	}

	if flagset.NArg() != 1 {
		fmt.Fprintf(os.Stderr, errmsg, flagset.Args())
		flagset.PrintDefaults()
		return nil
	}

	stngs.Dirname = flagset.Arg(0)

	return &stngs
}
