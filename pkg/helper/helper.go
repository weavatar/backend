// Package helper 存放辅助方法
package helper

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"net/url"
	"reflect"
)

// Empty 类似于 PHP 的 empty() 函数
func Empty(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

// FirstElement 安全获取 args[0]
func FirstElement(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

// RandomNumber 生成长度为 length 随机数字字符串
func RandomNumber(length int) string {
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

// RandomString 生成长度为 length 的随机字符串
func RandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	letters := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i, v := range b {
		b[i] = letters[v%byte(len(letters))]
	}
	return string(b)
}

// MD5 生成字符串的 MD5 值
func MD5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

// IsURL 判断字符串是否是 URL
func IsURL(str string) bool {
	unescape, err := url.QueryUnescape(str)
	if err == nil {
		str = unescape
	}
	parsed, err := url.Parse(str)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || len(parsed.Host) == 0 {
		return false
	}

	return true
}
