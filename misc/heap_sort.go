package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
)

func usage() {
	fmt.Println(`Usage: heap_sort list_count url
e.g.: 	heap_sort 20 http://localhost:6060/debug/pprof/heap\?debug\=1
`)
}

func rankByWordCount(wordFrequencies map[string]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func main() {
	usage()
	var err error

	list_cnt := 20
	url := "http://10.211.55.8:6060/debug/pprof/heap?debug=1"

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 2 {
		list_cnt, err = strconv.Atoi(argsWithoutProg[0])
		if err != nil {
			fmt.Println("list_cnt must be int")
			return
		}
		url = argsWithoutProg[1]
	}

	pat := regexp.MustCompile(`^#.*(github.com\\oliveagle\\.*:\d+)$`)

	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error to request: ", err)
		return
	}
	defer resp.Body.Close()

	lines := map[string]int{}

	r := bufio.NewReaderSize(resp.Body, 4*1024)
	line, isPrefix, err := r.ReadLine()
	for err == nil && !isPrefix {
		s := string(line)
		for _, ss := range pat.FindAllStringSubmatch(s, -1) {
			if len(ss) > 1 {
				// fmt.Println(ss[1])
				key := ss[1]
				_, ok := lines[key]
				if ok != true {
					lines[key] = 1
				} else {
					lines[key] += 1
				}
			}
		}
		// fmt.Println(s)
		line, isPrefix, err = r.ReadLine()
	}

	sm := rankByWordCount(lines)

	if list_cnt > len(sm) {
		list_cnt = len(sm)
	}
	for _, p := range sm[:list_cnt] {
		fmt.Printf("%d\t%s\n", p.Value, p.Key)
	}
}
