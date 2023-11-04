package main

import "gvb_server/service/es_ser"

func main() {
	var data = "### Geospatial\n- 朋友的定位，附近的人，打车距离计算\n- 底层是 Zset（**即可以使用 Zset 的命令操作 Geospatial**）\n\n`geoadd key 经度 纬度 名称`  添加地理位置\n\n- 两极无法添加\n- 经度：-180 ~ 180（度）\n- 纬度：-85.05112878 ~ 85.05112878（度）\n\n`geopos key 名称` 获取指定位置的地理位置\n`geodist key` 返回两个给定位置之间的距离（直线距离）\n--> (m, 米；km, 千米；mi, 英里；ft, 英尺)\n\n`georadius key 经度 纬度 半径 单位` 以给定值为半径，以经度和维度为中心，查找附近的人（获得所有附近的人的地址（开启定位））通过半径查询\n\n`georadiusbymember key 成员名 半径 单位` 以给定值为半径，以成员(城市名)为中心，查找\n\n`geohash key 成员1 成员2` 返回一个或多个位置元素的 geohash 表示（- 如果两个字符串越相似，表示两个地方越近）\n\n### Hyperloglog\n\n- 基数统计的算法\n- 占用内存是**固定的** \n\n**基数**：集合中元素的个数（先去重），如{1,2,2,3} 其基数为3（集合去重后为1,2,3 有3个元素）\n**网页的 UV**（一个人访问访问一个网站多次，但是还是算作一个人）\n传统实现UV：**Set**保存用户的Id，然后统计set中的元素的数量作为标准判断（这种需要保存大量用户的ID）\n![[Pasted image 20230827205857.png]]\n\n### Bitmaps\n- 位存储，位图（操作二进制）\n- 统计用户信息，活跃，不活跃！登录、未登录！打卡，365 打卡！(两个状态都可以使用)\n\n1: 打卡，0：未打卡\n`setbit sign 0 1`  周一\n`setbit sign 1 0`  周二\n`setbit sign 2 1`  周三\n`setbit sign 3 1`  周四\n...\n`getbit sign 2` 查看周三打卡情况\n`bitcount sign` 统计 key 为sign 所有打卡天数\n"
	es_ser.GetSearchIndexByContent("articles/crT7dIsBUhL2Yg9flgfi", "Redis内容", data)
}
