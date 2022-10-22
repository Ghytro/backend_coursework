package profile

import (
	"backend_coursework/internal/common"
	"backend_coursework/internal/view"
	"html/template"
	"sync"
)

var templates *common.SyncMap[string, *template.Template]

func init() {
	templatesMap, err := view.GenTemplatesMap(
		"profile/any.html",
		"profile/my.html",
	)
	common.LogFatalErr(err)
	templates = common.NewSyncMap[string, *template.Template](&sync.Mutex{})
	for k, v := range templatesMap {
		templates.Set(k, v)
	}
}
