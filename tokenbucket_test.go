package tokenbucket

import (
	"testing"
	"time"
)

func TestTokenBucketReserve(t *testing.T) {
	startTime := time.Now()
	// 每秒产生3个，最多可存5个，在startTime时桶内只有0个token
	bucket := New(3, time.Second, 5, 0, startTime)

	// 刚开始时获取不到token
	err := bucket.Reserve(1, startTime)
	if err == nil {
		t.Fatalf("there should be 0 token at %s", startTime.String())
	}

	// 1秒后可以获取3个
	t1 := startTime.Add(time.Second)
	err = bucket.Reserve(3, t1)
	if err != nil {
		t.Fatalf("there should be 3 tokens after 1 second, but '%s'", err)
	}
	// 此时已经没有了
	err = bucket.Reserve(1, t1)
	if err == nil {
		t.Fatal("there should be 0 tokens after reserving 3 tokens")
	}

	// 2秒后有5个token
	t1 = t1.Add(2 * time.Second)
	err = bucket.Reserve(5, t1)
	if err != nil {
		t.Fatalf("there should be 5 tokens after 2 seconds, but '%s'", err)
	}
	// 此时已经没有了
	err = bucket.Reserve(1, t1)
	if err == nil {
		t.Fatal("there should be 0 tokens after reserving all tokens")
	}
}

func BenchmarkTokenBucket_Reserve(b *testing.B) {
	startTime := time.Now()
	bucket := New(3, time.Second, 5, 0, startTime)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		startTime = startTime.Add(time.Second)
		err := bucket.Reserve(3, startTime)
		if err != nil {
			b.Fatal(err)
		}
	}
}
