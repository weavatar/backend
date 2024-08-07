package rules

import (
	"github.com/goravel-kit/geetest"
	geetestfacades "github.com/goravel-kit/geetest/facades"
	"github.com/goravel/framework/contracts/validation"
	"github.com/spf13/cast"
)

type Geetest struct {
}

// Signature The name of the rule.
func (receiver *Geetest) Signature() string {
	return "geetest"
}

// Passes Determine if the validation rule passes.
func (receiver *Geetest) Passes(data validation.Data, val any, options ...any) bool {
	ticket, err := cast.ToStringMapStringE(val)
	if err != nil {
		return false
	}

	keys := []string{"lot_number", "captcha_output", "pass_token", "gen_time"}
	for _, key := range keys {
		if _, ok := ticket[key]; !ok {
			return false
		}
	}

	verify, _ := geetestfacades.Geetest().Verify(geetest.Ticket{
		LotNumber:     ticket["lot_number"],
		CaptchaOutput: ticket["captcha_output"],
		PassToken:     ticket["pass_token"],
		GenTime:       ticket["gen_time"],
	})

	return verify
}

// Message Get the validation error message.
func (receiver *Geetest) Message() string {
	return ""
}
