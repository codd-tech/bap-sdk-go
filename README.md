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

Import the `bap` package into your Go code:

```golang
import "github.com/your/package/import/path/bap"
```

### Creating a BAP Client

To create a new instance of the BAP client, use the `NewBotAdvertisingPlatform` function:

```golang
config := bap.BAPConfig{
	APIKey: "your-api-key",
	Addr:   "api.production.bap.codd.io:8080", // optional, defaults to "api.production.bap.codd.io:8080"
}

client, err := bap.NewBotAdvertisingPlatform(config)
if err != nil {
	// handle error
}
```

**The `APIKey` field is mandatory and represents your Advertising Platform API key.** <br>
The `Addr` field is optional and specifies the address of the BAP API. If not provided, it defaults to "api.production.bap.codd.io:8080".

### Handling Updates (Example on Telegram Bot API)

[telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)

Use the `HandleUpdate` method of the BAP client to send update data to the BAP API:

```golang
// -- some code for setup telegram bot --
bot, err := tgBotApi.NewBotAPI(TELEGRAM_TOKEN)
if err != nil {
    log.Panic(err)
}

bot.Debug = true

u := tgBotApi.NewUpdate(0)
u.Timeout = 60
// -- end of setup bot code --

updates := bot.GetUpdatesChan(u)

for update := range updates {
    // Handle update
    err = bap.HandleUpdate(context.Background(), update)
    if err != nil {
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
