package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

type UserController struct {
	//Dependent services
}

func NewUserController() *UserController {
	return &UserController{
		//Inject services
	}
}

func (r *UserController) OauthLogin(ctx http.Context) {

}

func (r *UserController) OauthCallback(ctx http.Context) {

}

func (r *UserController) UpdateNickname(ctx http.Context) {

}

func (r *UserController) Logout(ctx http.Context) {

}

func (r *UserController) Refresh(ctx http.Context) {

}
