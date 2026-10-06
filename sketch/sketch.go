// Package sketch 从分桶直方图估分位。
package sketch

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
	total := 0
	for _, b := range buckets {
		c.Scanned++
		total += b.Count
	}
	if total == 0 {
		return 0, 0
	}
	want := total * p / 100
	seen := 0
	for _, b := range buckets {
		c.Scanned++
		seen += b.Count
		if seen > want {
			return b.High, total
		}
	}
	return buckets[len(buckets)-1].High, total
}

// All 算全部分位点。
func All(buckets []Bucket, c *Counter) map[string]int {
	out := map[string]int{}
	for _, p := range Percentiles {
		value, _ := Quantile(buckets, p, c)
		out["p"+itoa(p)] = value
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
