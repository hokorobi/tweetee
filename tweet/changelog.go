package tweet

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"

	goutils "github.com/hokorobi/go-utils"
)

type configChangelog struct {
	Path string `json:"path"`
}

func TweetChangelog(text string) error {
	d, err := loadConfigChangelog()
	if err != nil {
		return err
	}

	f, err := os.Open(d.Path)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	buf, err := buildChangelog(f, text)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	f.Close()

	// write file
	f, err = os.OpenFile(d.Path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, l := range buf {
		w.WriteString(l + "\r\n")
	}
	w.Flush()
	return nil
}

func getChangelogToday() string {
	weekdayja := strings.NewReplacer(
		"Sun", "日",
		"Mon", "月",
		"Tue", "火",
		"Wed", "水",
		"Thu", "木",
		"Fri", "金",
		"Sat", "土",
	)
	return weekdayja.Replace(time.Now().Format("2006-01-02 (Mon)"))
}

func matchDate(t string) bool {
	r := regexp.MustCompile(`^20[0-9][0-9]-[01][0-9]-[0-3][0-9]`)
	return r.MatchString(t)
}

func getTweet(text string) []string {
	text = time.Now().Format("15:04") + " " + text

	var tweetLines []string
	for _, b := range strings.Split(text, "\n") {
		tweetLines = append(tweetLines, "\t"+b)
	}
	tweetLines = append(tweetLines, "")

	return tweetLines
}

func buildChangelog(reader io.Reader, text string) ([]string, error) {
	category := "	* Diary: tweet"
	today := getChangelogToday()
	appendtweet := []string{today, "", category, ""}
	appendtweet = slices.Concat(appendtweet, getTweet(text))
	var line int
	var buf []string
	scanner := bufio.NewScanner(reader)

	// find today entry
	for scanner.Scan() {
		buf = append(buf, scanner.Text())
		// goutils.PrintTee(scanner.Text())
		// goutils.PrintTee(strings.Index(scanner.Text(), today))
		if strings.Index(scanner.Text(), today) == 0 {
			matchline := line
			// goutils.PrintTee("match today")
			// goutils.PrintTee(line)

			// find category
			for scanner.Scan() {
				buf = append(buf, scanner.Text())
				// goutils.PrintTee(line)
				// goutils.PrintTee(scanner.Text())
				if strings.Index(scanner.Text(), category) == 0 {
					// goutils.PrintTee("match category")

					// skip blank line
					scanner.Scan()
					buf = slices.Concat(buf, appendtweet[3:len(appendtweet)-1])
					break
				}

				// reached next day or read 100 lines
				if matchDate(scanner.Text()) || line > 100 {
					buf = slices.Concat(buf[0:matchline+2], appendtweet[2:], buf[matchline+2:])
					break
				}
				line++
			}
			break
		}
		line++

		if line > 100 {
			buf = slices.Concat(appendtweet, buf)
			break
		}
	}

	for scanner.Scan() {
		buf = append(buf, scanner.Text())
	}

	return buf, nil
}

func loadConfigChangelog() (configChangelog, error) {
	var d configChangelog

	f, err := os.Open(goutils.GetFilenameSameBase(".json"))
	if err != nil {
		return d, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(&d)
	if err != nil {
		return d, err
	}

	return d, nil
}
