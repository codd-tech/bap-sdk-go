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
	// handle error
    log.Panic(err)
}
defer bap.Close()
```

**The `APIKey` field is mandatory and represents your Advertising Platform API key.** <br>

### Handling Updates (Example on Telegram Bot API)

[telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)

Use the `HandleUpdate` method of the BAP client to send update data to the BAP API:

```golang
// -- some code for setup telegram bot --
bot, err := tgBotApi.NewBotAPI(TELEGRAM_TOKEN)
if err != nil {
	// handle error
    panic(err)
}

u := tgBotApi.NewUpdate(0)
u.Timeout = 60
// -- end of setup bot code --

updates := bot.GetUpdatesChan(u)

for update := range updates {
    // Handle update
    err = bap.HandleUpdate(context.Background(), update)
    if err != nil {
	    // handle error
        log.Print(err)
    }

    // any telegram bot logic
}
```

## About

### Submitting bugs and feature requests

Bugs and feature request are tracked on [GitHub](https://github.com/codd-tech/bap-sdk-go)

### License

Bot Advertising Platform SDK is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
