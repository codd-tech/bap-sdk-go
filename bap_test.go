package bap_test

import (
	"fmt"
	"testing"

	"github.com/codd-tech/bap-sdk-go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tele "gopkg.in/telebot.v3"
)

func TestValidateWithValidUpdate(t *testing.T) {
	cases := []struct {
		name   string
		update interface{}
	}{
		{
			name: "TGBotAPI",
			update: &tgbotapi.Update{
				UpdateID: 123,
				Message: &tgbotapi.Message{
					MessageID: 456,
					Text:      "Hello, world!",
				},
			},
		},
		{
			name: "TeleBot",
			update: &tele.Update{
				ID: 123,
				Message: &tele.Message{
					ID:   456,
					Text: "Hello, world!",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%v", c.name), func(t *testing.T) {
			_, err := bap.Validate(c.update)
			if err != nil {
				t.Errorf("unexpected error during validation: %v", err)
			}
		})
	}
}

func TestValidateWithInvalidUpdate(t *testing.T) {
	update := "jajaja"

	// Call the function under test
	_, err := bap.Validate(update)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
