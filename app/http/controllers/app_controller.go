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

func (r *AppController) Create(ctx http.Context) {

}

func (r *AppController) Get(ctx http.Context) {

}

func (r *AppController) GetSingle(ctx http.Context) {

}

func (r *AppController) Update(ctx http.Context) {

}

func (r *AppController) Delete(ctx http.Context) {

}
