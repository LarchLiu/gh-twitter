package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

// GetFileContent get file content from filePath.
func GetFileContent(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	//读取文件全部内容
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// WriteImage save image to dirPath.
func WriteImage(dirPath string, href string) string {
	var suffix string
	var name string
	var fileName string
	i := strings.LastIndex(href, "/")
	nameReg := regexp.MustCompile(`\/(.*?)(\.|\?)(?s)`)
	formatReg := regexp.MustCompile(`(?s)format=(.*?)&(?s)`)
	suffixReg := regexp.MustCompile(`\A((http)|(https))(.+)((.jpg)|(.png)|(.jpeg))\z`)

	if formatReg.MatchString(href) {
		// such like https://pbs.twimg.com/media/Eo9bCnnXcAAE_Wz?format=jpg&name=large
		name = nameReg.FindStringSubmatch(href[i:])[1]
		suffix = formatReg.FindStringSubmatch(href)[1]
		fileName = path.Join(dirPath, name+"."+suffix)
		fmt.Printf("fileName %s\n", fileName)
		_, err := os.Stat(fileName)
		if err == nil {
			return fileName
		}
	} else if suffixReg.MatchString(href) {
		// such like https://pbs.twimg.com/profile_images/1333096463916797954/7bzarkH2_normal.jpg
		j := strings.LastIndex(href, "/") + 1
		name = href[j:]
		fileName = path.Join(dirPath, name)
		_, err := os.Stat(fileName)
		fmt.Printf("fileName %s\n", fileName)
		if err == nil {
			return fileName
		}
	} else {
		suffix = "jpg"
		fileName = randFileName(dirPath, suffix)
	}
	imageFile, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("[writeImage create file]: fileName: %s\n href: %s\nerror: %s\n", fileName, href, err.Error())
		return ""
	}
	resp, err := http.Get(href)
	if err != nil {
		fmt.Printf("[writeImage http.Get: fileName]: %s\nhref: %s\n error: %s\n", fileName, href, err.Error())
		os.Remove(fileName)
		return ""
	}
	size, err := io.Copy(imageFile, resp.Body)
	if err != nil {
		fmt.Printf("io.Copy:error: %s  href: %s\n", err.Error(), href)
		os.Remove(fileName)
	}
	fmt.Printf("Get From %s: %d bytes\n", href, size)
	return fileName
}

func randFileName(dirPath string, suffix string) (fileName string) {
	var randChar string = "abcdefghimnjkqzxyt0123456789ioedkaldncalew0129387iue"
	var name string
	size := len(randChar)
	for i := 0; i < 12; i++ {
		name += string(randChar[rand.Intn(size)])
	}
	fileName = path.Join(dirPath, name)
	fileName += "." + suffix
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return fileName
		}
	}
	return randFileName(dirPath, suffix)
}

func getCorrectHref(href string, url string) (correctHref string, exist bool) {
	completeValidator := regexp.MustCompile(`\A((http)|(https))(.+)((.jpg)|(.png)|(.jpeg))\z`)
	imcompleteValidator := regexp.MustCompile(`\A(||)(.+)((.jpg)|(.jpeg)|(.png))\z`)
	if completeValidator.MatchString(href) {
		//fmt.Println(href)
		return href, true
	} else if imcompleteValidator.MatchString(href) {
		ret := href[1:]
		ret = url + ret
		//fmt.Println(ret)
		return string(ret), true
	}
	return "", false
}
