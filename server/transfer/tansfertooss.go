package main

import (
	"netdisk/mq"
)

func main() {

	//if err := setting.Init(); err != nil {
	//
	//	fmt.Printf("load config failed, err:%v\n", err)
	//	return
	//}
	//if err := oss.Init(setting.Conf.OSSConfig); err != nil {
	//	fmt.Printf("init oss failed, err:%v\n", err)
	//	return
	//}

	rabbit := mq.NewRabbitMQSimple("oss")
	for {
		rabbit.ConsumeSimple()
	}

}
