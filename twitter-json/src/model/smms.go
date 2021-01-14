package model

import (
	"os"
	"regexp"

	smms "github.com/sequix/smms-golang"
)

var token string = os.Getenv("SMMS_TOKEN")

// SmmsUpload upload
func SmmsUpload(filePath string, fileName string) (*smms.ImageRsp, error) {
	repeatReg := regexp.MustCompile(`statusCode: 200(.*?)code: image_repeated (.*?)this image exists at: (.*?)$`)
	client := smms.NewFromToken(token)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rsp, err := client.Upload(fileName, f)
	if err != nil && repeatReg.MatchString(err.Error()) {
		temp := smms.ImageRsp{}
		temp.URL = repeatReg.FindStringSubmatch(err.Error())[3]
		rsp = &temp
		err = nil
	}
	return rsp, err
}

// SmmsHistory get history
func SmmsHistory() ([]smms.ImageRsp, error) {
	client := smms.NewFromToken(token)
	rsp, err := client.History()
	return rsp, err
}

// SmmsDelete delete image
func SmmsDelete(hash string) error {
	client := smms.NewFromToken(token)
	err := client.Delete(hash)
	return err
}
