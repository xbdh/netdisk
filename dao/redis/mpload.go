package redis

import (
	"log"
	"strconv"
	"strings"
)

//func FileMultiPartUploadInfo(redisKey string, arg map[string]string) {
//
//	for k, v := range arg {
//		client.HSet(redisKey, k, v)
//	}
//
//}

func FileMultiPartUploadInfo(redisKey string, key string, value interface{}) {

	client.HSet(redisKey, key, value)
}

func FileMultiPartUploadIsComplete(redisKey string) bool {
	mapcmd := client.HGetAll("mpload")
	if mapcmd.Err() != nil {
		log.Panicln("redis 错误", mapcmd.Err())
		return
	}
	data := mapcmd.Val()
	totalCount := 0
	chunkCount := 0
	for k, v := range data {

		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}
	if totalCount != chunkCount {

		return false
	}

	return true
}
