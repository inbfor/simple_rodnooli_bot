package main

import (
	"fmt"

	"github.com/chromedp/chromedp"
)

func elementScreenshot(sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}

func clickAccept(class string) chromedp.Action {
	fmt.Println(class + "class")
	return chromedp.Tasks{
		chromedp.Click(class, chromedp.BySearch),
	}
}
