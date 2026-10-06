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

// histogram 是按左边界升序排好、带累计计数的直方图。
type histogram struct {
	buckets []Bucket
	cum     []int
	total   int
}

// prepare 复制并按左边界升序排列桶，单次扫描建立累计计数。
func prepare(buckets []Bucket, c *Counter) histogram {
	sorted := make([]Bucket, len(buckets))
	copy(sorted, buckets)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Low < sorted[j].Low })
	cum := make([]int, len(sorted))
	total := 0
	for i, b := range sorted {
		c.Scanned++
		total += b.Count
		cum[i] = total
	}
	return histogram{buckets: sorted, cum: cum, total: total}
}

// quantile 返回累计计数达到目标值的第一个桶的上界。
func (h histogram) quantile(p int) int {
	if h.total == 0 {
		return 0
	}
	want := h.total * p / 100
	idx := sort.Search(len(h.cum), func(i int) bool { return h.cum[i] >= want })
	if idx >= len(h.buckets) {
		return h.buckets[len(h.buckets)-1].High
	}
	return h.buckets[idx].High
}

// Quantile 估一个分位点，返回落在的桶上界。
func Quantile(buckets []Bucket, p int, c *Counter) (int, int) {
	h := prepare(buckets, c)
	return h.quantile(p), h.total
}

// All 算全部分位点。
func All(buckets []Bucket, c *Counter) map[string]int {
	h := prepare(buckets, c)
	out := map[string]int{}
	for _, p := range Percentiles {
		out["p"+itoa(p)] = h.quantile(p)
	}
	return out
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
