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

	var warn bool
	err := tw.TweetChangelog(text)
	if err != nil {
		goutils.PrintTee(err)
		warn = true
	}

	err = tw.PostBsky(text)
	if err != nil {
		goutils.PrintTee(err)
		warn = true
	}

	err = tw.PostNostr(text)
	if err != nil {
		goutils.PrintTee(err)
		warn = true
	}

	if warn {
		tw.ErrorMessageBox("Check log file.")
	}
}
