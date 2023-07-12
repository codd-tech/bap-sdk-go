package middleware

import (
	"context"
	"log"

	"github.com/codd-tech/bap-sdk-go"
	tele "gopkg.in/telebot.v3"
)

func TelebotBapMiddleware(token string) tele.MiddlewareFunc {
	bap, err := bap.NewBAPClient(token)
	if err != nil {
		log.Fatal(err)
	}

	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			update, err := bap.HandleUpdate(context.Background(), c.Update())
			if err != nil {
				return err
			}
			if update {
				return next(c)
			} else {
				return nil
			}
		}
	}
}
