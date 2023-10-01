package core

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"gvb_server/config"
	"gvb_server/global"
	"io/fs"
	"io/ioutil"
	"log"
)

const ConfigFile = "settings.yaml"

// 对conf的初始化
// 读.yaml文件
// InitConf 读取yaml文件的配置
func InitConf() {

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

// 修改配置文件
func SetYaml() error {
	// 这个函数接受一个任意类型的输入参数in，并返回一个字节切片out和一个错误err
	// 这里是将整个配置文件转换成字节切片
	byteData, err := yaml.Marshal(global.Config)
	if err != nil {
		return err
	}

	// func WriteFile(filename string, data []byte, perm fs.FileMode) error
	err = ioutil.WriteFile(ConfigFile, byteData, fs.ModePerm)
	if err != nil {
		return err
	}

	global.Log.Info("配置文件修改成功")
	return nil
}
