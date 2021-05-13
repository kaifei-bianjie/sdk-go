package ad

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/irisnet/sdk-go/util"
	"strings"
	"testing"
)

const (
	ossEndPoint     = "xxx"
	ossBucketName   = "xxx"
	ossAccessKey    = "xxx"
	ossAccessSecret = "xxx"
)

var (
	ossBucket *oss.Bucket
)

func initOssClient() {
	client, err := oss.New(ossEndPoint, ossAccessKey, ossAccessSecret)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 获取存储空间。
	c, err := client.Bucket(ossBucketName)
	if err != nil {
		panic(err)
	}
	ossBucket = c
}

func TestUploadFileIcon(t *testing.T) {
	initOssClient()
	iconNameFilePathMap := map[string]string{
		"icon_badge_bronze.png": "./source/icon_badge_bronze.png",
		"icon_badge_gold.png":   "./source/icon_badge_gold.png",
		"icon_badge_silver.png": "./source/icon_badge_silver.png",
	}

	option := oss.ContentType("image/jpg")
	objectAcl := oss.ObjectACL(oss.ACLPublicRead)

	if len(iconNameFilePathMap) > 0 {
		for objName, filePath := range iconNameFilePathMap {
			err := ossBucket.PutObjectFromFile(objName, filePath, option, objectAcl)
			if err != nil {
				t.Errorf("%s, Error:%s", filePath, err)
				return
			}
			t.Logf("%s:%s\n", filePath, genOssUrl(objName))
		}
	}
}

func TestUploadFileSvg(t *testing.T) {
	initOssClient()
	objectName := "gold.svg"
	fmt.Println(objectName)
	option := oss.ContentType("text/xml")
	objectAcl := oss.ObjectACL(oss.ACLPublicRead)

	err := ossBucket.PutObjectFromFile(objectName, "./source/gold.svg", option, objectAcl)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func TestUploadString(t *testing.T) {
	initOssClient()
	objectName := "goldString.svg"
	fmt.Println(objectName)
	option := oss.ContentType("text/xml")
	objectAcl := oss.ObjectACL(oss.ACLPublicRead)

	if c, err := util.ReadFile("./source/gold.svg"); err != nil {
		fmt.Println(err)
	} else {
		err = ossBucket.PutObject(objectName, strings.NewReader(c), option, objectAcl)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(genOssUrl(objectName))
	}
}

//生成外网访问url
func genOssUrl(objectKey string) string {
	return fmt.Sprintf("https://%s.%s/%s", ossBucketName, ossEndPoint, objectKey)
}
