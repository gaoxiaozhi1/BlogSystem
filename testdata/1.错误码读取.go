package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

const file = "models/res/err_code.json"

type ErrMap map[int]string

func main() {
	// 读取文件
	byteData, err := os.ReadFile(file)
	if err != nil {
		logrus.Error(err)
		return
	}

	var errMap = ErrMap{}
	// 将json数据转化成type ErrMap map[string]string类型
	err = json.Unmarshal(byteData, &errMap)
	if err != nil {
		logrus.Error(err)
		return
	}

	fmt.Println(errMap)
}
