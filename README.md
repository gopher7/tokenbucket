# tokenbucket
令牌桶

## 基本使用

```go
package main

import (
	"fmt"
	"time"
	"tokenbucket"
)

func main() {
	// 每秒产生3个，最多可存5个，当前时间开始产生令牌，此时桶内只有0个token
	bucket := tokenbucket.New(3, time.Second, 5, 0, time.Now())
	// 预约1个令牌
	if !bucket.Reserve(1) {
		fmt.Printf("无法获取令牌")
	}
	// 预约3个令牌
	if !bucket.Reserve(3) {
		fmt.Printf("无法获取令牌")
	}
}
````

