package main

import (
	"flag"
)

func cmdLine() (string, int, bool, bool) { //apiKey, time to sleep, is headless or not, debug

	keyPtr := flag.String("apiKey", "", "tg apiKey for the bot")

	timeToSleep := flag.Int("sleep", 0, "time to sleep before screenshoting page")

	headlessPtr := flag.Bool("headless", true, "is headless")

	debugPtr := flag.Bool("debug", false, "is debug mode")

	flag.Parse()

	return *keyPtr, *timeToSleep, *headlessPtr, *debugPtr

}
