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

// BAPClient is an interface that defines the methods for a BAP client.
type BAPClient interface {
	// HandleUpdate sends the update data to the BAP API.
	HandleUpdate(ctx context.Context, update interface{}) error
	// Close closes the BAP UDP connection.
	Close() error
}

type bap struct {
	apiKey string
	socket *net.UDPConn
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
			n, err := b.socket.Write(jsonData)
			if err != nil {
				log.Printf("Failed to write UDP data: %v\n", err)
			}
			log.Printf("Wrote %d bytes\n", n)
		}
	}()

	return nil
}

// Close closes the BAP UDP connection.
func (b *bap) Close() error {
	return b.socket.Close()
}
