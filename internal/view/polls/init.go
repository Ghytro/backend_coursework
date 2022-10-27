package polls

import (
	"backend_coursework/internal/common"
	"backend_coursework/internal/view"
	"html/template"
)

var templates *common.SyncMap[string, *template.Template]

func init() {
	var err error
	templates, err = view.GenTemplatesMap(
		"polls/new.html",
	)
	common.LogFatalErr(err)
}
