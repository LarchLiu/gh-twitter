package model

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/qiniu/x/rpc"
)

// QiniuRet qiniu return data type.
type QiniuRet struct {
	Key  string
	Hash string
	Type string
	URL  string
}

var domain string = os.Getenv("QINIU_DOMAIN")
var bucket string = os.Getenv("QINIU_BUCKET")
var accessKey string = os.Getenv("QINIU_ACCESS_KEY")
var secretKey string = os.Getenv("QINIU_SECRET_KEY")
var jsonPrefix string = strings.TrimPrefix(QiniuGetResourcePrefix()+"/json/", "/")

// QiniuGetResourcePrefix .
func QiniuGetResourcePrefix() string {
	// 删除 QINIU_RESOURCE_PREFIX 开头和结尾的 /
	prefix := strings.TrimPrefix(os.Getenv("QINIU_RESOURCE_PREFIX"), "/")
	prefix = strings.TrimSuffix(prefix, "/")
	return prefix
}

// QiniuUpload qiniu api to upload file.
func QiniuUpload(filePath string, key string) (ret QiniuRet, err error) {
	mac := qbox.NewMac(accessKey, secretKey)
	putPolicy := storage.PutPolicy{
		Scope:      fmt.Sprintf("%s:%s", bucket, key),
		InsertOnly: 0,
		ReturnBody: `{"key": $(key), "hash": $(etag), "type": $(mimeType)}`,
	}
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	formUploader := storage.NewFormUploader(&cfg)
	ret = QiniuRet{}
	err = formUploader.PutFile(context.Background(), &ret, upToken, key, filePath, nil)
	if err != nil {
		fmt.Println(err)
		return ret, err
	}
	ret.URL = storage.MakePublicURL(domain, ret.Key)
	return ret, nil
}

// QiniuCheckFile check if file Exist on qiniu.
func QiniuCheckFile(key string) (isExist bool, ret QiniuRet) {
	mac := qbox.NewMac(accessKey, secretKey)
	cfg := storage.Config{}
	bucketManager := storage.NewBucketManager(mac, &cfg)
	_, sErr := bucketManager.Stat(bucket, key)
	ret = QiniuRet{Key: key}
	if sErr != nil {
		return false, ret
	}
	ret.URL = storage.MakePublicURL(domain, ret.Key)
	return true, ret
}

// QiniuMakeURL make url by key.
func QiniuMakeURL(key string) string {
	return storage.MakePublicURL(domain, key)
}

// QiniuRemoveJSONFile delete json file for update
func QiniuRemoveJSONFile() {
	mac := qbox.NewMac(accessKey, secretKey)
	cfg := storage.Config{}
	bucketManager := storage.NewBucketManager(mac, &cfg)
	limit := 1000
	delimiter := ""
	//初始列举marker为空
	marker := ""
	for {
		entries, _, nextMarker, hasNext, err := bucketManager.ListFiles(bucket, jsonPrefix, delimiter, marker, limit)
		if err != nil {
			fmt.Println("list error,", err)
			break
		}
		keys := []string{}
		for _, entry := range entries {
			keys = append(keys, entry.Key)
		}
		deleteOps := make([]string, 0, len(keys))
		for _, key := range keys {
			deleteOps = append(deleteOps, storage.URIDelete(bucket, key))
		}
		rets, err := bucketManager.Batch(deleteOps)
		if err != nil {
			// 遇到错误
			if _, ok := err.(*rpc.ErrorInfo); ok {
				for _, ret := range rets {
					// 200 为成功
					fmt.Printf("%d\n", ret.Code)
					if ret.Code != 200 {
						fmt.Printf("%s\n", ret.Data.Error)
					}
				}
			} else {
				fmt.Printf("batch error, %s", err)
			}
		} else {
			// 完全成功
			for _, ret := range rets {
				// 200 为成功
				fmt.Printf("%d\n", ret.Code)
			}
		}

		if hasNext {
			marker = nextMarker
		} else {
			//list end
			break
		}
	}
}
