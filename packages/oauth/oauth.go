package oauth

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/imroc/req/v3"

	"weavatar/packages/helpers"
)

type BasicInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Nickname string `json:"nickname"`
		OpenID   string `json:"open_id"`
		UnionID  string `json:"union_id"`
	} `json:"data"`
}

type Token struct {
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int      `json:"expires_in"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

// GetAuthorizationState 获取授权状态码
func GetAuthorizationState(ip string) (string, error) {
	state := helpers.RandomString(32)
	err := facades.Cache.Put("oauth_state:"+state, ip, 10*time.Minute)
	if err != nil {
		return "", err
	}
	return state, nil
}

// GetToken 获取 AccessToken 和 RefreshToken 信息
func GetToken(code string) (map[string]string, error) {
	clientID := facades.Config.GetString("haozi.account.client_id")
	clientSecret := facades.Config.GetString("haozi.account.client_secret")
	redirectURI := facades.Config.GetString("http.url") + "/oauth/callback"

	client := req.C()
	resp, err := client.R().SetQueryParams(map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
		"redirect_uri":  redirectURI,
	}).Get(facades.Config.GetString("haozi.account.base_url") + "/api/oauth/token")
	if err != nil {
		facades.Log.Warning("耗子通行证 ", " [获取Token失败] ", err.Error())
		return nil, err
	}

	// 解析Token
	var token Token
	err = json.Unmarshal([]byte(resp.String()), &token)
	if err != nil {
		return nil, err
	}

	// 判断ExpiresIn是否为0
	if token.ExpiresIn == 0 {
		return nil, errors.New("获取Token失败")
	}

	retMap := make(map[string]string)
	retMap["access_token"] = token.AccessToken
	retMap["refresh_token"] = token.RefreshToken
	retMap["expires_in"] = strconv.Itoa(token.ExpiresIn)

	return retMap, nil
}

// GetUserInfo 获取用户信息
func GetUserInfo(accessToken string) (map[string]string, error) {
	client := req.C()
	resp, err := client.R().SetQueryParams(map[string]string{
		"access_token": accessToken,
	}).Get(facades.Config.GetString("haozi.account.base_url") + "/api/oauth/getBasicInfo")
	if err != nil {
		facades.Log.Warning("耗子通行证 ", " [获取用户信息失败] ", err.Error())
	}

	var basicInfo BasicInfo
	err = json.Unmarshal([]byte(resp.String()), &basicInfo)
	if err != nil {
		return nil, err
	}

	if basicInfo.Code != 0 {
		return nil, errors.New(basicInfo.Message)
	}

	retMap := make(map[string]string)
	retMap["nickname"] = basicInfo.Data.Nickname
	retMap["open_id"] = basicInfo.Data.OpenID
	retMap["union_id"] = basicInfo.Data.UnionID

	return retMap, nil
}
