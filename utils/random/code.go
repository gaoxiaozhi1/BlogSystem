package random

import (
	"fmt"
	"math/rand"
	"time"
)

var stringCode = ""

func Code(length int) string {
	rand.Seed(time.Now().UnixNano()) // 种子纳秒
	return fmt.Sprintf("%4v", rand.Intn(10000))
}
