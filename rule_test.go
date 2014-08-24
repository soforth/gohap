package filter

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"
)

var samples = []string{"condition", "false", "function", "query", "true", "abnormal"}

func TestFilter(t *testing.T) {
	fmt.Println("[!!NOTICE!!] IGNORE the error report if file name is 'abnormal'")
	for _, file := range samples {
		buf, err := ioutil.ReadFile("./sample/" + file)
		if err != nil {
			t.Fatal(err)
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
			h, err := NewParser(strings.NewReader(v[1]))
			if err != nil {
				fmt.Println("file:", file, "line:", i+1, v, err)
				continue
			}
			var symlist *SymList
			if v[2][0] == '{' {
				symlist, _ = JsonToSymlist(v[2])
			} else {
				symlist, _ = QueryToSymlist(v[2])
			}
			actual, err := h.Parse(symlist)
			if err != nil {
				fmt.Println("file:", file, "line:", i+1, "Parse error:", err)
			} else {
				expect, _ := strconv.Atoi(v[0])
				if actual != expect {
					t.Errorf("file: %s line: %d expect %d, actual %d\n", file, i+1, expect, actual)
				}
			}
		}
	}
}

func BenchmarkFilter(t *testing.B) {
	rand.Seed(int64(time.Now().Second()))
	buf, err := ioutil.ReadFile("./sample/bench")
	if err != nil {
		t.Fatal(err)
	}
	// pre-store Parser handle and symbol list
	var values [][]string
	var handle []*Parser
	lines := strings.Split(string(buf), "\n")
	for _, line := range lines {
		if len(line) == 0 || line[0] == '\n' || line[0] == '/' {
			continue
		}
		v := strings.Split(line, "%")
		if len(v) != 3 {
			continue
		}
		values = append(values, v)
		h, err := NewParser(strings.NewReader(v[1]))
		if err != nil {
			t.Fatal(err)
		}
		handle = append(handle, h)

	}

	for i := 0; i < t.N; i++ {
		index := rand.Intn(len(values))
		v := values[index]
		var symlist *SymList
		if v[2][0] == '{' {
			symlist, _ = JsonToSymlist(v[2])
		} else {
			symlist, _ = QueryToSymlist(v[2])
		}

		h := handle[index]
		actual, err := h.Parse(symlist)
		if err != nil {
			t.Error("line:", i+1, "Parse error:", err)
		} else {
			expect, _ := strconv.Atoi(v[0])
			if actual != expect {
				t.Errorf("line: %d expect %d, actual %d\n", i+1, expect, actual)
			}
		}
	}
}
