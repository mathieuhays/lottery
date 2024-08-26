package main

import (
	"errors"
	"io"
	"log"
	"os"
	"strconv"
)

const usage = "[usage] lottery <chances>"

func run(args []string, out io.Writer) error {
	if len(args) < 1 {
		return errors.New("missing chances")
	}

	chances, err := strconv.Atoi(args[0])

	return nil
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		log.Fatalf("error: %s", err)
	}

	log.Println("Done")
}
