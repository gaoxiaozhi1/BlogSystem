package core

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"gvb_server/config"
	"gvb_server/global"
	"io/ioutil"
	"log"
)

// 对conf的初始化
// 读.yaml文件
// InitConf 读取yaml文件的配置
func InitConf() {
	const ConfigFile = "settings.yaml"
	c := &config.Config{} // 读取的配置文件要存储在一个实体中
	yamlConf, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		panic(fmt.Errorf("get yamlconf error: %S", err))
	}
	err = yaml.Unmarshal(yamlConf, c)
	if err != nil {
		log.Fatalf("config Init Unmarshal: %V", err)
	}
	log.Println("config yamlFile Load Init success.")
	global.Config = c // 存到一个实体中
}
