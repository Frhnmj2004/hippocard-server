// IPFS client for medical history
package storage

import (
	"bytes"
	"context"
	"io"
	"log"

	"github.com/Frhnmj2004/hippocard-server/configs"

	ipfsapi "github.com/ipfs/go-ipfs-api"
)

// IPFSClient manages interactions with IPFS
type IPFSClient struct {
	Shell *ipfsapi.Shell
}

// NewIPFSClient initializes a new IPFS client using Pinata or a custom IPFS node
func NewIPFSClient(config *configs.Config) (*IPFSClient, error) {
	// Use Pinata's gateway with API key and secret
	// Pinata endpoint: https://api.pinata.cloud/psa
	shell := ipfsapi.NewShell("https://api.pinata.cloud/psa")

	// Set Pinata authentication headers
	shell.SetHeader("pinata_api_key", config.IPFS.APIKey)
	shell.SetHeader("pinata_secret_api_key", config.IPFS.Secret)

	// Verify connection
	ctx := context.Background()
	if _, err := shell.ID(ctx); err != nil {
		log.Printf("Failed to connect to IPFS (Pinata): %v", err)
		return nil, err
	}

	return &IPFSClient{Shell: shell}, nil
}

// AddData uploads data to IPFS and returns CID
func (c *IPFSClient) AddData(data []byte) (string, error) {
	reader := bytes.NewReader(data)
	cid, err := c.Shell.Add(reader)
	if err != nil {
		log.Printf("Failed to add data to IPFS: %v", err)
		return "", err
	}
	return cid, nil
}

// GetData retrieves data from IPFS by CID
func (c *IPFSClient) GetData(cid string) ([]byte, error) {
	reader, err := c.Shell.Cat(cid)
	if err != nil {
		log.Printf("Failed to retrieve data from IPFS: %v", err)
		return nil, err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("Failed to read IPFS data: %v", err)
		return nil, err
	}
	return data, nil
}
