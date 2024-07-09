package rules

import (
	"github.com/goravel-kit/geetest"
	geetestfacades "github.com/goravel-kit/geetest/facades"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/support/debug"
	"github.com/goravel/framework/support/json"
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
	debug.Dump(val)
	if t, ok := val.(geetest.Ticket); ok {
		verify, _ := geetestfacades.Geetest().Verify(t)
		return verify
	}

	ticket := make(map[string]string)
	if err := json.UnmarshalString(cast.ToString(val), &ticket); err != nil {
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
