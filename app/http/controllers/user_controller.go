package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	requests "weavatar/app/http/requests/user"
	"weavatar/app/models"
	"weavatar/packages/id"
	"weavatar/packages/oauth"
)

type UserController struct {
	//Dependent services
}

func NewUserController() *UserController {
	return &UserController{
		//Inject services
	}
}

// OauthLogin 通行证登录
func (r *UserController) OauthLogin(ctx http.Context) {
	state, stateErr := oauth.GetAuthorizationState(ctx.Request().Ip())
	if stateErr != nil {
		facades.Log.WithContext(ctx).Error("[UserController][OauthLogin] 获取State失败 ", stateErr)
		ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"code":    500,
			"message": stateErr.Error(),
		})
		return
	}

	url := facades.Config.GetString("haozi.account.base_url") + "/oauth/authorize?client_id=" + facades.Config.GetString("haozi.account.client_id") + "&redirect_uri=" + facades.Config.GetString("http.url") + "/oauth/callback&response_type=code&scope=basic&state=" + state

	ctx.Response().Success().Json(http.Json{
		"code":    0,
		"message": "获取成功",
		"data": http.Json{
			"url": url,
		},
	})
}

func (r *UserController) OauthCallback(ctx http.Context) {
	var oauthCallbackRequest requests.OauthCallbackRequest
	errors, err := ctx.Request().ValidateRequest(&oauthCallbackRequest)
	if err != nil {
		ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"code":    422,
			"message": err.Error(),
		})
		return
	}
	if errors != nil {
		ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"code":    422,
			"message": errors.All(),
		})
		return
	}

	// 验证 state
	if facades.Cache.GetString("oauth_state:"+oauthCallbackRequest.State, "1.1.1.1") != ctx.Request().Ip() {
		ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"code":    422,
			"message": "状态码不存在或已过期",
		})
		return
	}

	// 获取 token
	tokenMap, tokenErr := oauth.GetToken(oauthCallbackRequest.Code)
	if tokenErr != nil {
		facades.Log.WithContext(ctx).Error("[UserController][OauthCallback] 获取 access_token 失败 ", tokenErr)
		ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"code":    500,
			"message": "获取 access_token 失败",
		})
		return
	}
	accessToken := tokenMap["access_token"]

	// 获取用户信息
	userInfo, userInfoErr := oauth.GetUserInfo(accessToken)
	if userInfoErr != nil {
		facades.Log.WithContext(ctx).Error("[UserController][OauthCallback] 获取用户信息失败 ", userInfoErr)
		ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"code":    500,
			"message": "获取用户信息失败",
		})
		return
	}

	userID, idErr := id.NewRatID().Generate()
	if idErr != nil {
		facades.Log.WithContext(ctx).Error("[UserController][OauthCallback] 生成用户ID失败 ", idErr)
		ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"code":    500,
			"message": "系统内部错误",
		})
		return
	}

	// 检查用户是否存在
	var user models.User
	userErr := facades.Orm.Query().FirstOrCreate(&user, models.User{OpenID: userInfo["open_id"]}, models.User{ID: userID, UnionID: userInfo["union_id"], Nickname: userInfo["nickname"]})
	if userErr != nil {
		facades.Log.WithContext(ctx).Error("[UserController][OauthCallback] 查询用户失败 ", userErr)
		ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"code":    500,
			"message": "系统内部错误",
		})
		return
	}

	// 登录
	token, loginErr := facades.Auth.LoginUsingID(ctx, user.ID)
	if loginErr != nil {
		facades.Log.WithContext(ctx).Error("[UserController][OauthCallback] 登录失败 ", err)
		ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"code":    500,
			"message": loginErr.Error(),
		})
		return
	}

	ctx.Response().Success().Json(http.Json{
		"code":    0,
		"message": "登录成功",
		"data": http.Json{
			"token": token,
		},
	})
}

func (r *UserController) Info(ctx http.Context) {
	// 取出用户信息
	user, ok := ctx.Value("user").(models.User)
	if !ok {
		ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
			"code":    401,
			"message": "登录已过期",
		})
		return
	}

	ctx.Response().Success().Json(http.Json{
		"code":    0,
		"message": "获取成功",
		"data": http.Json{
			"id":         user.ID,
			"nickname":   user.Nickname,
			"created_at": user.CreatedAt,
		},
	})
}

func (r *UserController) UpdateNickname(ctx http.Context) {
	var updateNicknameRequest requests.UpdateNicknameRequest
	errors, err := ctx.Request().ValidateRequest(&updateNicknameRequest)
	if err != nil {
		ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"code":    422,
			"message": err.Error(),
		})
		return
	}
	if errors != nil {
		ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"code":    422,
			"message": errors.All(),
		})
		return
	}

	// 取出用户信息
	user, ok := ctx.Value("user").(models.User)
	if !ok {
		ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
			"code":    401,
			"message": "登录已过期",
		})
		return
	}

	user.Nickname = updateNicknameRequest.Nickname
	updateErr := facades.Orm.Query().Save(&user)
	if updateErr != nil {
		facades.Log.WithContext(ctx).Error("[UserController][UpdateNickname] 更新用户失败 ", updateErr)
		ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"code":    500,
			"message": "系统内部错误",
		})
		return
	}

	ctx.Response().Success().Json(http.Json{
		"code":    0,
		"message": "修改成功",
	})
}

func (r *UserController) Logout(ctx http.Context) {
	err := facades.Auth.Logout(ctx)
	if err != nil {
		facades.Log.WithContext(ctx).Error("[UserController][Logout] 登出失败 ", err)
		ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	ctx.Response().Success()
}

func (r *UserController) Refresh(ctx http.Context) {
	token, err := facades.Auth.Refresh(ctx)
	if err != nil {
		facades.Log.WithContext(ctx).Error("[UserController][Refresh] 刷新令牌失败 ", err)
		ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	ctx.Response().Success().Json(http.Json{
		"code":    0,
		"message": "刷新成功",
		"data": http.Json{
			"token": token,
		},
	})
}
