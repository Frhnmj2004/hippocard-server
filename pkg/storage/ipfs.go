package storage

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Frhnmj2004/hippocard-server/configs"

	ipfsapi "github.com/ipfs/go-ipfs-api"
)

// IPFSClient manages interactions with IPFS via Pinata
type IPFSClient struct {
	Shell    *ipfsapi.Shell
	MaxTries int
	Timeout  time.Duration
}

// NewIPFSClient initializes a new IPFS client using Pinata with authentication
func NewIPFSClient(config *configs.Config) (*IPFSClient, error) {
	// Create a custom HTTP client with Pinata authentication
	customClient := &http.Client{
		Transport: &roundTripperWithAuth{
			transport: http.DefaultTransport.(*http.Transport), // Type assertion to *http.Transport
			apiKey:    config.IPFS.APIKey,
			secret:    config.IPFS.Secret,
		},
		Timeout: 30 * time.Second,
	}

	// Initialize Shell with the custom client
	shell := ipfsapi.NewShellWithClient("https://api.pinata.cloud/psa", customClient)

	// Verify connection
	for attempt := 1; attempt <= 3; attempt++ {
		if _, err := shell.ID(); err == nil {
			break
		} else if attempt == 3 {
			log.Printf("Failed to connect to IPFS (Pinata) after %d attempts: %v", attempt, err)
			return nil, err
		}
		time.Sleep(time.Duration(attempt) * time.Second)
	}

	return &IPFSClient{
		Shell:    shell,
		MaxTries: 3,
		Timeout:  30 * time.Second,
	}, nil
}

// AddData uploads data to IPFS and returns the CID
func (c *IPFSClient) AddData(data []byte) (string, error) {
	reader := bytes.NewReader(data)
	for attempt := 1; attempt <= c.MaxTries; attempt++ {
		cid, err := c.Shell.Add(reader)
		if err == nil {
			log.Printf("Successfully added data to IPFS, CID: %s", cid)
			return cid, nil
		}
		log.Printf("Attempt %d failed to add data to IPFS: %v", attempt, err)
		if attempt < c.MaxTries {
			time.Sleep(time.Duration(attempt) * time.Second)
			reader.Reset(data) // Reset reader for retry
		}
	}
	return "", logError("Failed to add data to IPFS after %d attempts", c.MaxTries)
}

func (c *IPFSClient) AddBatch(dataChunks [][]byte) ([]string, error) {
	var cids []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, chunk := range dataChunks {
		wg.Add(1)
		go func(data []byte) {
			defer wg.Done()
			cid, err := c.AddData(data)
			if err != nil {
				log.Printf("Failed to add batch chunk to IPFS: %v", err)
				return
			}
			mu.Lock()
			cids = append(cids, cid)
			mu.Unlock()
		}(chunk)
	}

	wg.Wait()
	if len(cids) != len(dataChunks) {
		return nil, logError("Failed to add batch to IPFS: incomplete uploads")
	}
	return cids, nil
}

func (c *IPFSClient) GetData(cid string) ([]byte, error) {
	for attempt := 1; attempt <= c.MaxTries; attempt++ {
		reader, err := c.Shell.Cat(cid)
		if err == nil {
			defer reader.Close()
			data, err := io.ReadAll(reader)
			if err != nil {
				log.Printf("Attempt %d failed to read IPFS data for CID %s: %v", attempt, cid, err)
				continue
			}
			log.Printf("Successfully retrieved data from IPFS, CID: %s", cid)
			return data, nil
		}
		log.Printf("Attempt %d failed to retrieve data from IPFS for CID %s: %v", attempt, cid, err)
		if attempt < c.MaxTries {
			time.Sleep(time.Duration(attempt) * time.Second)
		}
	}
	return nil, logError("Failed to retrieve data from IPFS for CID %s after %d attempts", cid, c.MaxTries)
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

func logError(format string, args ...interface{}) error {
	errMsg := fmt.Sprintf(format, args...)
	err := fmt.Errorf(errMsg)
	log.Println(err)
	return err
}
