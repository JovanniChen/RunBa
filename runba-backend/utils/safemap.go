// SafeMap - 基于sync.Map的泛型线程安全Map包装器
//
// 本包提供了类型安全的并发Map实现，结合了sync.Map的高性能
// 和Go泛型的类型安全特性。
//
// 特性：
// - 类型安全：编译时类型检查，避免类型断言错误
// - 高性能：底层使用sync.Map，适合读多写少场景
// - 简洁API：提供类型化的方法，无需手动类型断言
// - 零依赖：仅依赖标准库
//
// Example:
//
//	// 创建一个存储用户会话的SafeMap
//	sessions := NewSafeMap[uint, *UserSession]()
//
//	// 存储数据 - 类型安全
//	sessions.Store(123, &UserSession{ID: 123, Name: "user1"})
//
//	// 读取数据 - 无需类型断言
//	if session, ok := sessions.Load(123); ok {
//		fmt.Printf("用户: %s", session.Name) // 直接使用，无需断言
//	}
//
//	// 遍历所有数据 - 类型安全的回调
//	sessions.Range(func(id uint, session *UserSession) bool {
//		fmt.Printf("ID: %d, Name: %s", id, session.Name)
//		return true
//	})
package utils

import "sync"

// SafeMap 泛型线程安全Map
// K必须是可比较类型，V可以是任意类型
type SafeMap[K comparable, V any] struct {
	m sync.Map
}

// NewSafeMap 创建新的SafeMap实例
func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{}
}

// Store 存储键值对
func (sm *SafeMap[K, V]) Store(key K, value V) {
	sm.m.Store(key, value)
}

// Load 根据键获取值
// 返回值和是否存在的布尔值
func (sm *SafeMap[K, V]) Load(key K) (V, bool) {
	value, ok := sm.m.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return value.(V), true
}

// LoadOrStore 获取键对应的值，如果不存在则存储并返回给定值
// 返回实际值和是否是新加载的布尔值
func (sm *SafeMap[K, V]) LoadOrStore(key K, value V) (V, bool) {
	actual, loaded := sm.m.LoadOrStore(key, value)
	return actual.(V), loaded
}

// LoadAndDelete 获取键对应的值并删除该键值对
// 返回值和是否存在的布尔值
func (sm *SafeMap[K, V]) LoadAndDelete(key K) (V, bool) {
	value, ok := sm.m.LoadAndDelete(key)
	if !ok {
		var zero V
		return zero, false
	}
	return value.(V), true
}

// Delete 删除键值对
func (sm *SafeMap[K, V]) Delete(key K) {
	sm.m.Delete(key)
}

// Range 遍历所有键值对
// 回调函数返回false时停止遍历
func (sm *SafeMap[K, V]) Range(f func(key K, value V) bool) {
	sm.m.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}

// Size 获取Map中元素的大概数量
// 注意：由于并发特性，这个值可能不是精确的
func (sm *SafeMap[K, V]) Size() int {
	count := 0
	sm.m.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}

// Keys 获取所有键的切片快照
// 注意：返回的是快照，不会随Map变化而更新
func (sm *SafeMap[K, V]) Keys() []K {
	var keys []K
	sm.m.Range(func(key, value any) bool {
		keys = append(keys, key.(K))
		return true
	})
	return keys
}

// Values 获取所有值的切片快照
// 注意：返回的是快照，不会随Map变化而更新
func (sm *SafeMap[K, V]) Values() []V {
	var values []V
	sm.m.Range(func(key, value any) bool {
		values = append(values, value.(V))
		return true
	})
	return values
}

// Clear 清空所有键值对
func (sm *SafeMap[K, V]) Clear() {
	sm.m.Range(func(key, value any) bool {
		sm.m.Delete(key)
		return true
	})
}

// Has 检查键是否存在
func (sm *SafeMap[K, V]) Has(key K) bool {
	_, ok := sm.m.Load(key)
	return ok
}
