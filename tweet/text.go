package tweet

import (
	"strconv"
	"strings"
	"syscall"

	"github.com/lxn/win"
)

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

func errorMessageBox(message string) {
	win.MessageBox(win.HWND(0), UTF16PtrFromString(message), UTF16PtrFromString("Error"), win.MB_OK+win.MB_ICONEXCLAMATION)
}

// "Go から Windows の MessageBox を呼び出す - Qiita" https://qiita.com/manymanyuni/items/867d7e0112ce22dec6d5
func UTF16PtrFromString(s string) *uint16 {
	result, _ := syscall.UTF16PtrFromString(s)
	return result
}
