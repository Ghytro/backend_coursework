package view

import (
	"fmt"
	"os"
	"strings"
)

func GenViewHTML(view string) (string, error) {
	htmlBytes, err := os.ReadFile(fmt.Sprintf("./internal/view/%s/%s.html", view, view))
	if err != nil {
		return "", err
	}
	cssBytes, err := os.ReadFile(fmt.Sprintf("./internal/view/%s/%s.css", view, view))
	if err != nil {
		return "", err
	}
	htmlString, cssString := string(htmlBytes), string(cssBytes)
	headEndIdx := strings.Index(htmlString, "</head>")
	return htmlString[:headEndIdx] + "<style>\n" + cssString + "\n</style>\n" + htmlString[headEndIdx+len("</head>"):], nil
}
