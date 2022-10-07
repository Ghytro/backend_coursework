package common

import "log"

func LogFatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
