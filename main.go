package main

import (
	"context"
	"fmt"
	"regexp"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/chromedp/chromedp"
)

func main() {

	apikey, timeToSleep, isHeadless, debug := cmdLine()

	onceClickedCookieButton := false

	opts := optsHeadOrNot(isHeadless)

	allocCtx, cancel := chromedp.NewExecAllocator(
		context.Background(), opts...,
	)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	bot, err := tgbotapi.NewBotAPI(apikey)

	bot.Debug = debug

	if err != nil {
		panic(err)
	}

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	wg := &sync.WaitGroup{}

	for update := range updates {

		if checkIfBotIsBlocked(update) {
			continue
		}

		wg.Add(1)

		if checkIfTwitter(update.Message.Text) {
			go func(update tgbotapi.Update, timeToSleep int) {

				sendCroppedPicture(&onceClickedCookieButton, ctx, update, bot, timeToSleep)
				wg.Done()

			}(update, timeToSleep)
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

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func findTwitterLinks(message string) [][]string {
	regexp, _ := regexp.Compile(`^([http.://]*twitter.com/*)$`)
	return regexp.FindAllStringSubmatch(message, -1)
}
