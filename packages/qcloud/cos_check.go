package qcloud

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/imroc/req/v3"
)

type CosChecker struct {
	accessKey, secretKey, bucket string
}

type CheckResponse struct {
	XMLName  xml.Name `xml:"RecognitionResult"`
	JobId    string   `xml:"JobId"`
	Result   int      `xml:"Result"`
	Label    string   `xml:"Label"`
	SubLabel string   `xml:"SubLabel"`
	Score    int      `xml:"Score"`
}

// once 单例模式
var once sync.Once

// internal 内部使用的对象
var internal *CosChecker

func NewCreator(accessKey, secretKey, bucket string) *CosChecker {
	once.Do(func() {
		internal = &CosChecker{accessKey: accessKey, secretKey: secretKey, bucket: bucket}
	})
	return internal
}

// Check 检查图片是否违规
func (cc *CosChecker) Check(url string) (bool, error) {
	authorization, err := cc.getAuthorization("GET", "/", 0)
	if err != nil {
		facades.Log().Error("COS审核 ", " [获取签名失败] "+err.Error())
		return false, err
	}

	client := req.C()
	resp, reqErr := client.R().SetQueryParams(map[string]string{
		"ci-process": "sensitive-content-recognition",
		"detect-url": url,
	}).SetHeader("Authorization", authorization).Get("https://" + cc.bucket + "/")
	if !resp.IsSuccessState() {
		if reqErr != nil {
			facades.Log().Error("COS审核 ", " [请求失败] URL:"+url+" "+reqErr.Error())
			return false, reqErr
		} else {
			facades.Log().Error("COS审核 ", " [请求失败] URL:"+url+" "+resp.String())
			return false, errors.New("COS审核[请求失败]")
		}
	}

	var checkResponse CheckResponse
	err = xml.Unmarshal(resp.Bytes(), &checkResponse)
	if err != nil {
		facades.Log().Error("COS审核 ", " [响应解析失败] "+err.Error())
		return false, err
	}

	if checkResponse.Result == 1 {
		return false, nil
	}

	return true, nil
}

// getAuthorization 获取签名
func (cc *CosChecker) getAuthorization(method, path string, expires time.Duration) (string, error) {
	if expires <= 0 {
		expires = 30 * time.Minute
	}
	signTimeStart := time.Now().Add(-time.Minute).Unix()
	signTimeEnd := time.Now().Add(expires).Unix()
	signTime := fmt.Sprintf("%d;%d", signTimeStart, signTimeEnd)

	pathUnescaped, err := url.PathUnescape(path)
	if err != nil {
		return "", err
	}

	httpString := strings.ToLower(method) + "\n" + pathUnescaped + "\n\n\n"
	hasher := sha1.New()
	hasher.Write([]byte(httpString))
	sha1edHttpString := hex.EncodeToString(hasher.Sum(nil))
	stringToSign := "sha1\n" + signTime + "\n" + sha1edHttpString + "\n"

	h := hmac.New(sha1.New, []byte(cc.secretKey))
	h.Write([]byte(signTime))
	signKey := hex.EncodeToString(h.Sum(nil))

	h2 := hmac.New(sha1.New, []byte(signKey))
	h2.Write([]byte(stringToSign))
	signature := hex.EncodeToString(h2.Sum(nil))

	return fmt.Sprintf("q-sign-algorithm=sha1&q-ak=%s&q-sign-time=%s&q-key-time=%s&q-header-list=&q-url-param-list=&q-signature=%s",
		cc.accessKey, signTime, signTime, signature), nil
}
