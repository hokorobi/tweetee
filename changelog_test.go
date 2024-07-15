package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestMatchDate(t *testing.T) {
	cases := []struct {
		date   string
		answer bool
	}{
		{"2023-01-01", true},
		{"\t2023-01-01", false},
		{"2023-1-01", false},
		{"2023-1-1", false},
		{"2023-1-01", false},
	}

	for _, c := range cases {
		if matchDate(c.date) != c.answer {
			t.Errorf("date: %s, result: %t, answer: %t", c.date, matchDate(c.date), c.answer)
		}
	}
}

func TestGetTweet(t *testing.T) {
	n := time.Now().Format("15:04")
	cases := []struct {
		args   []string
		answer []string
	}{
		{[]string{"test"}, []string{fmt.Sprintf("\t%s test", n), ""}},
		{[]string{"test test"}, []string{fmt.Sprintf("\t%s \"test test\"", n), ""}},
		{[]string{"test", "test"}, []string{fmt.Sprintf("\t%s test test", n), ""}},
		{[]string{"a b", "c", "d"}, []string{fmt.Sprintf("\t%s \"a b\" c d", n), ""}},
		{[]string{"a\nb"}, []string{fmt.Sprintf("\t%s a", n), "\tb", ""}},
		// {[]string{"a \n b"}, []string{fmt.Sprintf("\t%s a", n), "\tb", ""}},                               // 実際の動きと違う？
		// {[]string{"test test\ntest test"}, []string{fmt.Sprintf("\t%s test test", n), "\ttest test", ""}}, // 実際の動きと違う？
	}

	for _, c := range cases {
		if !reflect.DeepEqual(getTweet(genText(c.args)), c.answer) {
			t.Errorf("date: %s, result: %s, answer: %s", c.args, getTweet(genText(c.args)), c.answer)
		}
	}
}
