// reads input from stdin and writes compressed/uncompressed data to stdout
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bruston/squish"
)

func main() {
	compress := flag.Bool("c", true, "compress input if true, else decompress")
	flag.Parse()
	if *compress {
		if err := squish.Encode(os.Stdout, os.Stdin); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		if err := squish.Decode(os.Stdin, os.Stdout); err != nil {
			fmt.Println(os.Stderr, err)
			os.Exit(1)
		}
	}
}
