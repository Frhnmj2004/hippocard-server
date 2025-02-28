package storage

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/Frhnmj2004/hippocard-server/configs"

	ipfsapi "github.com/ipfs/go-ipfs-api"
)

// IPFSClient manages interactions with IPFS via Pinata
type IPFSClient struct {
	Shell *ipfsapi.Shell
}

// NewIPFSClient initializes a new IPFS client using Pinata with authentication
func NewIPFSClient(config *configs.IPFSConfig) (*IPFSClient, error) {
	// Create a custom HTTP client with Pinata authentication
	customClient := &http.Client{
		Transport: &roundTripperWithAuth{
			transport: http.DefaultTransport.(*http.Transport), // Type assertion to *http.Transport
			apiKey:    config.APIKey,
			secret:    config.Secret,
		},
	}

	// Initialize Shell with the custom client
	shell := ipfsapi.NewShellWithClient("https://api.pinata.cloud/psa", customClient)

	// Verify connection
	if _, err := shell.ID(); err != nil {
		log.Printf("Failed to connect to IPFS (Pinata): %v", err)
		return nil, err
	}

	return &IPFSClient{Shell: shell}, nil
}

// AddData uploads data to IPFS and returns the CID
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

// roundTripperWithAuth adds Pinata authentication headers to requests
type roundTripperWithAuth struct {
	transport *http.Transport
	apiKey    string
	secret    string
}

func (r roundTripperWithAuth) RoundTrip(req *http.Request) (*http.Response, error) {
	// Add Pinata authentication headers
	req.Header.Set("pinata_api_key", r.apiKey)
	req.Header.Set("pinata_secret_api_key", r.secret)
	return r.transport.RoundTrip(req)
}
