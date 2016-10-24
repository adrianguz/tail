package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	filename = flag.String("file", "", "filename to tail")
	lines    = flag.Int("lines", 10, "number of lines to tail")
)

func main() {
	flag.StringVar(filename, "f", "", "Filename to tail")
	flag.IntVar(lines, "n", 10, "Number of lines")
	flag.Parse()

	if *filename == "" || *lines == 0 {
		fmt.Println("Please supply an argument")
		os.Exit(-1)
	}
	if *lines < 0 {
		log.Fatalf("%d lines cannot be less than zero", *lines)
	}

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("Can't open file %s", err)
	}
	defer file.Close()
	var incSize int64 = 50
	bufSize := incSize
	stat, err := file.Stat()

	for {
		buffer := make([]byte, bufSize)
		start := stat.Size() - bufSize
		_, err = file.ReadAt(buffer, start)
		if err != nil {
			fmt.Printf("%s\n", err)
			break
		}
		if bytes.Count(buffer, []byte{'\n'}) > *lines {
			break
		}

		if stat.Size()-bufSize < incSize {
			bufSize += incSize
		} else {
			bufSize = stat.Size()
			break
		}

		if bufSize > stat.Size() {
			break
		}
	}
	buffer := make([]byte, bufSize)
	start := stat.Size() - bufSize
	_, err = file.ReadAt(buffer, start)
	if err != nil {
		log.Fatal(err)
	}

	out := strings.Split(string(buffer), "\n")
	l := bytes.Count(buffer, []byte{'\n'})
	for _, i := range out[l-*lines:] {
		fmt.Println(i)
	}
}
