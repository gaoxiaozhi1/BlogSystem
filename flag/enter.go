package flag

import (
	sys_flag "flag"
)

type Option struct {
	DB bool
}

// Parse 解析命令行参数
func Parse() Option {
	// 定义了一个名为 “db” 的命令行标志。这个标志的默认值是 false，并且它的描述是 “初始化数据库”。
	db := sys_flag.Bool("db", false, "初始化数据库")
	// 解析命令行参数写入注册的flag里
	sys_flag.Parse()
	return Option{
		// DB 字段被设置为 “db” 标志的值（通过 *db 解引用得到）。
		DB: *db, // false
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
	}
}
