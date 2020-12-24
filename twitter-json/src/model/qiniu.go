package model

import (
	"context"
	"fmt"
	"os"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
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

// QiniuUpload qiniu api to upload file.
func QiniuUpload(filePath string, key string) (ret QiniuRet, err error) {
	mac := qbox.NewMac(accessKey, secretKey)
	putPolicy := storage.PutPolicy{
		Scope:      bucket,
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
