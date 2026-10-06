package sketch

import "testing"

func TestSingleBucket(t *testing.T) {
	value, total := Quantile([]Bucket{{0, 10, 4}}, 50, &Counter{})
	if value != 10 || total != 4 {
		t.Fatalf("单桶应返回上界与总数，得到 %d/%d", value, total)
	}
}

func TestEmptyHistogram(t *testing.T) {
	value, total := Quantile([]Bucket{{0, 10, 0}}, 50, &Counter{})
	if value != 0 || total != 0 {
		t.Fatalf("空直方图应返回 0/0，得到 %d/%d", value, total)
	}
}

func TestAllCoversEveryPercentile(t *testing.T) {
	out := All([]Bucket{{0, 10, 10}}, &Counter{})
	for _, p := range Percentiles {
		if _, ok := out["p"+itoa(p)]; !ok {
			t.Fatalf("缺少 p%d：%+v", p, out)
		}
	}
}

func TestCountsAreSummed(t *testing.T) {
	_, total := Quantile([]Bucket{{0, 10, 3}, {10, 20, 4}}, 90, &Counter{})
	if total != 7 {
		t.Fatalf("总数应是 7，得到 %d", total)
	}
}

func TestLastBucketIsFallback(t *testing.T) {
	value, _ := Quantile([]Bucket{{0, 10, 1}, {10, 20, 1}}, 100, &Counter{})
	if value != 20 {
		t.Fatalf("最高分位应落到最后一个桶的上界，得到 %d", value)
	}
}
