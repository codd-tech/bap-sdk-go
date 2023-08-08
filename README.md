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
import "github.com/codd-tech/bap-sdk-go"
```

### Creating a BAP Client

To create a new instance of the BAP client, use the `NewBAPClient` function:

```golang
// Setup BAP
bapClient, err := bap.NewBAPClient("your-api-key")
if err != nil {
    log.Fatal(err)
}
defer bapClient.Close()
```

**The `APIKey` field is mandatory and represents your Advertising Platform API key.** <br>

### Handling Updates

Use the `HandleUpdate` method of the BAP client to send update data to the BAP API:

```golang
needHandle, err := bapClient.HandleUpdate(context.Background(), update)
if err != nil {
	log.Print(err)
}
```

If your advertisement mode is set to manual you can mark ad placement in your code by calling:

```golang
err = bapClient.SendAdvertisement(context.Background(), update)
if err != nil {
	log.Print(err)
}
```

#### Interrupting control flow

At times, BAP may introduce telegram updates within its advertisement flow. To maintain the logical consistency of your bot, it is necessary to ignore such updates.

The `bap.HandleUpdate` method returns a boolean value indicating whether you should proceed with handling the request or skip it as an internal BAP request.

When the method returns `false`, it signifies that the current request should not be processed by your bot.

### Handling Updates Using Middleware (Example on Telebot)

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

## About

### Submitting bugs and feature requests

Bugs and feature request are tracked on [GitHub](https://github.com/codd-tech/bap-sdk-go)

### License

Bot Advertising Platform SDK is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
