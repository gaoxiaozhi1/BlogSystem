package redis_ser

import (
	"gvb_server/global"
	"strconv"
)

// 针对点赞
const diggPrefix = "digg"

// Digg 点赞某一篇文章,充当缓存的作用，因为如果每次直接对es操作，就会很消耗时间，所以将这些存到redis中
func Digg(id string) error {
	num, _ := global.Redis.HGet(diggPrefix, id).Int()
	num++
	err := global.Redis.HSet(diggPrefix, id, num).Err()
	return err
}

// GetDigg 获取某一篇文章下的点赞数
func GetDigg(id string) int {
	num, _ := global.Redis.HGet(diggPrefix, id).Int()
	return num
}

// GetDiggInfo 取出点赞数据
func GetDiggInfo() map[string]int {
	var DiggInfo = map[string]int{}
	maps := global.Redis.HGetAll(diggPrefix).Val()
	for id, val := range maps {
		num, _ := strconv.Atoi(val) // 转换成int类型
		DiggInfo[id] = num
	}
	return DiggInfo
}

// 每隔一段时间同步点赞数据到es
// 这个操作写在service中es对应的操作处...

// DiggClear 同步完之后要清空redis
func DiggClear() {
	global.Redis.Del(diggPrefix)
}
