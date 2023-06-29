package bap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
)

const defaultAddr = "api.production.bap.codd.io:8080"

// BAPConfig represents the configuration for the BAP client.
type BAPConfig struct {
	APIKey string // APIKey is the Advertising Platform API key.
	Addr   string // Addr is the address of the BAP API. Defaults to "api.production.bap.codd.io:8080".
}

// BAPClient is an interface that defines the methods for a BAP client.
type BAPClient interface {
	// HandleUpdate sends the update data to the BAP API.
	HandleUpdate(ctx context.Context, update interface{}) error
	// Close closes the BAP UDP client connection.
	Close() error
}

type bap struct {
	apiKey string
	socket *net.UDPConn
}

// Creates a new instance of the BAP client.
func NewBotAdvertisingPlatform(config BAPConfig) (BAPClient, error) {
	if config.APIKey == "" {
		return nil, errors.New("AdvertisingPlatform API key is empty")
	}
	if config.Addr == "" {
		config.Addr = defaultAddr
	}

	addr, err := net.ResolveUDPAddr("udp", config.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	// Create a UDP connection
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial UDP: %w", err)
	}

	return &bap{apiKey: config.APIKey, socket: conn}, nil
}

// HandleUpdate handles the update received from the BAP API.
func (b *bap) HandleUpdate(ctx context.Context, update interface{}) error {
	data := map[string]interface{}{
		"api_key": b.apiKey,
		"update":  update,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data: %w", err)
	}

	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			_, err := b.socket.Write(jsonData)
			if err != nil {
				log.Printf("Failed to write UDP data: %v\n", err)
			}
		}
	}()

	return nil
}

// Close closes the BAP client connection.
func (b *bap) Close() error {
	return b.socket.Close()
}
