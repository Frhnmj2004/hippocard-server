// Polygon client setup
package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/Frhnmj2004/hippocard-server/configs"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// PrescriptionNFTCaller is the interface for read-only contract calls
type PrescriptionNFTCaller interface {
	Mint(opts *bind.CallOpts, to common.Address, medication string, dosage *big.Int) (*big.Int, error)
	Burn(opts *bind.CallOpts, tokenId *big.Int) error
}

// PrescriptionNFTTransactor is the interface for contract transactions
type PrescriptionNFTTransactor interface {
	Mint(opts *bind.TransactOpts, to common.Address, medication string, dosage *big.Int) (*types.Transaction, error)
	Burn(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error)
}

// Client manages blockchain interactions with Polygon for PrescriptionNFT
type Client struct {
	EthClient    *ethclient.Client
	Contract     *prescriptionnft.PrescriptionNFT // Generated binding
	ContractAddr common.Address
	ChainID      *big.Int
	PrivateKey   *ecdsa.PrivateKey // For signing transactions
	FromAddress  common.Address    // Sender’s address
}

// NewClient initializes a new Polygon blockchain client
func NewClient(config *configs.Config) (*Client, error) {
	// Connect to Polygon RPC (e.g., Mumbai testnet)
	ethClient, err := ethclient.Dial(config.Blockchain.RPCURL)
	if err != nil {
		log.Printf("Failed to connect to Polygon RPC: %v", err)
		return nil, err
	}

	// Get chain ID to verify network
	chainID, err := ethClient.NetworkID(context.Background())
	if err != nil {
		log.Printf("Failed to get chain ID: %v", err)
		return nil, err
	}

	// Convert contract address string to common.Address
	contractAddr := common.HexToAddress(config.Blockchain.ContractAddress)

	// Load private key from environment variable for transaction signing (securely manage this!)
	privateKey, err := crypto.HexToECDSA(os.Getenv("POLYGON_PRIVATE_KEY"))
	if err != nil {
		log.Printf("Failed to load private key: %v", err)
		return nil, err
	}

	// Derive public address from private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, logError("Invalid public key type")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Bind to PrescriptionNFT contract (requires generated binding)
	contract, err := prescriptionnft.NewPrescriptionNFT(contractAddr, ethClient)
	if err != nil {
		log.Printf("Failed to bind PrescriptionNFT contract: %v", err)
		return nil, err
	}

	return &Client{
		EthClient:    ethClient,
		Contract:     contract,
		ContractAddr: contractAddr,
		ChainID:      chainID,
		PrivateKey:   privateKey,
		FromAddress:  fromAddress,
	}, nil
}

// MintPrescription mints a new prescription NFT for a patient
func (c *Client) MintPrescription(patientAddr string, medication string, dosage uint64) (uint64, error) {
	// Create transaction options with private key
	auth, err := bind.NewKeyedTransactorWithChainID(c.PrivateKey, c.ChainID)
	if err != nil {
		log.Printf("Failed to create transactor: %v", err)
		return 0, err
	}

	// Convert patient address to common.Address
	to := common.HexToAddress(patientAddr)

	// Mint the NFT (tokenID is returned as *big.Int)
	tokenID, err := c.Contract.Mint(auth, to, medication, big.NewInt(int64(dosage)))
	if err != nil {
		log.Printf("Failed to mint prescription NFT: %v", err)
		return 0, err
	}

	return tokenID.Int64(), nil
}

// DispensePrescription burns a prescription NFT
func (c *Client) DispensePrescription(tokenID uint64) error {
	// Create transaction options with private key
	auth, err := bind.NewKeyedTransactorWithChainID(c.PrivateKey, c.ChainID)
	if err != nil {
		log.Printf("Failed to create transactor: %v", err)
		return err
	}

	// Burn the NFT
	tx, err := c.Contract.Burn(auth, big.NewInt(int64(tokenID)))
	if err != nil {
		log.Printf("Failed to burn prescription NFT: %v", err)
		return err
	}

	log.Printf("Dispensed prescription NFT (tokenID: %d), transaction: %s", tokenID, tx.Hash().Hex())
	return nil
}

// logError helper function for consistent error logging
func logError(msg string) error {
	err := fmt.Errorf(msg)
	log.Println(err)
	return err
}
