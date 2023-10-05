package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"

	"weavatar/pkg/wangsu/common/constant"
	"weavatar/pkg/wangsu/common/model"
	"weavatar/pkg/wangsu/common/util"
)

type AkskConfig struct {
	AccessKey     string
	SecretKey     string
	Uri           string
	Method        string
	EndPoint      string
	SignedHeaders string // 参与计算的头部
}

func TransferRequestMsg(config AkskConfig) model.HttpRequestMsg {
	var requestMsg = model.HttpRequestMsg{Params: map[string]string{}, Headers: map[string]string{}}
	requestMsg.Uri = config.Uri
	requestMsg.Method = config.Method
	requestMsg.Url = constant.HttpRequestPrefix + config.Uri
	if len(config.EndPoint) == 0 || "{endPoint}" == config.EndPoint {
		requestMsg.Host = constant.HttpRequestDomain
		requestMsg.Url = "https://" + constant.HttpRequestDomain + config.Uri
	} else {
		requestMsg.Host = config.EndPoint
		requestMsg.Url = "https://" + config.EndPoint + config.Uri
	}
	requestMsg.SignedHeaders = getSignedHeaders(config.SignedHeaders)
	return requestMsg
}

func Invoke(config AkskConfig, jsonStr string, params ...string) string {
	var requestMsg = TransferRequestMsg(config)

	if config.Method == "POST" || config.Method == "PUT" || config.Method == "PATCH" || config.Method == "DELETE" {
		requestMsg.Body = jsonStr
	}

	if len(params) > 0 {
		decodeParams := make(map[string]string)
		err := sonic.UnmarshalString(params[0], &decodeParams)
		if err == nil {
			requestMsg.Params = decodeParams
		}
	}

	timeStamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	requestMsg.Headers[constant.HeadSignTimeStamp] = timeStamp
	requestMsg.Headers["Host"] = requestMsg.Host
	requestMsg.Headers[constant.ContentType] = constant.ApplicationJson
	requestMsg.Headers[constant.HeadSignAccessKey] = config.AccessKey
	requestMsg.Headers[constant.XCncAuthMethod] = constant.AKSK
	signature := getSignature(requestMsg, config.SecretKey, timeStamp)
	requestMsg.Headers[constant.Authorization] = genAuthorization(config.AccessKey, requestMsg.SignedHeaders, signature)

	return util.Call(requestMsg)
}

// 拼接最后签名
func genAuthorization(accessKey string, signedHeaders string, signature string) string {
	var build strings.Builder
	build.WriteString(constant.HeadSignAlgorithm)
	build.WriteString(" ")
	build.WriteString("Credential=")
	build.WriteString(accessKey)
	build.WriteString(", ")
	build.WriteString("SignedHeaders=")
	build.WriteString(signedHeaders)
	build.WriteString(", ")
	build.WriteString("Signature=")
	build.WriteString(signature)
	return build.String()
}

func getSignature(requestMsg model.HttpRequestMsg, secretKey string, timeStamp string) string {
	var bodyStr = requestMsg.Body
	if len(requestMsg.Body) == 0 || "GET" == requestMsg.Method {
		bodyStr = ""
	}
	hashedRequestPayload := hmacSha256(bodyStr)
	canonicalRequest := requestMsg.Method + "\n" +
		getRequestUri(requestMsg) + "\n" +
		getQueryString(requestMsg) + "\n" +
		getCanonicalHeaders(requestMsg.Headers, requestMsg.SignedHeaders) + "\n" +
		getSignedHeaders(requestMsg.SignedHeaders) + "\n" +
		hashedRequestPayload
	stringToSign := constant.HeadSignAlgorithm + "\n" + timeStamp + "\n" + hmacSha256(canonicalRequest)
	return hmac256(secretKey, stringToSign)
}

// 获取uri
func getRequestUri(requestMsg model.HttpRequestMsg) string {
	indexOfQueryStringSeparator := strings.Index(requestMsg.Uri, "?")
	if indexOfQueryStringSeparator == -1 {
		return requestMsg.Uri
	}
	return string([]rune(requestMsg.Uri)[:indexOfQueryStringSeparator])
}

// 获取uri参数
func getQueryString(requestMsg model.HttpRequestMsg) string {
	indexOfQueryStringSeparator := strings.Index(requestMsg.Uri, "?")
	if "POST" == requestMsg.Method || indexOfQueryStringSeparator == -1 {
		return ""
	}
	s, err := url.QueryUnescape(requestMsg.Uri[indexOfQueryStringSeparator+1 : len(requestMsg.Uri)])
	if err != nil {
		fmt.Println("decode请求参数失败.")
	}
	return s
}

// 获取并排序参与签名计算的头部
func getSignedHeaders(signedHeaders string) string {
	if len(signedHeaders) == 0 {
		return "content-type;host"
	}
	headers := strings.Split(strings.ToLower(signedHeaders), ";")
	sort.Strings(headers)
	return strings.Join(headers, ";")
}

// 获取k-v字符串
func getCanonicalHeaders(headerMap map[string]string, signedHeaders string) string {
	keys := strings.Split(signedHeaders, ";")
	var headers = make(map[string]string)
	for k, v := range headerMap {
		headers[strings.ToLower(k)] = v
	}
	var build strings.Builder
	for i := 0; i < len(keys); i++ {
		build.WriteString(keys[i])
		build.WriteString(":")
		build.WriteString(strings.ToLower(headers[keys[i]]))
		build.WriteString("\n")
	}
	return build.String()
}

// 加密算法
func hmacSha256(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	hashCode := hash.Sum(nil)
	result := hex.EncodeToString(hashCode)
	return strings.ToLower(result)
}

func hmac256(secretKey string, stringToSign string) string {
	value := []byte(secretKey)
	key := hmac.New(sha256.New, value)
	key.Write([]byte(stringToSign))
	result := hex.EncodeToString(key.Sum(nil))
	return strings.ToLower(result)
}
