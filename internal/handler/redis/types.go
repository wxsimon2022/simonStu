package redis

// RedisMultiGetRequest 批量查询 Redis 的请求参数。
// Keys 是要查询的 key 列表，最少 1 个，没有上限（但实际受 Redis 性能和网络带宽影响）。
type RedisMultiGetRequest struct {
	Keys []string `json:"keys" binding:"required,min=1"`
}

// RedisMultiSetRequest 批量写入 Redis 的请求参数。
// Values 是 key-value 对，全部并发写入（非原子，部分成功部分失败是可能的）。
type RedisMultiSetRequest struct {
	Values map[string]string `json:"values" binding:"required"`
}

// redisGetResult 单个 key 的查询结果。
// Key/Value 自解释；Err 为空表示成功；CostMs 是单次查询耗时，可用于分析慢 key。
type redisGetResult struct {
	Key    string `json:"key"`
	Value  string `json:"value,omitempty"`
	Err    string `json:"err,omitempty"`
	CostMs int64  `json:"cost_ms"`
}
