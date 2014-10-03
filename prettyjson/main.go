package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	indent := flag.String("indent", "  ", "How the JSON should be indented. Defaults to two spaces.")
	flag.Parse()
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("couldn't read from Stdin: %s", err.Error())
	}
	var out bytes.Buffer
	if err = json.Indent(&out, b, "", *indent); err != nil {
		log.Fatalf("bad json: %s", err.Error())
	}
	fmt.Println(out.String())
}
