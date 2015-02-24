package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	feedURL          = "https://www.archlinux.org/feeds/news/"
	rssDateLayout    = "Mon, 02 Jan 2006 15:04:05 -0700"
	dateOutputLayout = "2006-Jan-02"
	boldPrefix       = "\033[1m"
	boldSuffix       = "\033[0m"
)

type rss struct {
	Items items `xml:"channel"`
}

type items struct {
	List []item `xml:"item"`
}

type item struct {
	Title     string `xml:"title"`
	Link      string `xml:"link"`
	Published string `xml:"pubDate"`
}

func (i item) String() string {
	return fmt.Sprintf("%s: %s %s", parseDate(i.Published), i.Title, i.Link)
}

func parseDate(date string) string {
	t, err := time.Parse(rssDateLayout, date)
	if err != nil {
		return ""
	}
	return t.Format(dateOutputLayout)
}

func bold(msg string) string {
	return fmt.Sprintf("%s%s%s", boldPrefix, msg, boldSuffix)
}

func main() {
	flagEmphasis := flag.Bool("emphasis", true, "Disables bold emphasis of important news items if set to false.")
	flag.Parse()

	var news rss
	resp, err := http.Get(feedURL)
	if err != nil {
		log.Fatalf("unable to reach the news feed: %s", err)
	}
	defer resp.Body.Close()

	if err := xml.NewDecoder(resp.Body).Decode(&news); err != nil {
		log.Fatalf("unable to parse feed: %s", err)
	}

	for i := len(news.Items.List) - 1; i >= 0; i-- {
		itm := news.Items.List[i]
		if strings.Contains(strings.ToLower(itm.Title), "intervention") && *flagEmphasis {
			fmt.Printf(bold("%s\n"), itm)
		} else {
			fmt.Println(itm)
		}
	}
}
