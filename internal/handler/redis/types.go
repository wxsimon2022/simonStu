package redis

// RedisMultiGetRequest 批量查询 Redis 请求参数。
type RedisMultiGetRequest struct {
	Keys []string `json:"keys" binding:"required,min=1"`
}

// RedisMultiSetRequest 批量写入 Redis 请求参数。
type RedisMultiSetRequest struct {
	Values map[string]string `json:"values" binding:"required"`
}

// redisGetResult 单个 key 的查询结果。
type redisGetResult struct {
	Key    string `json:"key"`
	Value  string `json:"value,omitempty"`
	Err    string `json:"err,omitempty"`
	CostMs int64  `json:"cost_ms"`
}
