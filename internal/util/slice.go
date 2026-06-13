// Package util 提供通用的工具函数。
package util

// ToAny 将任意类型切片转换为 []any，用于无法直接展开泛型切片的场景（如 go-redis 的变参）。
// Go 的类型系统不允许直接将 []T 展开为 ...any，这里通过循环做一次内存拷贝。
func ToAny[T any](s []T) []any {
	r := make([]any, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}
