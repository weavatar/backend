package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

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
func (r *UserController) OauthLogin(ctx http.Context) {
	state, stateErr := oauth.GetAuthorizationState(ctx.Request().Ip())
	if stateErr != nil {
		facades.Log().Error("[UserController][OauthLogin] 获取State失败 ", stateErr.Error())
		Error(ctx, http.StatusInternalServerError, stateErr.Error())
		return
	}

	url := facades.Config().GetString("haozi.account.base_url") + "/oauth/authorize?client_id=" + facades.Config().GetString("haozi.account.client_id") + "&redirect_uri=" + facades.Config().GetString("http.url") + "/oauth/callback&response_type=code&scope=basic&state=" + state

	Success(ctx, http.Json{
		"url": url,
	})
}

func (r *UserController) OauthCallback(ctx http.Context) {
	var oauthCallbackRequest requests.OauthCallbackRequest
	if !Sanitize(ctx, &oauthCallbackRequest) {
		return
	}

	// 验证 state
	if facades.Cache().GetString("oauth_state:"+oauthCallbackRequest.State) != ctx.Request().Ip() {
		Error(ctx, http.StatusUnprocessableEntity, "状态码不存在或已过期")
		return
	}

	// 获取 token
	oauthToken, tokenErr := oauth.GetToken(oauthCallbackRequest.Code)
	if tokenErr != nil {
		facades.Log().Error("[UserController][OauthCallback] 获取 access_token 失败 ", tokenErr.Error())
		Error(ctx, http.StatusInternalServerError, "获取 access_token 失败")
		return
	}

	// 获取用户信息
	userInfo, userInfoErr := oauth.GetUserInfo(oauthToken.AccessToken)
	if userInfoErr != nil {
		facades.Log().Error("[UserController][OauthCallback] 获取用户信息失败 ", userInfoErr.Error())
		Error(ctx, http.StatusInternalServerError, "获取用户信息失败")
		return
	}

	userID, idErr := id.NewRatID().Generate()
	if idErr != nil {
		facades.Log().Error("[UserController][OauthCallback] 生成用户ID失败 ", idErr.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	// 检查用户是否存在
	var user models.User
	if err := facades.Orm().Query().FirstOrCreate(&user, models.User{OpenID: userInfo.Data.OpenID}, models.User{ID: userID, UnionID: userInfo.Data.UnionID, Nickname: userInfo.Data.Nickname, Avatar: "https://weavatar.com/avatar/?d=mp"}); err != nil {
		facades.Log().Error("[UserController][OauthCallback] 查询用户失败 ", err.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	// 登录
	token, loginErr := facades.Auth().LoginUsingID(ctx, user.ID)
	if loginErr != nil {
		facades.Log().Error("[UserController][OauthCallback] 登录失败 ", loginErr.Error())
		Error(ctx, http.StatusInternalServerError, loginErr.Error())
		return
	}

	Success(ctx, http.Json{
		"token": token,
	})
}

func (r *UserController) GetInfo(ctx http.Context) {
	var user models.User
	if err := facades.Auth().User(ctx, &user); err != nil {
		Error(ctx, http.StatusUnauthorized, "登录已过期")
		return
	}

	Success(ctx, http.Json{
		"id":         user.ID,
		"nickname":   user.Nickname,
		"avatar":     user.Avatar,
		"created_at": user.CreatedAt,
	})
}

func (r *UserController) UpdateInfo(ctx http.Context) {
	var updateInfoRequest requests.UpdateInfoRequest
	if !Sanitize(ctx, &updateInfoRequest) {
		return
	}

	var user models.User
	if err := facades.Auth().User(ctx, &user); err != nil {
		Error(ctx, http.StatusUnauthorized, "登录已过期")
		return
	}

	user.Nickname = updateInfoRequest.Nickname
	user.Avatar = updateInfoRequest.Avatar
	if err := facades.Orm().Query().Save(&user); err != nil {
		facades.Log().Error("[UserController][UpdateNickname] 更新用户失败 ", err.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, nil)
}

func (r *UserController) Logout(ctx http.Context) {
	if err := facades.Auth().Logout(ctx); err != nil {
		facades.Log().Error("[UserController][Logout] 登出失败 ", err.Error())
		Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	Success(ctx, nil)
}
