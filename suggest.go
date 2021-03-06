package main

import (
	"fmt"
	"sort"

	"github.com/cloud66/cli"
)

type Suggestion struct {
	s string
	d int
}

type Suggestions []Suggestion

func (a Suggestions) Len() int           { return len(a) }
func (a Suggestions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Suggestions) Less(i, j int) bool { return a[i].d < a[j].d }

// suggest returns command names that are similar to s.
func suggest(context *cli.Context, s string) {
	a := []string{}
	var g Suggestions
	for _, c := range context.App.Commands {
		if d := editDistance(s, c.Name); d < 4 {
			g = append(g, Suggestion{c.Name, d})
		}
	}
	sort.Sort(g)
	for i, s := range g {
		a = append(a, s.s)
		if i >= 4 {
			break
		}
	}

	fmt.Printf("Can't find command %s. Did you mean any of these?\n", s)
	for _, c := range a {
		fmt.Println(c)
	}
}

func editDistance(a, b string) int {
	var d [][]int
	d = append(d, make([]int, len(b)+1))
	for i := range b {
		d[0][i+1] = i + 1
	}
	for i := range a {
		v := make([]int, len(b)+1)
		d = append(d, v)
		v[0] = i + 1
	}
	for j, cb := range []byte(b) {
		for i, ca := range []byte(a) {
			if ca == cb {
				d[i+1][j+1] = d[i][j]
			} else {
				cost := d[i][j+1]             // delete
				if v := d[i+1][j]; v < cost { // insert
					cost = v
				}
				if v := d[i][j]; v < cost { // substitute
					cost = v
				}
				if i > 0 && j > 0 {
					if ca == b[j-1] && cb == a[i-1] {
						if v := d[i-1][j-1]; v < cost { // transpose
							cost = v
						}
					}
				}
				d[i+1][j+1] = cost + 1
			}
		}
	}
	return d[len(a)][len(b)]
}
