# Bot Advertising Platform SDK

This repository holds SDK related to [Bot Advertising Platform](https://publisher.socialjet.pro/).

## Requirements

## Installation

To install the `bap` package, use the following command:

```bash
go get github.com/codd-tech/bap-sdk-go
```

## Usage

### Importing the Package

Import the `bap-sdk-go` package into your Go code:

```golang
import bapSdk "github.com/codd-tech/bap-sdk-go"
```

### Creating a BAP Client

To create a new instance of the BAP client, use the `NewBAPClient` function:

```golang
// Setup BAP
bap, err := bapSdk.NewBAPClient("your-api-key")
if err != nil {
    log.Fatal(err)
}
defer bap.Close()
```

**The `APIKey` field is mandatory and represents your Advertising Platform API key.** <br>

### Handling Updates (Example on Telebot)

[telebot](https://github.com/tucnak/telebot/tree/v3)

Use the `BapMiddleware`

```golang
package main

import (
	"log"
	"time"

	"github.com/codd-tech/bap-sdk-go/middleware"

	tele "gopkg.in/telebot.v3"
)

func main() {
	pref := tele.Settings{
		Token:  "TELEGRAM_TOKEN",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	// Bap middleware
	b.Use(middleware.TelebotBapMiddleware("your-api-key"))

	b.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello Moderator!")
	})

	b.Start()
}
```

### Handling Updates (Example on Telegram Bot API)

[telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)

Use the `HandleUpdate` method of the BAP client to send update data to the BAP API:

```golang
package main

import (
	"context"
	"log"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	bapSdk "github.com/codd-tech/bap-sdk-go"
)

func main() {
	log.Println("Starting bot")

	bot, err := tgBotApi.NewBotAPI("TELEGRAM_TOKEN")
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Setup BAP
	bap, err := bapSdk.NewBAPClient("your-api-key")
	if err != nil {
		log.Fatal(err)
	}
	defer bap.Close()

	u := tgBotApi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// Handle update
		err = bap.HandleUpdate(context.Background(), update)
		if err != nil {
			log.Print(err)
		}

		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgBotApi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

```

## About

### Submitting bugs and feature requests

Bugs and feature request are tracked on [GitHub](https://github.com/codd-tech/bap-sdk-go)

### License

Bot Advertising Platform SDK is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
