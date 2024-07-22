package main

import (
	"os"

	goutils "github.com/hokorobi/go-utils"
	tw "github.com/hokorobi/tweetee/tweet"
)

func main() {
	if len(os.Args) < 2 {
		tw.ErrorMessageBox("Needs tweet text.")
		return
	}

	text := tw.GenText(os.Args[1:])

	err := tw.PostBsky(text)
	if err != nil {
		goutils.PrintTee(err)
		tw.ErrorMessageBox("Check log file.")
	}
}
