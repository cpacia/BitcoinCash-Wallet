package spvwallet

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"time"
)

type Checkpoint struct {
	Height uint32
	Header wire.BlockHeader
}

var mainnetCheckpoints []Checkpoint
var testnet3Checkpoints []Checkpoint
var regtestCheckpoint Checkpoint

var BitcoinCashForkBlock *chainhash.Hash

func init() {
	// Bitcoin Cash Fork
	BitcoinCashForkBlock, _ = chainhash.NewHashFromStr("000000000000000000651ef99cb9fcbe0dadde1d424bd9f15ff20136191a5eec")

	// Mainnet
	mainnetPrev, _ := chainhash.NewHashFromStr("0000000000000000010e0373719b7538e713e47d8d7189826dce4264d85a79b8")
	mainnetMerk, _ := chainhash.NewHashFromStr("28e3ae14925aeead5b44a885bca88e06ec37483ae905e65d57aa8cd60cb74b2f")
	mainnetCheckpoints = append(mainnetCheckpoints, Checkpoint{
		Height: 477792,
		Header: wire.BlockHeader{
			Version:    536870914,
			PrevBlock:  *mainnetPrev,
			MerkleRoot: *mainnetMerk,
			Timestamp:  time.Unix(1501153434, 0),
			Bits:       402736949,
			Nonce:      165193412,
		},
	})
	if mainnetCheckpoints[0].Header.BlockHash().String() != "00000000000000000016ba7786309176445b838b36a16bd1ef3c3e3020473206" {
		panic("Invalid checkpoint")
	}

	// Testnet3
	testnet3Prev, _ := chainhash.NewHashFromStr("0000000000001e8cdb2d98471a5c60bdbddbe644b9ad08e17a97b3a7dce1e332")
	testnet3Merk, _ := chainhash.NewHashFromStr("f675c565b293be2ad808b01b0a763557c8874e4aefe7f2eea0dab91b1f60ec45")
	testnet3Checkpoints = append(testnet3Checkpoints, Checkpoint{
		Height: 1151136,
		Header: wire.BlockHeader{
			Version:    536870912,
			PrevBlock:  *testnet3Prev,
			MerkleRoot: *testnet3Merk,
			Timestamp:  time.Unix(1498950206, 0),
			Bits:       436724869,
			Nonce:      2247874206,
		},
	})
	if testnet3Checkpoints[0].Header.BlockHash().String() != "00000000000002c04de174cf25c993b4dd221eb087c0601a599ff1977e230c99" {
		panic("Invalid checkpoint")
	}

	// Regtest
	regtestCheckpoint = Checkpoint{0, chaincfg.RegressionNetParams.GenesisBlock.Header}
}

func GetCheckpoint(walletCreationDate time.Time, params *chaincfg.Params) Checkpoint {
	switch params.Name {
	case chaincfg.MainNetParams.Name:
		for i := len(mainnetCheckpoints) - 1; i >= 0; i-- {
			if walletCreationDate.After(mainnetCheckpoints[i].Header.Timestamp) {
				return mainnetCheckpoints[i]
			}
		}
		return mainnetCheckpoints[0]
	case chaincfg.TestNet3Params.Name:
		for i := len(testnet3Checkpoints) - 1; i >= 0; i-- {
			if walletCreationDate.After(testnet3Checkpoints[i].Header.Timestamp) {
				return testnet3Checkpoints[i]
			}
		}
		return testnet3Checkpoints[0]

	default:
		return regtestCheckpoint
	}
}
