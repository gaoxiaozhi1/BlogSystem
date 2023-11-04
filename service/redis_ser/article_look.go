package redis_ser

import (
	"gvb_server/global"
	"strconv"
)

// 针对点赞
const lookPrefix = "look"

// Look 浏览某一篇文章,充当缓存的作用，因为如果每次直接对es操作，就会很消耗时间，所以将这些存到redis中
func Look(id string) error {
	num, _ := global.Redis.HGet(lookPrefix, id).Int()
	num++
	err := global.Redis.HSet(lookPrefix, id, num).Err()
	return err
}

// GetLook 获取某一篇文章下的浏览量
func GetLook(id string) int {
	num, _ := global.Redis.HGet(lookPrefix, id).Int()
	return num
}

// GetLookInfo 取出浏览量数据
func GetLookInfo() map[string]int {
	var DiggInfo = map[string]int{}
	maps := global.Redis.HGetAll(lookPrefix).Val()
	for id, val := range maps {
		num, _ := strconv.Atoi(val) // 转换成int类型
		DiggInfo[id] = num
	}
	return DiggInfo
}

// 每隔一段时间同步点赞数据到es
// 这个操作写在service中es对应的操作处...

// LookClear 同步完之后要清空redis
func LookClear() {
	global.Redis.Del(lookPrefix)
}
