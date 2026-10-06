// Package sketch 从分桶直方图估分位。
package sketch

import "sort"

// Percentiles 是要算的分位点。
var Percentiles = []int{50, 90, 99}

// Bucket 是一个左闭右开的区间桶。
type Bucket struct {
	Low   int
	High  int
	Count int
}

// Counter 记录本轮扫描过的桶数。
type Counter struct {
	Scanned int
}

// Quantile 估一个分位点，返回落在的桶上界。
func Quantile(buckets []Bucket, p int, c *Counter) (int, int) {
	bs := sorted(buckets)
	total := totalOf(bs, c)
	if total == 0 {
		return 0, 0
	}
	want := total * p / 100
	seen := 0
	for _, b := range bs {
		c.Scanned++
		seen += b.Count
		if seen >= want {
			return b.High, total
		}
	}
	return bs[len(bs)-1].High, total
}

// All 算全部分位点。
func All(buckets []Bucket, c *Counter) map[string]int {
	out := map[string]int{}
	bs := sorted(buckets)
	total := totalOf(bs, c)
	if total == 0 {
		for _, p := range Percentiles {
			out["p"+itoa(p)] = 0
		}
		return out
	}
	next := 0
	seen := 0
	for _, b := range bs {
		c.Scanned++
		seen += b.Count
		for next < len(Percentiles) && seen >= total*Percentiles[next]/100 {
			out["p"+itoa(Percentiles[next])] = b.High
			next++
		}
	}
	for ; next < len(Percentiles); next++ {
		out["p"+itoa(Percentiles[next])] = bs[len(bs)-1].High
	}
	return out
}

// sorted 返回按左边界升序排列的桶副本。
func sorted(buckets []Bucket) []Bucket {
	out := make([]Bucket, len(buckets))
	copy(out, buckets)
	sort.Slice(out, func(i, j int) bool { return out[i].Low < out[j].Low })
	return out
}

// totalOf 统计总数，每看一个桶计一次扫描。
func totalOf(buckets []Bucket, c *Counter) int {
	total := 0
	for _, b := range buckets {
		c.Scanned++
		total += b.Count
	}
	return total
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	digits := []byte{}
	for v > 0 {
		digits = append([]byte{byte('0' + v%10)}, digits...)
		v /= 10
	}
	return string(digits)
}
