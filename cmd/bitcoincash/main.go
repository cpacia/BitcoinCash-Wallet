package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OpenBazaar/wallet-interface"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
	"github.com/atotto/clipboard"
	"github.com/gcash/bchd/bchec"
	"github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchutil"
	bc "github.com/bubbajoe/BitcoinCash-Wallet"
	"github.com/bubbajoe/BitcoinCash-Wallet/api"
	"github.com/bubbajoe/BitcoinCash-Wallet/cli"
	"github.com/bubbajoe/BitcoinCash-Wallet/db"
	"github.com/bubbajoe/BitcoinCash-Wallet/exchangerates"
	"github.com/bubbajoe/BitcoinCash-Wallet/gui"
	"github.com/bubbajoe/BitcoinCash-Wallet/gui/bootstrap"
	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
	"github.com/natefinch/lumberjack"
	"github.com/op/go-logging"
	"github.com/skratchdot/open-golang/open"
	"github.com/yawning/bulb"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"time"
)

var parser = flags.NewParser(nil, flags.Default)

type Start struct {
	DataDir            string `short:"d" long:"datadir" description:"specify the data directory to be used"`
	Testnet            bool   `short:"t" long:"testnet" description:"use the test network"`
	Regtest            bool   `short:"r" long:"regtest" description:"run in regression test mode"`
	Mnemonic           string `short:"m" long:"mnemonic" description:"specify a mnemonic seed to use to derive the keychain"`
	WalletCreationDate string `short:"w" long:"walletcreationdate" description:"specify the date the seed was created. if omitted the wallet will sync from the oldest checkpoint."`
	TrustedPeer        string `short:"i" long:"trustedpeer" description:"specify a single trusted peer to connect to"`
	Tor                bool   `long:"tor" description:"connect via a running Tor daemon"`
	FeeAPI             string `short:"f" long:"feeapi" description:"fee API to use to fetch current fee rates. set as empty string to disable API lookups." default:""`
	MaxFee             uint64 `short:"x" long:"maxfee" description:"the fee-per-byte ceiling beyond which fees cannot go" default:"2000"`
	LowDefaultFee      uint64 `short:"e" long:"economicfee" description:"the default low fee-per-byte" default:"20"`
	MediumDefaultFee   uint64 `short:"n" long:"normalfee" description:"the default medium fee-per-byte" default:"90"`
	HighDefaultFee     uint64 `short:"p" long:"priorityfee" description:"the default high fee-per-byte" default:"180"`
	Gui                bool   `long:"gui" description:"launch an experimental GUI"`
	Verbose            bool   `short:"v" long:"verbose" description:"print to standard out"`
}
type Version struct{}

var start Start
var version Version
var cashWallet *bc.SPVWallet

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("BitcoinCash wallet shutting down...")
			cashWallet.Close()
			os.Exit(1)
		}
	}()
	if len(os.Args) == 1 {
		start.Gui = true
		start.Execute([]string{"useDefaults"})
	} else {
		parser.AddCommand("start",
			"start the wallet",
			"The start command starts the wallet daemon",
			&start)
		parser.AddCommand("version",
			"print the version number",
			"Print the version number and exit",
			&version)
		cli.SetupCli(parser)
		if _, err := parser.Parse(); err != nil {
			os.Exit(1)
		}
	}
}

func (x *Version) Execute(args []string) error {
	fmt.Println(bc.WALLET_VERSION)
	return nil
}

func (x *Start) Execute(args []string) error {
	var err error
	// Create a new config
	config := bc.NewDefaultConfig()
	if x.DataDir != "" {
		config.RepoPath = x.DataDir
	}
	if x.Testnet && x.Regtest {
		return errors.New("Invalid combination of testnet and regtest modes")
	}
	basepath := config.RepoPath
	if x.Testnet {
		config.Params = &chaincfg.TestNet3Params
		config.RepoPath = path.Join(config.RepoPath, "testnet")
	}
	if x.Regtest {
		config.Params = &chaincfg.RegressionNetParams
		config.RepoPath = path.Join(config.RepoPath, "regtest")
	}

	_, ferr := os.Stat(config.RepoPath)
	if os.IsNotExist(ferr) {
		os.Mkdir(config.RepoPath, os.ModePerm)
	}
	if x.Mnemonic != "" {
		config.Mnemonic = x.Mnemonic
	}
	if x.TrustedPeer != "" {
		addr, err := net.ResolveTCPAddr("tcp", x.TrustedPeer)
		if err != nil {
			return err
		}
		config.TrustedPeer = addr
	}
	if x.Tor {
		var conn *bulb.Conn
		conn, err = bulb.Dial("tcp4", "127.0.0.1:9151")
		if err != nil {
			conn, err = bulb.Dial("tcp4", "127.0.0.1:9151")
			if err != nil {
				return errors.New("Tor daemon not found")
			}
		}
		dialer, err := conn.Dialer(nil)
		if err != nil {
			return err
		}
		config.Proxy = dialer
	}
	if x.FeeAPI != "" {
		u, err := url.Parse(x.FeeAPI)
		if err != nil {
			return err
		}
		config.FeeAPI = *u
	}

	if len(args) == 0 {
		config.MaxFee = x.MaxFee
		config.LowFee = x.LowDefaultFee
		config.MediumFee = x.MediumDefaultFee
		config.HighFee = x.HighDefaultFee
	}

	// Make the logging a little prettier
	var fileLogFormat = logging.MustStringFormatter(`%{time:15:04:05.000} [%{shortfunc}] [%{level}] %{message}`)
	var stdoutLogFormat = logging.MustStringFormatter(`%{color:reset}%{color}%{time:15:04:05.000} [%{shortfunc}] [%{level}] %{message}`)
	w := &lumberjack.Logger{
		Filename:   path.Join(config.RepoPath, "logs", "bitcoin.log"),
		MaxSize:    10, // Megabytes
		MaxBackups: 3,
		MaxAge:     30, // Days
	}
	bitcoinFile := logging.NewLogBackend(w, "", 0)
	bitcoinFileFormatter := logging.NewBackendFormatter(bitcoinFile, fileLogFormat)
	config.Logger = logging.MultiLogger(logging.MultiLogger(bitcoinFileFormatter))
	if x.Verbose {
		stdoutLog := logging.NewLogBackend(os.Stdout, "", 0)
		stdoutFormatter := logging.NewBackendFormatter(stdoutLog, stdoutLogFormat)
		config.Logger = logging.MultiLogger(logging.MultiLogger(stdoutFormatter))
	}

	// Select wallet datastore
	sqliteDatastore, _ := db.Create(config.RepoPath)
	config.DB = sqliteDatastore

	mn, _ := sqliteDatastore.GetMnemonic()
	if mn != "" {
		config.Mnemonic = mn
	}

	// Write version file
	f, err := os.Create(path.Join(basepath, "version"))
	if err != nil {
		return err
	}
	f.Write([]byte("1"))
	f.Close()

	exchangeRates := exchangerates.NewBitcoinCashPriceFetcher(config.Proxy)
	config.ExchangeRateProvider = exchangeRates

	// Load settings
	type Fees struct {
		Priority      uint64 `json:"priority"`
		Normal        uint64 `json:"normal"`
		Economic      uint64 `json:"economic"`
		MaxFee        uint64 `json:"maxFee"`
		feeEstimation bool   `json:"feeEstimation"`
	}

	type Settings struct {
		FiatCode      string `json:"fiatCode"`
		FiatSymbol    string `json:"fiatSymbol"`
		FeeLevel      string `json:"feeLevel"`
		SelectBox     string `json:"selectBox"`
		BitcoinUnit   string `json:"bitcoinUnit"`
		DecimalPlaces int    `json:"decimalPlaces"`
		TrustedPeer   string `json:"trustedPeer"`
		Proxy         string `json:"proxy"`
		Fees          Fees   `json:"fees"`
	}

	var settings Settings
	s, err := ioutil.ReadFile(path.Join(basepath, "settings.json"))
	if err != nil {
		settings = Settings{
			FiatCode:      "USD",
			FiatSymbol:    "$",
			FeeLevel:      "priority",
			SelectBox:     "bitcoin",
			BitcoinUnit:   "BCH",
			DecimalPlaces: 5,
			Fees: Fees{
				Priority: config.HighFee,
				Normal:   config.MediumFee,
				Economic: config.LowFee,
			},
		}
		f, err := os.Create(path.Join(basepath, "settings.json"))
		if err != nil {
			return err
		}
		s, err := json.MarshalIndent(&settings, "", "    ")
		if err != nil {
			return err
		}
		f.Write(s)
		f.Close()
	} else {
		err := json.Unmarshal([]byte(s), &settings)
		if err != nil {
			return err
		}
		switch settings.BitcoinUnit {
		case "BTC":
			settings.BitcoinUnit = "BCH"
		case "mBTC":
			settings.BitcoinUnit = "mBCH"
		}
	}
	if settings.TrustedPeer != "" {
		var tp net.Addr
		tp, err = net.ResolveTCPAddr("tcp", settings.TrustedPeer)
		if err != nil {
			return err
		}
		config.TrustedPeer = tp
	}

	if settings.Proxy != "" {
		tbProxyURL, err := url.Parse("socks5://" + settings.Proxy)
		if err != nil {
			return err
		}
		tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
		if err != nil {
			return err
		}
		config.Proxy = tbDialer
	}
	config.HighFee = settings.Fees.Priority
	config.MediumFee = settings.Fees.Normal
	config.LowFee = settings.Fees.Economic

	creationDate := time.Time{}
	if x.WalletCreationDate != "" {
		creationDate, err = time.Parse(time.RFC3339, x.WalletCreationDate)
		if err != nil {
			return errors.New("Wallet creation date timestamp must be in RFC3339 format")
		}
	} else {
		creationDate, _ = sqliteDatastore.GetCreationDate()
	}
	config.CreationDate = creationDate

	// Create the wallet
	cashWallet, err = bc.NewSPVWallet(config)
	if err != nil {
		return err
	}

	if err := sqliteDatastore.SetMnemonic(config.Mnemonic); err != nil {
		return err
	}
	if err := sqliteDatastore.SetCreationDate(config.CreationDate); err != nil {
		return err
	}

	go api.ServeAPI(cashWallet)

	// Start it!
	printSplashScreen()

	if x.Gui {
		go cashWallet.Start()

		type Stats struct {
			Confirmed    int64  `json:"confirmed"`
			Fiat         string `json:"fiat"`
			Transactions int    `json:"transactions"`
			Height       uint32 `json:"height"`
			ExchangeRate string `json:"exchangeRate"`
		}

		txc := make(chan uint32)
		listener := func(wallet.TransactionCallback) {
			h, _ := cashWallet.ChainTip()
			txc <- h
		}
		cashWallet.AddTransactionListener(listener)

		tc := make(chan struct{})
		rc := make(chan int)

		os.RemoveAll(path.Join(basepath, "resources"))
		iconPath := path.Join(basepath, "icon.png")
		_, err := os.Stat(iconPath)
		if os.IsNotExist(err) {
			f, err := os.Create(iconPath)
			if err != nil {
				return err
			}
			icon, err := gui.AppIconPngBytes()
			if err != nil {
				return err
			}
			f.Write(icon)
			defer f.Close()
		}

		// Run bootstrap
		if err := bootstrap.Run(bootstrap.Options{
			AstilectronOptions: astilectron.Options{
				AppName:            "Bitcoin Cash Wallet",
				AppIconDefaultPath: iconPath,
				//AppIconDarwinPath:  p + "/gopher.icns",
				BaseDirectoryPath: basepath,
			},
			Homepage: "index.html",
			MessageHandler: func(w *astilectron.Window, m bootstrap.MessageIn) {
				switch m.Name {
				case "getStats":
					type P struct {
						CurrencyCode string `json:"currencyCode"`
					}
					var p P
					if err := json.Unmarshal(m.Payload, &p); err != nil {
						astilog.Errorf("Unmarshaling %s failed", m.Payload)
						return
					}
					confirmed, _ := cashWallet.Balance()
					txs, err := cashWallet.Transactions()
					if err != nil {
						astilog.Errorf(err.Error())
						return
					}
					rate, err := exchangeRates.GetExchangeRate(p.CurrencyCode)
					if err != nil {
						astilog.Errorf("Failed to get exchange rate")
						return
					}
					bchVal := float64(confirmed) / 100000000
					fiatVal := float64(bchVal) * rate
					height, _ := cashWallet.ChainTip()

					st := Stats{
						Confirmed:    confirmed,
						Fiat:         fmt.Sprintf("%.2f", fiatVal),
						Transactions: len(txs),
						Height:       height,
						ExchangeRate: fmt.Sprintf("%.2f", rate),
					}
					w.SendMessage(bootstrap.MessageOut{Name: "statsUpdate", Payload: st})
				case "getAddress":
					addr := cashWallet.CurrentAddress(wallet.EXTERNAL)
					w.SendMessage(bootstrap.MessageOut{Name: "address", Payload: addr.EncodeAddress()})
				case "send":
					type P struct {
						Address  string  `json:"address"`
						Amount   float64 `json:"amount"`
						Note     string  `json:"note"`
						FeeLevel string  `json:"feeLevel"`
					}
					var p P
					if err := json.Unmarshal(m.Payload, &p); err != nil {
						astilog.Errorf("Unmarshaling %s failed", m.Payload)
						return
					}
					var feeLevel wallet.FeeLevel
					switch strings.ToLower(p.FeeLevel) {
					case "priority":
						feeLevel = wallet.PRIOIRTY
					case "normal":
						feeLevel = wallet.NORMAL
					case "economic":
						feeLevel = wallet.ECONOMIC
					default:
						feeLevel = wallet.NORMAL
					}
					addr, err := bchutil.DecodeAddress(p.Address, cashWallet.Params())
					if err != nil {
						w.SendMessage(bootstrap.MessageOut{Name: "spendError", Payload: "Invalid address"})
						return
					}
					_, err = cashWallet.Spend(int64(p.Amount), addr, feeLevel)
					if err != nil {
						w.SendMessage(bootstrap.MessageOut{Name: "spendError", Payload: err.Error()})
					}
				case "clipboard":
					type P struct {
						Data string `json:"data"`
					}
					var p P
					if err := json.Unmarshal(m.Payload, &p); err != nil {
						astilog.Errorf("Unmarshaling %s failed", m.Payload)
						return
					}
					clipboard.WriteAll(p.Data)
				case "putSettings":
					var setstr string
					if err := json.Unmarshal(m.Payload, &setstr); err != nil {
						astilog.Errorf("Unmarshaling %s failed", m.Payload)
						return
					}
					var settings Settings
					if err := json.Unmarshal([]byte(setstr), &settings); err != nil {
						astilog.Errorf("Unmarshaling %s failed", m.Payload)
						return
					}

					f, err := os.Create(path.Join(basepath, "settings.json"))
					if err != nil {
						astilog.Error(err.Error())
						return
					}
					defer f.Close()
					b, err := json.MarshalIndent(&settings, "", "    ")
					if err != nil {
						astilog.Error(err.Error())
						return
					}
					f.Write(b)
				case "getSettings":
					settings, err := ioutil.ReadFile(path.Join(basepath, "settings.json"))
					if err != nil {
						astilog.Error(err.Error())
					}

					var set Settings
					err = json.Unmarshal(settings, &set)
					if err != nil {
						astilog.Error(err.Error())
					}

					switch set.BitcoinUnit {
					case "BTC":
						set.BitcoinUnit = "BCH"
					case "mBTC":
						set.BitcoinUnit = "mBCH"
					}
					ret, err := json.MarshalIndent(set, "", "    ")
					if err != nil {
						astilog.Error(err.Error())
					}
					w.SendMessage(bootstrap.MessageOut{Name: "settings", Payload: string(ret)})
				case "openbrowser":
					var url string
					if err := json.Unmarshal(m.Payload, &url); err != nil {
						astilog.Errorf("Unmarshaling %s failed", m.Payload)
						return
					}
					open.Run(url)
				case "resync":
					cashWallet.ReSyncBlockchain(config.CreationDate)
				case "importKey":
					type P struct {
						Key  string `json:"key"`
						Date string `json:"date"`
					}
					var p P
					if err := json.Unmarshal(m.Payload, &p); err != nil {
						astilog.Errorf("Unmarshaling %s failed", m.Payload)
						return
					}
					var privKey *bchec.PrivateKey
					compress := true
					keyBytes, err := hex.DecodeString(p.Key)
					if err != nil {
						wif, err := bchutil.DecodeWIF(p.Key)
						if err != nil {
							w.SendMessage(bootstrap.MessageOut{Name: "importError", Payload: err.Error()})
							return
						}
						privKey = wif.PrivKey
						compress = wif.CompressPubKey
					} else {
						if len(keyBytes) != 32 {
							w.SendMessage(bootstrap.MessageOut{Name: "importError", Payload: "Invalid key"})
							return
						}
						privKey, _ = bchec.PrivKeyFromBytes(bchec.S256(), keyBytes)
					}
					cashWallet.ImportKey(privKey, compress)
					var t time.Time
					if p.Date != "" {
						t, _ = time.Parse("2006-01-2", p.Date)
					}
					cashWallet.ReSyncBlockchain(t)
				case "restore":
					type P struct {
						Mnemonic string `json:"mnemonic"`
						Date     string `json:"date"`
					}
					var p P
					if err := json.Unmarshal(m.Payload, &p); err != nil {
						astilog.Errorf("Unmarshaling %s failed", m.Payload)
						return
					}
					var t time.Time
					if p.Date != "" {
						t, _ = time.Parse("2006-01-2", p.Date)
					}

					cashWallet.Close()
					os.Remove(path.Join(config.RepoPath, "wallet.db"))
					os.Remove(path.Join(config.RepoPath, "headers.bin"))
					sqliteDatastore, _ := db.Create(config.RepoPath)
					config.DB = sqliteDatastore
					config.Mnemonic = p.Mnemonic
					config.CreationDate = t
					cashWallet, err = bc.NewSPVWallet(config)
					if err != nil {
						astilog.Errorf("Unmarshaling %s failed", m.Payload)
						return
					}
					sqliteDatastore.SetMnemonic(p.Mnemonic)
					sqliteDatastore.SetCreationDate(t)
					go cashWallet.Start()
				case "minimize":
					go func() {
						w.Hide()
						tc <- struct{}{}
					}()
				case "showTransactions":
					go func() {
						rc <- 649
					}()
					txs, err := cashWallet.Transactions()
					if err != nil {
						w.SendMessage(bootstrap.MessageOut{Name: "txError", Payload: err.Error()})
					}
					w.SendMessage(bootstrap.MessageOut{Name: "transactions", Payload: txs})
				case "getTransactions":
					txs, err := cashWallet.Transactions()
					if err != nil {
						w.SendMessage(bootstrap.MessageOut{Name: "txError", Payload: err.Error()})
					}
					w.SendMessage(bootstrap.MessageOut{Name: "transactions", Payload: txs})
				case "hide":
					go func() {
						rc <- 341
					}()
				case "showSettings":
					go func() {
						rc <- 649
					}()
				case "getMnemonic":
					w.SendMessage(bootstrap.MessageOut{Name: "mnemonic", Payload: cashWallet.Mnemonic()})
				}
			},
			RestoreAssets: gui.RestoreAssets,
			WindowOptions: &astilectron.WindowOptions{
				Center:         astilectron.PtrBool(true),
				Height:         astilectron.PtrInt(340),
				Width:          astilectron.PtrInt(621),
				Maximizable:    astilectron.PtrBool(false),
				Fullscreenable: astilectron.PtrBool(false),
				Resizable:      astilectron.PtrBool(false),
			},
			TrayOptions: &astilectron.TrayOptions{
				Image: astilectron.PtrStr(iconPath),
			},
			TrayChan:          tc,
			ResizeChan:        rc,
			TransactionChan:   txc,
			BaseDirectoryPath: basepath,
			Wallet:            cashWallet,
			//Debug:             true,
		}); err != nil {
			astilog.Fatal(err)
		}
	} else {
		cashWallet.Start()
		var wg sync.WaitGroup
		wg.Add(1)
		wg.Wait()
	}
	return nil
}

func printSplashScreen() {
	blue := color.New(color.FgBlue)
	white := color.New(color.FgWhite)
	white.Printf("  _______________________   ______")
	blue.Println("      __        .__  .__          __")
	white.Printf(` /   _____/\______   \   \ /   /`)
	blue.Println(`  \    /  \_____  |  | |  |   _____/  |_`)
	white.Printf(` \_____  \  |     ___/\   Y   /`)
	blue.Println(`\   \/\/   /\__  \ |  | |  | _/ __ \   __\`)
	white.Printf(` /        \ |    |     \     / `)
	blue.Println(` \        /  / __ \|  |_|  |_\  ___/|  |`)
	white.Printf(`/_______  / |____|      \___/ `)
	blue.Println(`   \__/\  /  (____  /____/____/\___  >__|`)
	white.Printf(`	\/ `)
	blue.Println(`                           \/        \/               \/`)
	blue.DisableColor()
	white.DisableColor()
	fmt.Println("")
	fmt.Println("BitcoinCash wallet v" + bc.WALLET_VERSION + " starting...")
	fmt.Println("[Press Ctrl+C to exit]")
}
