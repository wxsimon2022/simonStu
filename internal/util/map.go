// Package util 工具函数。
package util

import "encoding/json"

// SliceToMap 将 slice 转为 map，slice 内每个元素通过 keyFn 提取 key。
//
//	users := []model.Admin{...}
//	m := util.SliceToMap(users, func(v model.Admin) int { return v.ID })
//	// map[int]model.Admin{ 1: {ID:1, Username:"admin", ...}, ... }
//
// 按 username 索引：
//
//	util.SliceToMap(users, func(v model.Admin) string { return v.Username })
func SliceToMap[K comparable, V any](items []V, keyFn func(V) K) map[K]V {
	m := make(map[K]V, len(items))
	for _, v := range items {
		m[keyFn(v)] = v
	}
	return m
}

// SliceToMapGeneric 将任意 slice 转为 map[K]map[string]any，
// 序列化时按 struct 的 json tag 提取字段（json:"-" 字段会自动排除）。
//
//	users := []model.Admin{...}
//	m := util.SliceToMapGeneric(users, func(v map[string]any) int { return int(v["id"].(float64)) })
//	// map[int]map[string]any{ 1: {"id":1, "username":"admin", ...}, ... }
func SliceToMapGeneric[K comparable](items any, keyFn func(map[string]any) K) map[K]map[string]any {
	raw, err := json.Marshal(items)
	if err != nil {
		return nil
	}
	var list []map[string]any
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil
	}
	m := make(map[K]map[string]any, len(list))
	for _, v := range list {
		m[keyFn(v)] = v
	}
	return m
}
