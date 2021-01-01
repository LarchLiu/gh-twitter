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

// CheckDirExist check whether dir exist, if not make it
func CheckDirExist(name string) (exist bool) {
	_, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(name, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return false
			}
		} else {
			fmt.Println(err)
			return false
		}
	}
	return true
}

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

// GetImageInfo get image file path and file name.
func GetImageInfo(dirPath string, href string) (filePath string, fileName string, err error) {
	exist := CheckDirExist(dirPath)

	if !exist {
		return "", "", err
	}
	var suffix string
	var name string

	i := strings.LastIndex(href, "/")
	nameReg := regexp.MustCompile(`\/(.*?)(\.|\?)(?s)`)
	formatReg := regexp.MustCompile(`(?s)format=(.*?)&(?s)`)
	suffixReg := regexp.MustCompile(`\A((http)|(https))(.+)((.jpg)|(.png)|(.jpeg))\z`)

	if formatReg.MatchString(href) {
		// such like https://pbs.twimg.com/media/Eo9bCnnXcAAE_Wz?format=jpg&name=large
		name = nameReg.FindStringSubmatch(href[i:])[1]
		suffix = formatReg.FindStringSubmatch(href)[1]
		fileName = name + "." + suffix
		filePath = path.Join(dirPath, fileName)
		// fmt.Printf("filePath %s\n", filePath)
	} else if suffixReg.MatchString(href) {
		// such like https://pbs.twimg.com/profile_images/1333096463916797954/7bzarkH2_normal.jpg
		j := strings.LastIndex(href, "/") + 1
		fileName = href[j:]
		filePath = path.Join(dirPath, fileName)
		// fmt.Printf("fileName %s\n", filePath)
	} else {
		suffix = "jpg"
		filePath, fileName = randFileName(dirPath, suffix)
	}
	return filePath, fileName, nil
}

// WriteImage save image to dirPath.
func WriteImage(filePath string, href string) (err error) {
	imageFile, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("[writeImage create file]: filePath: %s\n href: %s\nerror: %s\n", filePath, href, err.Error())
		return err
	}
	resp, err := http.Get(href)
	if err != nil {
		fmt.Printf("[writeImage http.Get: filePath]: %s\nhref: %s\n error: %s\n", filePath, href, err.Error())
		os.Remove(filePath)
		return err
	}
	size, err := io.Copy(imageFile, resp.Body)
	if err != nil {
		fmt.Printf("io.Copy:error: %s  href: %s\n", err.Error(), href)
		os.Remove(filePath)
	}
	fmt.Printf("Get From %s: %d bytes\n", href, size)
	return nil
}

func randFileName(dirPath string, suffix string) (filePath string, fileName string) {
	var randChar string = "abcdefghimnjkqzxyt0123456789ioedkaldncalew0129387iue"
	var name string
	size := len(randChar)
	for i := 0; i < 12; i++ {
		name += string(randChar[rand.Intn(size)])
	}
	fileName = name + "." + suffix
	filePath = path.Join(dirPath, fileName)
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return filePath, fileName
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

// WriteJSONFile write json to file
func WriteJSONFile(dir string, fileName string, jsonBytes []byte) {
	exist := CheckDirExist(dir)
	if exist {
		file, er := os.OpenFile(dir+fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		defer func() { file.Close() }()
		if er != nil {
			fmt.Println(er)
		}
		_, err := file.Write(jsonBytes)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// WriteJSONFileByPath write json file by file path
func WriteJSONFileByPath(filePath string, jsonBytes []byte) {
	file, er := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer func() { file.Close() }()
	if er != nil {
		fmt.Println(er)
	}
	_, err := file.Write(jsonBytes)
	if err != nil {
		fmt.Println(err)
	}
}
