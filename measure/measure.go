package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
)

type Datapoint struct {
	BlockHash       string   `json:"hash"`
	BlockCount      int64    `json:"blockCount"`
	ConnectionCount int64    `json:"connectionCount"`
	TimeElapsed     float64  `json:"timeElapsed"`
	ConnectedPeers  []string `json:"connectedPeers"`
}

type DataContainer struct {
	Data []Datapoint `json:"data"`
}

func main() {
	// Only override the handlers for notifications you care about.
	// Also note most of these handlers will only be called if you register
	// for notifications.  See the documentation of the rpcclient
	// NotificationHandlers type for more details about each handler.
	ntfnHandlers := rpcclient.NotificationHandlers{
		OnFilteredBlockConnected: func(height int32, header *wire.BlockHeader, txns []*btcutil.Tx) {
			log.Printf("Block connected: %v (%d) %v",
				header.BlockHash(), height, header.Timestamp)
		},
		OnFilteredBlockDisconnected: func(height int32, header *wire.BlockHeader) {
			log.Printf("Block disconnected: %v (%d) %v",
				header.BlockHash(), height, header.Timestamp)
		},
	}

	log.Println("Registered the notification handler")
	// Connect to local btcd RPC server using websockets.
	btcdHomeDir := btcutil.AppDataDir("btcd", false)
	certs, err := ioutil.ReadFile(filepath.Join(btcdHomeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:8334",
		Endpoint:     "ws",
		User:         "btcd",
		Pass:         "123",
		Certificates: certs,
	}
	client, err := rpcclient.New(connCfg, &ntfnHandlers)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to RPC server")

	// Register for block connect and disconnect notifications.
	if err := client.NotifyBlocks(); err != nil {
		log.Fatal(err)
	}
	log.Println("NotifyBlocks: Registration Complete")

	//////////////////////////////////////////////////////////////////
	dataContainer := DataContainer{}

	// Initialize the start time
	startTime := time.Now()

	// Start a loop that runs every 10 seconds
	for {
		// Calculate the time elapsed since the first loop iteration
		timeElapsed := time.Since(startTime).Seconds()

		// newBlockHash, err := client.Generate(1)
		// if err != nil {
		// 	print(err, newBlockHash)
		// 	return
		// }

		newestBlockHash, err := client.GetBestBlockHash()
		if err != nil {
			print(err)
			return
		}

		connectionCount, err := client.GetConnectionCount()
		if err != nil {
			print(err)
			return
		}

		blockCount, err := client.GetBlockCount()
		if err != nil {
			print(err)
			return
		}

		peerJSON, err := client.GetPeerInfo()
		if err != nil {
			print(err)
			return
		}

		connectedPeers := []string{}
		for _, peers := range peerJSON {
			connectedPeers = append(connectedPeers, peers.Addr)
		}

		// create data point
		newDatapoint := Datapoint{
			BlockHash:       newestBlockHash.String(),
			BlockCount:      blockCount,
			ConnectionCount: connectionCount,
			TimeElapsed:     timeElapsed,
			ConnectedPeers:  connectedPeers,
		}
		dataContainer.Data = append(dataContainer.Data, newDatapoint)

		log.Println("Writing datapoint to JSON")
		print(timeElapsed)
		// Write data to the JSON file
		if err := writeDataToFile(dataContainer, "/root/.btcd/logs/mainnet/log.json"); err != nil {
			log.Printf("Error writing data to file: %v", err)
		}
		time.Sleep(10 * time.Second)
	}

}
func writeDataToFile(dataContainer DataContainer, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(dataContainer); err != nil {
		return err
	}

	return nil
}
