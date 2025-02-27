package configs

import (
	"fmt"
	"log"
	"os"
)

// FirebaseConfig holds Firebase-related settings
type FirebaseConfig struct {
	CredentialsPath string
}

// BlockchainConfig holds Polygon blockchain settings
type BlockchainConfig struct {
	RPCURL          string
	ContractAddress string
}

// IPFSConfig holds IPFS settings
type IPFSConfig struct {
	APIKey string
	Secret string
}

// Config aggregates all configurations
type Config struct {
	ServerPort string
	Firebase   FirebaseConfig
	Blockchain BlockchainConfig
	IPFS       IPFSConfig
}

// LoadConfig retrieves environment variables and returns a validated Config struct
func LoadConfig() (*Config, error) {
	// Assume godotenv.Load() is called in main.go, so env vars are already available
	config := &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		Firebase: FirebaseConfig{
			CredentialsPath: getEnv("FIREBASE_CREDENTIALS_PATH", ""),
		},
		Blockchain: BlockchainConfig{
			RPCURL:          getEnv("POLYGON_RPC", "https://rpc-mumbai.maticvigil.com"),
			ContractAddress: getEnv("CONTRACT_ADDRESS", ""),
		},
		IPFS: IPFSConfig{
			APIKey: getEnv("IPFS_API_KEY", ""),
			Secret: getEnv("IPFS_SECRET", ""),
		},
	}

	// Validate required fields
	if config.Firebase.CredentialsPath == "" {
		return nil, logError("FIREBASE_CREDENTIALS_PATH is required")
	}
	if config.Blockchain.ContractAddress == "" {
		return nil, logError("CONTRACT_ADDRESS is required")
	}
	if config.IPFS.APIKey == "" || config.IPFS.Secret == "" {
		return nil, logError("IPFS_API_KEY and IPFS_SECRET are required")
	}

	return config, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// logError logs and returns an error
func logError(msg string) error {
	err := fmt.Errorf(msg)
	log.Println(err)
	return err
}
