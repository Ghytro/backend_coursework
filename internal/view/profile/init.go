package profile

import (
	"backend_coursework/internal/common"
	"backend_coursework/internal/view"
)

var template string

func init() {
	var err error
	template, err = view.GenViewHTML("profile")
	common.LogFatalErr(err)
}
