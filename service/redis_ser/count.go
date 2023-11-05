package redis_ser

import (
	"gvb_server/global"
	"strconv"
)

type CountDB struct {
	Index string // 索引前缀

}

// Set 设置某一个的数据，重复执行，重复累加
func (c CountDB) Set(id string) error {
	num, _ := global.Redis.HGet(c.Index, id).Int()
	num++
	err := global.Redis.HSet(c.Index, id, num).Err()
	return err
}

// SetCount 在原有基础上增加多少
func (c CountDB) SetCount(id string, cnt int) error {
	num, _ := global.Redis.HGet(c.Index, id).Int()
	num += cnt
	err := global.Redis.HSet(c.Index, id, num).Err()
	return err
}

// Get 获取某个的数据
func (c CountDB) Get(id string) int {
	num, _ := global.Redis.HGet(c.Index, id).Int()
	return num
}

// GetInfo 取出数据
func (c CountDB) GetInfo() map[string]int {
	var Info = map[string]int{}
	maps := global.Redis.HGetAll(c.Index).Val()
	for id, val := range maps {
		num, _ := strconv.Atoi(val) // 转换成int类型
		Info[id] = num
	}
	return Info
}

// 每隔一段时间同步点赞数据到es
// 这个操作写在service中es对应的操作处...

// Clear 同步完之后要清空redis
func (c CountDB) Clear() {
	global.Redis.Del(c.Index)
}
