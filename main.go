package main

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

func main() {

	apikey, timeToSleep, isHeadless := cmdLine()

	onceClickedCookieButton := false

	opts := optsHeadOrNot(isHeadless)

	allocCtx, cancel := chromedp.NewExecAllocator(
		context.Background(), opts...,
	)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	bot, err := tgbotapi.NewBotAPI(apikey)

	if err != nil {
		panic(err)
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	wg := &sync.WaitGroup{}

	for update := range updates {

		wg.Add(1)

		switch checkIfTwitter(update.Message.Text) {
		case true:
			go func(update tgbotapi.Update) {

				sendCroppedPicture(&onceClickedCookieButton, ctx, update, bot, timeToSleep)
				wg.Done()

			}(update)
		case false:
			fmt.Println(update.Message.Text + "not a twitter link")
		}
	}
}

func optsHeadOrNot(isHeadless bool) []func(*chromedp.ExecAllocator) {
	if isHeadless {
		return []chromedp.ExecAllocatorOption{
			chromedp.NoFirstRun,
			chromedp.NoDefaultBrowserCheck,
			chromedp.Headless,
			chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36"),
		}
	} else {
		return []chromedp.ExecAllocatorOption{
			chromedp.NoFirstRun,
			chromedp.NoDefaultBrowserCheck,
			chromedp.Flag("headless", false),
			chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36"),
		}
	}
}

func checkIfTwitter(url string) bool {

	reg, err := regexp.MatchString("[http.]*://twitter.com/*", url)
	check(err)

	return reg
}

func sendCroppedPicture(once *bool, ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, timeToSleep int) {

	xpathCookie := `/html/body/div[1]/div/div/div[1]/div[1]/div/div/div/div/div[2]/div[1]`
	xpathScreenshot := `//*[@id="react-root"]/div/div/div[2]/main/div/div/div/div/div/section/div/div/div/div/div[1]/div/div/article`

	var buf []byte
	if !*once {
		err := chromedp.Run(ctx, emulation.SetDeviceMetricsOverride(400, 800, 1.0, true), chromedp.Navigate(update.Message.Text), chromedp.Sleep(time.Second*time.Duration(timeToSleep)), clickAccept(xpathCookie), elementScreenshot(xpathScreenshot, &buf))
		*once = true
		check(err)
	} else {
		err := chromedp.Run(ctx, emulation.SetDeviceMetricsOverride(400, 800, 1.0, true), chromedp.Navigate(update.Message.Text), chromedp.Sleep(time.Second*time.Duration(timeToSleep)), elementScreenshot(xpathScreenshot, &buf))
		check(err)
	}

	if checkIfMedia(ctx, update.Message.Text) {

		pict := cutPicture(buf, 400, 700)
		photoFileBytes := tgbotapi.FileBytes{
			Name:  "picture",
			Bytes: pict,
		}

		bot.Send(tgbotapi.NewPhoto(update.Message.Chat.ID, photoFileBytes))
	} else {

		photoFileBytes := tgbotapi.FileBytes{
			Name:  "picture",
			Bytes: buf,
		}

		bot.Send(tgbotapi.NewPhoto(update.Message.Chat.ID, photoFileBytes))
	}
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

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
