package model

import (
	"os"

	smms "github.com/sequix/smms-golang"
)

var token string = os.Getenv("SMMS_TOKEN")

// SmmsUpload upload
func SmmsUpload(filePath string, fileName string) (*smms.ImageRsp, error) {
	client := smms.NewFromToken(token)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rsp, err := client.Upload(fileName, f)
	return rsp, err
}
