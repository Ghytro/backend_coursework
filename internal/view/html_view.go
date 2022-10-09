package view

import (
	"backend_coursework/internal/common"
	"html/template"
)

func GenTemplatesMap(paths ...string) (map[string]*template.Template, error) {
	m := make(map[string]*template.Template)
	for _, p := range paths {
		t, err := template.ParseFiles("./web/" + p)
		common.LogFatalErr(err)
		m[p] = t
	}
	return m, nil
}
