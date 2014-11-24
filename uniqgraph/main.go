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

var html = `
<!DOCTYPE HTML>
<html>
<head>
  <script type="text/javascript">
      window.onload = function () {
          var chart = new CanvasJS.Chart("chartContainer", {
              theme: "theme2",//theme1
              title:{
                  text: "'uniq -c' Graph"
             },
              axisY:{
                  interval: {{ .Interval }}
              },
              data: [
              {
                  type: "{{ .Type }}",
                  dataPoints: {{ .DataPoints }}
              }]
          });

          chart.render();
      }
  </script>
</head>
<body>
  <div id="chartContainer" style="height: 800px; width: 100%;">
  </div>
  <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/canvasjs/1.4.1/canvas.min.js"></script>
</body>
</html>
`

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
	customTemplate := flag.String("template", "", "Custom template to use to generate the graph HTML. Not required.")
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
	var t *template.Template
	var err error
	if *customTemplate != "" {
		t, err = template.ParseFiles(*customTemplate)
	} else {
		t, err = template.New("html").Parse(html)
	}
	if err != nil {
		log.Fatalf("unable to parse template file:", err)
	}
	t.Execute(os.Stdout, g)
}
