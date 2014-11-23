package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type Result struct {
	text   string
	line   int
	length int
}

func (r Result) String() string {
	return fmt.Sprintf("%d %d %s", r.line, r.length, r.text)
}

func main() {
	file := flag.String("file", "", "Input file to use. If none is supplied input from stdin is assumed.")
	flag.Parse()
	var result Result
	if *file == "" {
		result = Longest(os.Stdin)
	} else {
		f, err := os.Open(*file)
		if err != nil {
			log.Fatal("unable to open input file:", err)
		}
		result = Longest(f)
		f.Close()
	}
	fmt.Println(result)
}

func Longest(r io.Reader) Result {
	scanner := bufio.NewScanner(r)
	var result Result
	var lineNum int
	for scanner.Scan() {
		lineNum++
		if len(scanner.Text()) < result.length {
			continue
		}
		result = Result{
			text:   scanner.Text(),
			line:   lineNum,
			length: len(scanner.Text()),
		}
	}
	return result
}
