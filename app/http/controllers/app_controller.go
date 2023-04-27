package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

type AppController struct {
	//Dependent services
}

func NewAppController() *AppController {
	return &AppController{
		//Inject services
	}
}

// Index 获取应用列表
func (r *AppController) Index(ctx http.Context) {

}

// Show 获取应用详情
func (r *AppController) Show(ctx http.Context) {

}

// Store 创建应用
func (r *AppController) Store(ctx http.Context) {

}

// Update 更新应用
func (r *AppController) Update(ctx http.Context) {

}

// Destroy 删除应用
func (r *AppController) Destroy(ctx http.Context) {

}

// Delete 删除应用
func (r *AppController) Delete(ctx http.Context) {

}
