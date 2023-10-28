package main

import (
	"fmt"
	"github.com/liu-cn/json-filter/filter"
	"time"
)

type User struct {
	UID        uint      `json:"uid,select(article)"`
	Avatar     string    `json:"avatar,select(article)"`
	Nickname   string    `json:"nickname,select(article|profile)"`
	Sex        int       `json:"sex,omit(article)"`
	VipEndTime time.Time `json:"vip_end_time,select(profile)"`
	Price      string    `json:"price,select(profile)"`
}

func NewUser() User {
	return User{
		UID:        0,
		Avatar:     "https://gimg2.baidu.com",
		Nickname:   "昵称",
		Sex:        1,
		VipEndTime: time.Now().Add(time.Hour * 24 * 365),
		Price:      "19999.9",
	}
}

func main() {
	// 显示有article标签的
	// 显示有article标签的
	fmt.Println(filter.Select("article", NewUser()))
	fmt.Println(filter.Omit("article", NewUser()))

	/*
		我们定义了一个User结构体，并使用select(scene)标签来指定在哪些场景下我们想要输出这个字段。
		然后，我们使用filter.SelectMarshal函数来选择我们想要输出的字段，并使用MustJSON方法来获取JSON字符串。
	*/
}
