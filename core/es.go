package core

import (
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"gvb_server/global"
)

func EsConnect() *elastic.Client {
	var err error
	sniffOpt := elastic.SetSniff(false)

	client, err := elastic.NewClient( // 可以设置超时时间
		elastic.SetURL(global.Config.ES.URL()),
		sniffOpt,
		elastic.SetBasicAuth(global.Config.ES.User, global.Config.ES.Password),
	)

	if err != nil {
		logrus.Fatalf("es连接失败%s", err.Error())
	}
	return client
}
