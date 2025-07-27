package main

import (
	"chesser/bot"
	"chesser/storage"
	"fmt"
)

func main() {
	storage.LoadTrackedUser()

	err := bot.Start()
	if err != nil {
		panic("❌ Failed to start the bot : " + err.Error())
	}

	fmt.Println("✅ Bot is running ...")

	select {}
}
