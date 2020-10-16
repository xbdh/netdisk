package main

import (
	"fmt"
	"netdisk/dao/mysql"
	"netdisk/logger"
	"netdisk/mq"
	"netdisk/pkg/snowflake"
	"netdisk/router"
	"netdisk/setting"
	"netdisk/store/oss"
)

func main() {
	if err := setting.Init(); err != nil {

		fmt.Printf("load config failed, err:%v\n", err)
		return
	}

	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	if err := mysql.Init(setting.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	if err := oss.Init(setting.Conf.OSSConfig); err != nil {
		fmt.Printf("init oss failed, err:%v\n", err)
		return
	}
	if err := mq.Init(setting.Conf.MqConfig); err != nil {
		fmt.Printf("init rabbitmq failed, err:%v\n", err)
		return
	}
	////defer mysql.Close() // 程序退出关闭数据库连接
	//if err := redis.Init(setting.Conf.RedisConfig); err != nil {
	//	fmt.Printf("init redis failed, err:%v\n", err)
	//	return
	//}
	//defer redis.Close()

	if err := snowflake.Init(1); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	//// 初始化gin框架内置的校验器使用的翻译器
	//if err := controller.InitTrans("zh"); err != nil {
	//	fmt.Printf("init validator trans failed, err:%v\n", err)
	//	return
	//}
	// 注册路由
	r := router.SetupRouter(setting.Conf.Mode)
	err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
}
