package providers

import (
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"

	"weavatar/app/rules"
)

type ValidationServiceProvider struct {
}

func (receiver *ValidationServiceProvider) Register() {

}

func (receiver *ValidationServiceProvider) Boot() {
	if err := facades.Validation.AddRules(receiver.rules()); err != nil {
		facades.Log.Errorf("add rules error: %+v", err)
	}
}

func (receiver *ValidationServiceProvider) rules() []validation.Rule {
	return []validation.Rule{
		&rules.Exists{},
		&rules.NotExists{},
		&rules.Captcha{},
		&rules.Phone{},
		&rules.VerifyCode{},
		&rules.ReCaptcha{},
	}
}
