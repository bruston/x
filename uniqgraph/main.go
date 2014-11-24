package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type DataPoint struct {
	Label string `json:"label"`
	Y     int    `json:"y"`
}

type DataCollection []DataPoint

type Graph struct {
	Type       string
	DataPoints string
	Interval   int
}

func main() {
	graphType := flag.String("t", "bar", "The type of graph to generate.")
	interval := flag.Int("i", 0, "Marker interval.")
	flag.Parse()

	var points DataCollection
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		args := strings.Split(strings.TrimLeft(scanner.Text(), " \t"), " ")
		num, err := strconv.Atoi(args[0])
		if len(args) != 2 || err != nil {
			log.Print("skipping bad input line:", scanner.Text())
			continue
		}
		points = append(points, DataPoint{Label: args[1], Y: num})
	}

	b, _ := json.Marshal(points)
	g := Graph{
		Type:       *graphType,
		DataPoints: string(b),
		Interval:   *interval,
	}
	t, err := template.ParseFiles("graph.tmpl.html")
	if err != nil {
		log.Fatalf("unable to parse template file:", err)
	}
	t.Execute(os.Stdout, g)
}
