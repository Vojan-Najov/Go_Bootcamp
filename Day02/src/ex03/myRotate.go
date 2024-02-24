package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

	var wg sync.WaitGroup
	filenames := flag.Args()
	for _, filename := range filenames {
		wg.Add(1)
		go func(filename, dirname string) {
			defer wg.Done()
			err := createArchive(filename, dirname)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}(dir.path, filename)
	}
	wg.Wait()
}

func destFilename(dirname, filename string) (string, error) {
	filename, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}

	fileInfo, err := os.Stat(filename)
	if err != nil {
		return "", err
	}

	if len(dirname) == 0 {
		dirname = filepath.Dir(filename)
	} else {
		dirname, err = filepath.Abs(dirname)
	}
	if err != nil {
		return "", err
	}

	basename := filepath.Base(filename)
	basename, _ = strings.CutSuffix(basename, filepath.Ext(basename))
	basename += fmt.Sprintf("_%d.tar.gz", fileInfo.ModTime().Unix())

	destpath := filepath.Join(dirname, basename)

	return destpath, nil
}

func createArchive(dirname, filename string) error {
	destname, err := destFilename(dirname, filename)
	if err != nil {
		return err
	}

	out, err := os.Create(destname)
	if err != nil {
		return err
	}
	defer out.Close()

	gw := gzip.NewWriter(out)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}

	return nil
}
