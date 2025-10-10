package pie

import (
	"math/rand"
	"time"
)

// randomInt 生成指定范围内的随机整数
func randomInt(min, max int64) int64 {
	if min >= max {
		return min
	}
	return min + rand.Int63n(max-min+1)
}

// init 初始化随机数生成器
func init() {
	rand.Seed(time.Now().UnixNano())
}
