package tokenbucket

import (
	"errors"
	"sync"
	"time"
)

// ErrOverflow 计算token时溢出
var ErrOverflow = errors.New("overflow")

// TokenBucket 令牌桶
type TokenBucket struct {
	lastModifiedTime int64  // 上次修改时间
	storedTokens     uint64 // 桶中存储的令牌数
	count            uint64 // 每inter时间内产生count个令牌
	inter            int64  // 产生count个令牌的时间
	maxTokens        uint64 // 最大的令牌数
	sync.RWMutex
}

// New 创建一个令牌桶，每隔inter时间产生count个令牌，桶内最多存储maxTokens个令牌，初始状态从startTime开始，桶内有tokensNow个令牌
func New(count uint64, inter time.Duration, maxTokens uint64, tokensNow uint64, startTime time.Time) *TokenBucket {
	return &TokenBucket{
		count:            count,
		inter:            inter.Nanoseconds(),
		maxTokens:        maxTokens,
		storedTokens:     tokensNow,
		lastModifiedTime: startTime.UnixNano(),
	}
}

// Reserve 预约count个令牌，返回是否预约成功
func (b *TokenBucket) Reserve(count uint64) bool {
	return b.ReserveWithTime(count, time.Now())
}

// ReserveWithTime 预约count个令牌，err == nil 表示成功，否则失败
func (b *TokenBucket) ReserveWithTime(count uint64, now time.Time) bool {
	if count <= 0 {
		return true
	}

	b.Lock()

	b.sync(now.UnixNano())

	storedTokens := b.storedTokens
	if storedTokens < count {
		b.Unlock()
		return false
	}
	b.storedTokens -= count
	b.Unlock()
	return true
}

// 同步状态到now时间点
func (b *TokenBucket) sync(nowNano int64) {
	diff := nowNano - b.lastModifiedTime
	if diff < 0 {
		return
	}
	tokensToPut := uint64(diff/b.inter) * b.count
	if tokensToPut < 1 {
		return
	}

	if sum, e := b.checkedAddUint64(b.storedTokens, tokensToPut); e == nil {
		if sum > b.maxTokens {
			sum = b.maxTokens
		}
		b.storedTokens = sum
	} else {
		return
	}
	b.lastModifiedTime = nowNano
	return
}

// sum = a + b，如果溢出，err不为nil。
func (b *TokenBucket) checkedAddUint64(n1, n2 uint64) (sum uint64, err error) {
	sum = n1 + n2
	if !(((n1 ^ n2) < 0) || ((n1 ^ sum) >= 0)) {
		err = ErrOverflow
	}
	return
}

// SetRate 设置令牌产生的速度为rate。
// 当前时间之前，令牌产生速度按照之前的设置。
func (b *TokenBucket) SetRate(count uint64, inter time.Duration) {
	b.Lock()
	// 先将状态更新到当前时间，然后在设置速度
	b.sync(time.Now().UnixNano())
	b.count = count
	b.inter = inter.Nanoseconds()
	b.Unlock()
}

// SetMaxTokens 设置最大令牌数
func (b *TokenBucket) SetMaxTokens(max uint64) {
	b.Lock()
	// 先将状态更新到当前时间，然后在设置速度
	b.sync(time.Now().UnixNano())
	b.maxTokens = max
	b.Unlock()
	return
}

// GetStoredTokensNow 获取当前存储的令牌数，该方法会将令牌桶同步到当前时间点。
func (b *TokenBucket) GetStoredTokensNow() (tokens uint64) {
	b.Lock()
	b.sync(time.Now().UnixNano())
	tokens = b.storedTokens
	b.Unlock()
	return
}
