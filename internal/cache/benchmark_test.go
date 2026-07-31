package cache

import (
	"fmt"
	"testing"
)

const benchmarkKeyCount = 100 // Cache's default capacity.

func benchmarkCache() *Cache {
	c := New()
	for i := 0; i < benchmarkKeyCount; i++ {
		c.Set(fmt.Sprintf("key-%d", i), "value", 0)
	}
	return c
}

func BenchmarkSetOverwrite(b *testing.B) {
	c := benchmarkCache()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Set(fmt.Sprintf("key-%d", i%benchmarkKeyCount), "value", 0)
	}
}

func BenchmarkGetHit(b *testing.B) {
	c := benchmarkCache()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, ok := c.Get(fmt.Sprintf("key-%d", i%benchmarkKeyCount)); !ok {
			b.Fatal("expected cache hit")
		}
	}
}

// BenchmarkMixedParallel models concurrent cache use: three reads for each write.
func BenchmarkMixedParallel(b *testing.B) {
	c := benchmarkCache()
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key-%d", i%benchmarkKeyCount)
			if i%4 == 0 {
				c.Set(key, "value", 0)
			} else if _, ok := c.Get(key); !ok {
				b.Fatal("expected cache hit")
			}
			i++
		}
	})
}
