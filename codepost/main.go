package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	comment    = "// "
	startBlock = "```Go"
	endBlock   = "```"
	spaceJunk  = " \t"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	inCodeBlock := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			fmt.Fprintln(os.Stdout, "")
			continue
		}
		if strings.HasPrefix(strings.TrimLeft(line, spaceJunk), comment) {
			if inCodeBlock {
				inCodeBlock = false
				fmt.Fprintln(os.Stdout, endBlock)
			}
			fmt.Fprintln(os.Stdout, strings.TrimPrefix(line, comment))
		} else {
			if !inCodeBlock {
				fmt.Fprintln(os.Stdout, startBlock)
				inCodeBlock = true
			}
			fmt.Fprintln(os.Stdout, line)
		}
	}
	if scanner.Err() != nil {
		log.Print(scanner.Err())
	}
}
