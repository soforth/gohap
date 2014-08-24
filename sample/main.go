package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"util/filter/filter"
)

func main() {
	buf, err := ioutil.ReadFile("./test_file")
	if err != nil {
		fmt.Println(err)
		return
	}
	lines := strings.Split(string(buf), "\n")
	for i, line := range lines {
		if len(line) == 0 || line[0] == '\n' || line[0] == '/' {
			continue
		}
		v := strings.Split(line, "%")
		if len(v) != 3 {
			continue
		}
		h, err := filter.NewParser(strings.NewReader(v[1]))
		if err != nil {
			fmt.Println("line:", i+1, v, err)
			continue
		}
		var symlist *filter.SymList
		if v[2][0] == '{' {
			symlist, _ = filter.JsonToSymlist(v[2])
		} else {
			symlist, _ = filter.QueryToSymlist(v[2])
		}
		actual, err := h.Parse(symlist)
		if err != nil {
			fmt.Println("line:", i+1, "Parse error:", err)
		} else {
			expect, _ := strconv.Atoi(v[0])
			if actual != expect {
				fmt.Printf("line: %d expect %d, actual %d\n", i+1, expect, actual)
			}
		}
	}
}
