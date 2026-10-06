# quantile-sketch

从分桶直方图估分位点。只用标准库。

```
go test ./...
go vet ./...
go run ./cmd/qs --file buckets.csv
go run ./cmd/qs --sample=boundary
go run ./cmd/qs --sample=order
go run ./cmd/qs --sample=work
```

CSV 列为 `low`、`high`、`count`，第一行是表头。

## 口径

- **桶**：左闭右开区间；分位落在哪个桶就报那个桶的上界。
- **顺序**：桶按左边界升序使用。
- **落点**：累计计数**达到**目标值（`总数 × p / 100`）的那个桶就是答案，正好相等也算。
- **空直方图**：总数为 0 时返回 0。
- **分位点**：固定算 p50、p90、p99。

## 不变量

- 分位值一定落在某个桶的上界上，且不超过最大上界。
- 结果与桶在 CSV 里的书写顺序无关。
- `scanned` 不随分位点数乘桶数放大：两千个桶算三个分位点的 `scanned` 不超过 4000。

## 输出契约

打印一行 JSON：`p50`、`p90`、`p99`；`--sample=work` 另含 `buckets` 与 `scanned`。
