package main

import (
	"fmt"
	"os"

	"github.com/gcash/bchd/chaincfg"
	"github.com/bubbajoe/BitcoinCash-Wallet"
	"github.com/bubbajoe/BitcoinCash-Wallet/db"
	"github.com/op/go-logging"
)

func main() {
	// Create a new config
	config := bitcoincash.NewDefaultConfig()

	// Make the logging a little prettier
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	formatter := logging.MustStringFormatter(`%{color:reset}%{color}%{time:15:04:05.000} [%{shortfunc}] [%{level}] %{message}`)
	stdoutFormatter := logging.NewBackendFormatter(backend, formatter)
	config.Logger = logging.MultiLogger(stdoutFormatter)

	// Use testnet
	config.Params = &chaincfg.TestNet3Params

	// Select wallet datastore
	sqliteDatastore, _ := db.Create(config.RepoPath)
	config.DB = sqliteDatastore

	// Create the wallet
	wallet, err := bitcoincash.NewSPVWallet(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Start it!
	wallet.Start()
}
