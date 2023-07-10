package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendCroppedPicture(once *bool, ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, timeToSleep int) {

	regexp, _ := regexp.Compile(`(twitter\.com\/([\d\w]*\/*)*)`)
	links := regexp.FindAllStringSubmatch(update.Message.Text, -1)

	lnksCaption := ""

	for _, el := range links {
		lnksCaption += fmt.Sprintf("Ссылка на тви: %s \n", el[1])
	}

	log.Println(links)

	var buf []byte

	var sliceOfPhotos []interface{}

	xpathCookie := `/html/body/div[1]/div/div/div[1]/div[1]/div/div/div/div/div[2]/div[1]`
	xpathScreenshot := `//*[@id="react-root"]/div/div/div[2]/main/div/div/div/div/div/section/div/div/div/div/div[1]/div/div/article`

	for idx, el := range links {

		url := el[1]

		log.Println(url)

		if !*once {
			err := chromedp.Run(ctx, emulation.SetDeviceMetricsOverride(400, 800, 1.0, true), chromedp.Navigate("https://"+url), chromedp.Sleep(time.Second*time.Duration(timeToSleep)), clickAccept(xpathCookie), elementScreenshot(xpathScreenshot, &buf))
			*once = true
			check(err)
		} else {
			err := chromedp.Run(ctx, emulation.SetDeviceMetricsOverride(400, 800, 1.0, true), chromedp.Navigate("https://"+url), chromedp.Sleep(time.Second*time.Duration(timeToSleep)), elementScreenshot(xpathScreenshot, &buf))
			check(err)
		}

		if checkIfMedia(ctx, update.Message.Text) {

			pict := cutPicture(buf, 400, 700)
			photoFileBytes := tgbotapi.FileBytes{
				Name:  "picture",
				Bytes: pict,
			}

			inputPhoto := tgbotapi.NewInputMediaPhoto(photoFileBytes)

			if idx == 0 {
				inputPhoto.Caption = lnksCaption
			}

			sliceOfPhotos = append(sliceOfPhotos, inputPhoto)
		} else {

			photoFileBytes := tgbotapi.FileBytes{
				Name:  "picture",
				Bytes: buf,
			}

			inputPhoto := tgbotapi.NewInputMediaPhoto(photoFileBytes)

			if idx == 0 {
				inputPhoto.Caption = lnksCaption
			}

			sliceOfPhotos = append(sliceOfPhotos, inputPhoto)
		}
	}

	mediaGroup := tgbotapi.NewMediaGroup(update.Message.Chat.ID, sliceOfPhotos)
	mediaGroup.ReplyToMessageID = update.Message.MessageID
	bot.Send(mediaGroup)
}
