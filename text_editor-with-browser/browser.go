package main

import (
	"os"
	"runtime"

	webview "github.com/webview/webview_go"
)

func main() {
	url := "https://google.com"
	if len(os.Args) > 1 {
		url = os.Args[1]
	}
	runtime.LockOSThread()
	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("Browser")
	w.SetSize(800, 600, webview.Hint(webview.HintNone))
	w.Navigate(url)
	w.Run()
}
