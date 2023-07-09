package main

import (
	"flag"
)

func cmdLine() (string, int, bool) { //apiKey, time to sleep, is headless or not

	keyPtr := flag.String("apiKey", "", "tg apiKey for the bot")

	timeToSleep := flag.Int("sleep", 0, "time to sleep before screenshoting page")

	headlessPtr := flag.Bool("headless", true, "is headless")

	flag.Parse()

	return *keyPtr, *timeToSleep, *headlessPtr

}
