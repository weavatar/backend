package oauth

import (
	"errors"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/imroc/req/v3"

	"weavatar/pkg/helper"
)

type BasicInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Nickname  string `json:"nickname"`
		OpenID    string `json:"open_id"`
		UnionID   string `json:"union_id"`
		PhoneBind bool   `json:"phone_bind"`
		RealName  bool   `json:"real_name"`
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
	state := helper.RandomString(32)
	err := facades.Cache().Put("oauth_state:"+state, ip, 10*time.Minute)
	if err != nil {
		return "", err
	}
	return state, nil
}

// GetToken 获取 AccessToken 和 RefreshToken 信息
func GetToken(code string) (Token, error) {
	clientID := facades.Config().GetString("haozi.account.client_id")
	clientSecret := facades.Config().GetString("haozi.account.client_secret")
	redirectURI := facades.Config().GetString("http.url") + "/oauth/callback"

	var token Token

	client := req.C()
	resp, err := client.R().SetQueryParams(map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
		"redirect_uri":  redirectURI,
	}).SetSuccessResult(&token).Get(facades.Config().GetString("haozi.account.base_url") + "/api/oauth/token")
	if err != nil {
		return token, err
	}
	if !resp.IsSuccessState() {
		return token, errors.New("获取Token失败: " + resp.String())
	}

	if token.ExpiresIn == 0 {
		return token, errors.New("获取Token失败: " + resp.String())
	}

	return token, nil
}

// GetUserInfo 获取用户信息
func GetUserInfo(accessToken string) (BasicInfo, error) {
	var basicInfo BasicInfo

	client := req.C()
	resp, err := client.R().SetQueryParams(map[string]string{
		"access_token": accessToken,
	}).SetSuccessResult(basicInfo).Get(facades.Config().GetString("haozi.account.base_url") + "/api/oauth/getBasicInfo")
	if err != nil {
		return basicInfo, err
	}
	if !resp.IsSuccessState() {
		return basicInfo, errors.New("获取用户信息失败: " + resp.String())
	}

	if basicInfo.Code != 0 {
		return basicInfo, errors.New(basicInfo.Message)
	}

	return basicInfo, nil
}
