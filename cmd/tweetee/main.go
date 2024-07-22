package main

import (
	"os"

	goutils "github.com/hokorobi/go-utils"
)

func main() {
	if len(os.Args) < 2 {
		tweet.errorMessageBox("Needs tweet text.")
		return
	}

	text := tweet.genText(os.Args[1:])

	var warn bool
	err := tweet.tweetChangelog(text)
	if err != nil {
		goutils.PrintTee(err)
		warn = true
	}

	err = tweet.postBsky(text)
	if err != nil {
		goutils.PrintTee(err)
		warn = true
	}

	if warn {
		tweet.errorMessageBox("Check log file.")
	}
}
