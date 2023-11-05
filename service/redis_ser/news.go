package redis_ser

import (
	"encoding/json"
	"fmt"
	"gvb_server/global"
	"time"
)

const newIndex = "new_index"

type NewData struct {
	Index    string `json:"index"`
	Title    string `json:"title"`
	HotValue string `json:"hotValue"`
	Link     string `json:"link"`
}

// SetNews 设置某一个的数据，重复执行，重复累加
func SetNews(key string, newdata []NewData) error {
	byteData, _ := json.Marshal(newdata) // 转化成字符串
	err := global.Redis.Set(fmt.Sprintf("%s_%s", newIndex, key), byteData, 1*time.Hour).Err()
	return err
}

func GetNews(key string) (newData []NewData, err error) {
	res := global.Redis.Get(fmt.Sprintf("%s_%s", newIndex, key)).Val()
	err = json.Unmarshal([]byte(res), &newData)
	return
}
