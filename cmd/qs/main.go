// Command qs 是分位估计的场景入口。
package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"example.com/quantile-sketch/sketch"
)

func emit(v any) {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func load(path string) []sketch.Bucket {
	fh, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fh.Close()
	rows, err := csv.NewReader(fh).ReadAll()
	if err != nil {
		panic(err)
	}
	out := []sketch.Bucket{}
	for i, row := range rows {
		if i == 0 || len(row) < 3 {
			continue
		}
		low, _ := strconv.Atoi(row[0])
		high, _ := strconv.Atoi(row[1])
		count, _ := strconv.Atoi(row[2])
		out = append(out, sketch.Bucket{Low: low, High: high, Count: count})
	}
	return out
}

func main() {
	sample := flag.String("sample", "", "boundary|order|work")
	file := flag.String("file", "", "桶 CSV 路径")
	flag.Parse()

	c := &sketch.Counter{}
	switch *sample {
	case "boundary":
		buckets := []sketch.Bucket{{Low: 0, High: 10, Count: 5}, {Low: 10, High: 20, Count: 5}}
		value, _ := sketch.Quantile(buckets, 50, c)
		emit(map[string]int{"p50": value})
	case "order":
		buckets := []sketch.Bucket{{Low: 10, High: 20, Count: 3}, {Low: 0, High: 10, Count: 3}}
		value, _ := sketch.Quantile(buckets, 70, c)
		emit(map[string]int{"p70": value})
	case "work":
		buckets := []sketch.Bucket{}
		for i := 0; i < 2000; i++ {
			buckets = append(buckets, sketch.Bucket{Low: i, High: i + 1, Count: i % 7})
		}
		out := sketch.All(buckets, c)
		emit(map[string]any{"buckets": 2000, "values": out, "scanned": c.Scanned})
	case "":
		if *file == "" {
			flag.Usage()
			os.Exit(2)
		}
		emit(sketch.All(load(*file), c))
	default:
		fmt.Fprintln(os.Stderr, "未知场景："+*sample)
		os.Exit(2)
	}
}
