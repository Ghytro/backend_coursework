package profile

import (
	"backend_coursework/internal/common"
)

var template string

func init() {
	var err error
	template, err = common.GenViewHTML("profile")
	common.LogFatalErr(err)
}
