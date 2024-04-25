package controllers

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/imroc/req/v3"

	requests "weavatar/app/http/requests/user"
	"weavatar/app/models"
	"weavatar/pkg/id"
	"weavatar/pkg/oauth"
)

type UserController struct {
	// Dependent services
}

func NewUserController() *UserController {
	return &UserController{
		// Inject services
	}
}

// OauthLogin 通行证登录
func (r *UserController) OauthLogin(ctx http.Context) http.Response {
	state, stateErr := oauth.GetAuthorizationState(ctx.Request().Ip())
	if stateErr != nil {
		facades.Log().Error("[UserController][OauthLogin] 获取State失败 ", stateErr.Error())
		return Error(ctx, http.StatusInternalServerError, stateErr.Error())
	}

	url := facades.Config().GetString("haozi.account.base_url") + "/oauth/authorize?client_id=" + facades.Config().GetString("haozi.account.client_id") + "&redirect_uri=" + facades.Config().GetString("http.url") + "/oauth/callback&response_type=code&scope=basic&state=" + state

	return Success(ctx, http.Json{
		"url": url,
	})
}

func (r *UserController) OauthCallback(ctx http.Context) http.Response {
	var oauthCallbackRequest requests.OauthCallbackRequest
	sanitize := Sanitize(ctx, &oauthCallbackRequest)
	if sanitize != nil {
		return sanitize
	}

	// 验证 state
	if facades.Cache().GetString("oauth_state:"+oauthCallbackRequest.State) != ctx.Request().Ip() {
		return Error(ctx, http.StatusUnprocessableEntity, "状态码不存在或已过期")
	}

	// 获取 token
	oauthToken, tokenErr := oauth.GetToken(oauthCallbackRequest.Code)
	if tokenErr != nil {
		facades.Log().Error("[UserController][OauthCallback] 获取 access_token 失败 ", tokenErr.Error())
		return Error(ctx, http.StatusInternalServerError, "获取 access_token 失败")
	}

	// 获取用户信息
	userInfo, userInfoErr := oauth.GetUserInfo(oauthToken.AccessToken)
	if userInfoErr != nil {
		facades.Log().Error("[UserController][OauthCallback] 获取用户信息失败 ", userInfoErr.Error())
		return Error(ctx, http.StatusInternalServerError, "获取用户信息失败")
	}

	var user models.User
	if err := facades.Orm().Query().FirstOrCreate(&user, models.User{OpenID: userInfo.Data.OpenID, UnionID: userInfo.Data.UnionID}, models.User{ID: id.Generate(), Nickname: userInfo.Data.Nickname, Avatar: "https://weavatar.com/avatar/?d=mp"}); err != nil {
		facades.Log().Error("[UserController][OauthCallback] 查询用户失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}
	user.RealName = userInfo.Data.RealName
	if err := facades.Orm().Query().Save(&user); err != nil {
		facades.Log().Error("[UserController][OauthCallback] 更新用户失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	token, loginErr := facades.Auth(ctx).LoginUsingID(user.ID)
	if loginErr != nil {
		facades.Log().Error("[UserController][OauthCallback] 登录失败 ", loginErr.Error())
		return Error(ctx, http.StatusInternalServerError, loginErr.Error())
	}

	return Success(ctx, http.Json{
		"token": token,
	})
}

func (r *UserController) GetInfo(ctx http.Context) http.Response {
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return Error(ctx, http.StatusUnauthorized, "登录已过期")
	}

	return Success(ctx, http.Json{
		"id":         user.ID,
		"nickname":   user.Nickname,
		"avatar":     user.Avatar,
		"real_name":  user.RealName,
		"created_at": user.CreatedAt,
	})
}

func (r *UserController) UpdateInfo(ctx http.Context) http.Response {
	var updateInfoRequest requests.UpdateInfoRequest
	sanitize := Sanitize(ctx, &updateInfoRequest)
	if sanitize != nil {
		return sanitize
	}

	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return Error(ctx, http.StatusUnauthorized, "登录已过期")
	}

	user.Nickname = updateInfoRequest.Nickname
	user.Avatar = updateInfoRequest.Avatar
	if err := facades.Orm().Query().Save(&user); err != nil {
		facades.Log().Error("[UserController][UpdateNickname] 更新用户失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	return Success(ctx, nil)
}

func (r *UserController) Logout(ctx http.Context) http.Response {
	if err := facades.Auth(ctx).Logout(); err != nil {
		facades.Log().Error("[UserController][Logout] 登出失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, nil)
}

// GetQQAvatar 获取 QQ 头像
func (r *UserController) GetQQAvatar(ctx http.Context) http.Response {
	client := req.C()
	resp, err := client.R().SetQueryParams(map[string]string{
		"b":  "qq",
		"nk": ctx.Request().Input("qq"),
		"s":  "640",
	}).Get("http://q1.qlogo.cn/g")

	length, lengthErr := strconv.Atoi(resp.GetHeader("Content-Length"))
	if length < 6400 || lengthErr != nil {
		resp, err = client.R().SetQueryParams(map[string]string{
			"b":  "qq",
			"nk": ctx.Request().Input("qq"),
			"s":  "100",
		}).Get("http://q1.qlogo.cn/g")
	}

	if err != nil || !resp.IsSuccessState() {
		return Error(ctx, http.StatusInternalServerError, "获取失败请检查输入")
	}

	img, _, err := image.Decode(bytes.NewReader(resp.Bytes()))
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, "解析头像图片失败")
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, "解析头像图片失败")
	}

	return Success(ctx, base64.StdEncoding.EncodeToString(buf.Bytes()))
}
