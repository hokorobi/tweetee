package main

import (
	"os"
	"strconv"
	"strings"

	goutils "github.com/hokorobi/go-utils"
)

func main() {
	if len(os.Args) < 2 {
		errorMessageBox("Needs tweet text.")
		return
	}

	text := genText(os.Args[1:])

	var warn bool
	err := tweetChangelog(text)
	if err != nil {
		goutils.PrintTee(err)
		warn = true
	}

	err = postBsky(text)
	if err != nil {
		goutils.PrintTee(err)
		warn = true
	}

	if warn {
		errorMessageBox("Check log file.")
	}
}

func genText(args []string) string {
	var text string
	for _, a := range args {
		// スペースがあったら "" でくくる
		if strings.Index(a, " ") >= 0 {
			a = strconv.Quote(a)
		}
		text = text + " " + a
	}
	return text
}
