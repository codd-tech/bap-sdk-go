package bap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	defaultAddr = "api.production.bap.codd.io:8080"
	ApiVersion  = 2
	BapPrefix   = "/__bap"
)

// BAPClient is an interface that defines the methods for a BAP client.
type BAPClient interface {
	// HandleUpdate sends the update data to the BAP API.
	//
	// Return false, if you do not need to handle this update (because it is a bap command).
	HandleUpdate(ctx context.Context, update interface{}) (bool, error)
	// Close closes the BAP UDP connection.
	Close() error
}

type bap struct {
	apiKey string
	socket *net.UDPConn
}

type updateStub struct {
	Callback *callback `json:"callback_query,omitempty"`
}

type callback struct {
	Data string `json:"data"`
}

// Creates a new instance of the BAP client.
func NewBAPClient(apiKey string) (BAPClient, error) {
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

// HandleUpdate handles the update received from the BAP API.
//
// Return false, if you do not need to handle this update (because it is a bap command).
func (b *bap) HandleUpdate(ctx context.Context, update interface{}) (bool, error) {
	data := map[string]interface{}{
		"api_key": b.apiKey,
		"version": ApiVersion,
		"update":  update,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return true, fmt.Errorf("failed to marshal JSON data: %w", err)
	}

	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := b.socket.Write(jsonData)
			if err != nil {
				log.Printf("Failed to write UDP data: %v\n", err)
			}
			log.Printf("Wrote %d bytes\n", n)
		}
	}()

	// Serialize the update.
	updateBytes, err := json.Marshal(update)
	if err != nil {
		return true, fmt.Errorf("failed to marshal JSON data: %w", err)
	}

	updateObj := &updateStub{}

	// Deserialize the update.
	err = json.Unmarshal(updateBytes, updateObj)
	if err != nil {
		return true, fmt.Errorf("failed to marshal JSON data: %w", err)
	}

	// Check if the update is a bap command.
	// if yes, do not need to handle this update
	if updateObj.Callback != nil && strings.HasPrefix(updateObj.Callback.Data, BapPrefix) {
		return false, nil
	}

	return true, nil
}

// Close closes the BAP UDP connection.
func (b *bap) Close() error {
	return b.socket.Close()
}
