package bap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	defaultAddr             = "api.production.bap.codd.io:8080"
	APIVersion              = 3
	CallbackQueryDataPrefix = "/__bap"
)

// Client is an interface that defines the methods for a BAP client.
type Client interface {
	// HandleUpdate sends the update data to the BAP API.
	//
	// Return false, if you do not need to handle this update (because it is a bap command).
	HandleUpdate(ctx context.Context, update interface{}) (bool, error)
	// SendAdvertisement send advertisement to the BAP API, allows you to immediately display ads
	SendAdvertisement(ctx context.Context, update interface{}) error
	// Close closes the BAP UDP connection.
	Close() error
}

type bap struct {
	apiKey string
	socket *net.UDPConn
}

// Validate the update received Update.
func Validate(update interface{}) (*tgbotapi.Update, error) {
	// Convert update to JSON bytes
	updateJSON, err := json.Marshal(update)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update to JSON: %w", err)
	}

	// Unmarshal JSON bytes back to Update struct
	var updatedUpdate *tgbotapi.Update
	err = json.Unmarshal(updateJSON, &updatedUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to Update struct: %w", err)
	}

	return updatedUpdate, nil
}

// NewBAPClient Creates a new instance of the BAP client.
func NewBAPClient(apiKey string) (Client, error) {
	if apiKey == "" {
		return nil, errors.New("AdvertisingPlatform API key is empty")
	}

	addr, err := net.ResolveUDPAddr("udp", defaultAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	// Create a UDP connection
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial UDP: %w", err)
	}

	return &bap{apiKey: apiKey, socket: conn}, nil
}

func (b *bap) sendToBAP(update *tgbotapi.Update, method string) error {
	data := map[string]interface{}{
		"api_key": b.apiKey,
		"version": APIVersion,
		"update":  update,
		"method":  method,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data: %w", err)
	}

	_, err = b.socket.Write(jsonData)
	if err != nil {
		return fmt.Errorf("Failed to write UDP data: %v\n", err)
	}

	return nil
}

// HandleUpdate handles the update received from the BAP API.
//
// Return false, if you do not need to handle this update (because it is a bap command).
func (b *bap) HandleUpdate(ctx context.Context, update interface{}) (bool, error) {
	updateObj, err := Validate(update)
	if err != nil {
		return true, fmt.Errorf("failed to validate update: %w", err)
	}

	if err := b.sendToBAP(updateObj, "activity"); err != nil {
		return true, fmt.Errorf("failed to send update to BAP: %w", err)
	}

	return !b.isBAPUpdate(updateObj), nil
}

// SendAdvertisement send advertisement to the BAP API, allows you to immediately display ads
func (b *bap) SendAdvertisement(ctx context.Context, update interface{}) error {
	updateObj, err := Validate(update)
	if err != nil {
		return fmt.Errorf("failed to validate update: %w", err)
	}

	if b.isBAPUpdate(updateObj) {
		return nil
	}

	if err := b.sendToBAP(updateObj, "advertisement"); err != nil {
		return fmt.Errorf("failed to send update to BAP: %w", err)
	}

	return nil
}

func (b *bap) isBAPUpdate(update *tgbotapi.Update) bool {
	return update.CallbackQuery != nil && strings.HasPrefix(update.CallbackQuery.Data, CallbackQueryDataPrefix)
}

// Close closes the BAP UDP connection.
func (b *bap) Close() error {
	return b.socket.Close()
}
