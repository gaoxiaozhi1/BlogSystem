package main

import (
	"fmt"
)

func main() {
	var (
		userName string
		email    string
	)
	fmt.Printf("请输入用户名")
	fmt.Scan(&userName)
	fmt.Printf("请输入邮箱")
	fmt.Scanln(&email) // 不用必须输入

	fmt.Println(userName, email)

}
