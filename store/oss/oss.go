package oss

import (
	"fmt"
	"io"
	"log"
	"netdisk/setting"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var Client *oss.Client
var BuketName string
var err error

func Init(cfg *setting.OSSConfig) error {

	Client, err = oss.New(cfg.EndPoint, cfg.AccessKeyId, cfg.AccessKeySecret)

	BuketName = cfg.BucketName
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func Upload(fd io.Reader, objectkey string) error {
	bucket, err := Client.Bucket(BuketName)
	if err != nil {
		log.Panicln("buket not found", err)
		return err
	}
	err = bucket.PutObject(objectkey, fd)
	if err != nil {
		log.Panicln("put object wrong", err)
		return err
	}
	return nil
}

func UploadLocalFile(objectname string, filename string) error {
	bucket, err := Client.Bucket(BuketName)
	if err != nil {
		log.Panicln("buket not found", err)
		return err
	}
	err = bucket.PutObjectFromFile(objectname, filename)
	if err != nil {
		log.Panicln("put object wrong", err)
		return err
	}
	return nil
}

// Client : 创建oss client对象
//func Client() *oss.Client {
//	if ossCli != nil {
//		return ossCli
//	}
//	ossCli, err := oss.New(setting.Conf.OSSConfig.EndPoint,setting.)
//	if err != nil {
//		fmt.Println(err.Error())
//		return nil
//	}
//	return ossCli
//}

// Bucket : 获取bucket存储空间
//func Bucket() *oss.Bucket {
//	cli := Client()
//	if cli != nil {
//		bucket, err := cli.Bucket(cfg.OSSBucket)
//		if err != nil {
//			fmt.Println(err.Error())
//			return nil
//		}
//		return bucket
//	}
//	return nil
//}
