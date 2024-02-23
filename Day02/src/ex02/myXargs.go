package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{"echo"}
	}

	stdinScanner := bufio.NewScanner(os.Stdin)
	for stdinScanner.Scan() {
		args = append(args, strings.TrimSpace(stdinScanner.Text()))
	}

	path, err := exec.LookPath(args[0])
	if errors.Is(err, exec.ErrDot) {
		err = nil
	}
	if err != nil {
		log.Fatal(err)
	}

	err = syscall.Exec(path, args, os.Environ())
	if err != nil {
		log.Fatal(err)
	}
}
