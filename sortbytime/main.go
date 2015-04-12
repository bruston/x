package main

import (
	"fmt"
	"sort"
	"time"
)

type thing struct {
	time.Time
}

type things []thing

func (t things) Len() int           { return len(t) }
func (t things) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t things) Less(i, j int) bool { return t[i].Time.Before(t[j].Time) }

func newDate(month int, day int) time.Time {
	return time.Date(2015, time.Month(month), day, 0, 0, 0, 0, time.Local)
}

func main() {
	myThings := things{
		{newDate(10, 25)},
		{newDate(10, 10)},
		{newDate(10, 16)},
		{newDate(4, 20)},
	}

	sort.Sort(myThings)
	fmt.Println("Ascending:")

	for _, v := range myThings {
		fmt.Println(v)
	}
	// Ascending:
	// 2015-04-20 00:00:00 +0100 BST
	// 2015-10-10 00:00:00 +0100 BST
	// 2015-10-16 00:00:00 +0100 BST
	// 2015-10-25 00:00:00 +0100 BST

	sort.Sort(sort.Reverse(myThings))
	fmt.Println("Descending:")

	for _, v := range myThings {
		fmt.Println(v)
	}
	// Descending:
	// 2015-10-25 00:00:00 +0100 BST
	// 2015-10-16 00:00:00 +0100 BST
	// 2015-10-10 00:00:00 +0100 BST
	// 2015-04-20 00:00:00 +0100 BST
}
