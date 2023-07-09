package main

import (
	"context"
	"regexp"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func checkIfMedia(ctx context.Context, urlstr string) bool {
	xpathMedia := `/html/body/div[1]/div/div/div[2]/main/div/div/div/div/div/section/div/div/div/div/div[1]/div/div/article/div/div/div[3]/div[2]/div/div/div/div/div/div/div/div/div/div[2]/div/div/div/div[2]/div/div[2]/div/div/div[1]`

	var res []*cdp.Node

	chromedp.Run(ctx,
		chromedp.Navigate(urlstr),
		chromedp.Nodes(xpathMedia, &res, chromedp.AtLeast(0)),
	)

	if len(res) == 0 {
		return false
	} else {
		return true
	}
}

func checkIfTwitter(message string) bool {

	reg, err := regexp.MatchString("[http.://]*twitter.com/*", message)
	check(err)

	return reg
}

func checkIfBotIsBlocked(update tgbotapi.Update) bool {
	if update.Message == nil {
		return true
	} else {
		return false
	}
}
