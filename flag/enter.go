package flag

import (
	sys_flag "flag"
)

type Option struct {
	DB   bool
	User string
	// -u admin (创建一个叫admin的管理员)
	// -u user (创建一个叫user的用户)
}

// Parse 解析命令行参数
func Parse() Option {
	// 定义了一个名为 “db” 的命令行标志。这个标志的默认值是 false，并且它的描述是 “初始化数据库”。
	db := sys_flag.Bool("db", false, "初始化数据库")
	user := sys_flag.String("u", "", "创建用户")
	// 解析命令行参数写入注册的flag里
	sys_flag.Parse()
	return Option{
		// DB 字段被设置为 “db” 标志的值（通过 *db 解引用得到）。
		DB:   *db, // false
		User: *user,
	}
	// 如果你在命令行中使用 -db=true 运行你的程序，那么这个 Parse 函数将会返回 { DB: true }。
	// 如果你没有指定 -db 标志，那么它将返回 { DB: false }，因为 false 是 “db” 标志的默认值。
}

// IsWebStop 是否停止web项目
func IsWebStop(option Option) bool {
	if option.DB {
		return true
	}
	return false // 默认返回这个
}

// SwitchOption 根据命令执行不同的函数
func SwitchOption(option Option) {
	if option.DB {
		Makemigrations()
		return
	}

	// 创建用户调用(管理员或普通用户)
	// 因为如果option.User为""，那么就不会走这里，必须是个字符串
	if option.User == "admin" || option.User == "user" {
		CreateUser(option.User)
		return
	}

	// 不符合上面的任意一种情况，那就显示，类似于提示帮助
	//sys_flag.Usage()
}
