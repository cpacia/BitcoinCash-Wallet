package bitcoincash

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

var BitcoinCashMainnetForkBlock *chainhash.Hash
var BitcoinCashTestnet3ForkBlock *chainhash.Hash

func init() {
	// Bitcoin Cash Fork
	BitcoinCashMainnetForkBlock, _ = chainhash.NewHashFromStr("000000000000000000651ef99cb9fcbe0dadde1d424bd9f15ff20136191a5eec")

	// Mainnet
	mainnetCheckpoints = append(mainnetCheckpoints, Checkpoint{
		Height: 0,
		Header: chaincfg.MainNetParams.GenesisBlock.Header,
	})

	mainnetPrev1, _ := chainhash.NewHashFromStr("00000000103226f85fe2b68795f087dcec345e523363f18017e60b5c94175355")
	mainnetMerk1, _ := chainhash.NewHashFromStr("f90471766508bfe36158b4f2a5209dc5609d496f01b2d9b8f8fef1869fd8d382")
	mainnetCheckpoints = append(mainnetCheckpoints, Checkpoint{
		Height: 48384,
		Header: wire.BlockHeader{
			Version:    1,
			PrevBlock:  *mainnetPrev1,
			MerkleRoot: *mainnetMerk1,
			Timestamp:  time.Unix(1270120042, 0),
			Bits:       472518933,
			Nonce:      317008837,
		},
	})
	if mainnetCheckpoints[1].Header.BlockHash().String() != "000000000f8bad9641203f9d788e72d55ff164f1a62a3e230df6f411f7d7d8fc" {
		panic("Invalid checkpoint 1")
	}

	mainnetPrev2, _ := chainhash.NewHashFromStr("00000000000292b850b8f8578e6b4d03cbb4a78ada44afbb4d2f80a16490e8f9")
	mainnetMerk2, _ := chainhash.NewHashFromStr("d1b77218087bf1c892d76ac0a6591e411322fb3399ee52a47b300158c61a6993")
	mainnetCheckpoints = append(mainnetCheckpoints, Checkpoint{
		Height: 98784,
		Header: wire.BlockHeader{
			Version:    1,
			PrevBlock:  *mainnetPrev2,
			MerkleRoot: *mainnetMerk2,
			Timestamp:  time.Unix(1292956443, 0),
			Bits:       453281356,
			Nonce:      553309409,
		},
	})
	if mainnetCheckpoints[2].Header.BlockHash().String() != "0000000000043c4b0ed186109c36ce09e9af29b8ed0fdebcb909fa338877743d" {
		panic("Invalid checkpoint 2")
	}

	mainnetPrev3, _ := chainhash.NewHashFromStr("00000000000001b0d8d885e4d77d7c51e8f1fdaba68f229ac04d191915845f09")
	mainnetMerk3, _ := chainhash.NewHashFromStr("6664abb20174f8efcc8ae76b85c571919f6376c9083600d5e3fce7b94a985b12")
	mainnetCheckpoints = append(mainnetCheckpoints, Checkpoint{
		Height: 149184,
		Header: wire.BlockHeader{
			Version:    1,
			PrevBlock:  *mainnetPrev3,
			MerkleRoot: *mainnetMerk3,
			Timestamp:  time.Unix(1318556675, 0),
			Bits:       436956491,
			Nonce:      1274332420,
		},
	})
	if mainnetCheckpoints[3].Header.BlockHash().String() != "000000000000028c550f87b05cddb1f4d461451bb26f0706f1ed5b20e3889fb9" {
		panic("Invalid checkpoint 3")
	}

	mainnetPrev4, _ := chainhash.NewHashFromStr("0000000000000304760bf583f4ac241c5ffe77312fa213634eba252c720530f1")
	mainnetMerk4, _ := chainhash.NewHashFromStr("4d91953d90f8ca7a60c40ec4f40a6bcd3832e3a6a079297352777d998371793c")
	mainnetCheckpoints = append(mainnetCheckpoints, Checkpoint{
		Height: 199584,
		Header: wire.BlockHeader{
			Version:    1,
			PrevBlock:  *mainnetPrev4,
			MerkleRoot: *mainnetMerk4,
			Timestamp:  time.Unix(1348092851, 0),
			Bits:       436591499,
			Nonce:      887613184,
		},
	})
	if mainnetCheckpoints[4].Header.BlockHash().String() != "000000000000002e00a243fe9aa49c78f573091d17372c2ae0ae5e0f24f55b52" {
		panic("Invalid checkpoint 4")
	}

	mainnetPrev5, _ := chainhash.NewHashFromStr("0000000000000055e146e473b49fe656a1f2f4b8c33e72b80acc18f84d9fcc26")
	mainnetMerk5, _ := chainhash.NewHashFromStr("3e21fe7f80ac41dbcbdeb4f13fe9a0d78157625ba7c0d797d39918a5929bbf51")
	mainnetCheckpoints = append(mainnetCheckpoints, Checkpoint{
		Height: 249984,
		Header: wire.BlockHeader{
			Version:    2,
			PrevBlock:  *mainnetPrev5,
			MerkleRoot: *mainnetMerk5,
			Timestamp:  time.Unix(1375527115, 0),
			Bits:       426957810,
			Nonce:      3431028336,
		},
	})
	if mainnetCheckpoints[5].Header.BlockHash().String() != "000000000000006f68295a01c059e45c38fe3ff5095107f067c52ae6437c3d4b" {
		panic("Invalid checkpoint 5")
	}

	mainnetPrev6, _ := chainhash.NewHashFromStr("000000000000000098686ab04cc22fec77e4fa2d76d5a3cc0eb8cbf4ed800cdc")
	mainnetMerk6, _ := chainhash.NewHashFromStr("e9087641b6f19e49dc37be1d35ec5b670b1baa4529883602b59068e1799adb44")
	mainnetCheckpoints = append(mainnetCheckpoints, Checkpoint{
		Height: 298368,
		Header: wire.BlockHeader{
			Version:    2,
			PrevBlock:  *mainnetPrev6,
			MerkleRoot: *mainnetMerk6,
			Timestamp:  time.Unix(1398811175, 0),
			Bits:       419465580,
			Nonce:      952935459,
		},
	})
	if mainnetCheckpoints[6].Header.BlockHash().String() != "0000000000000000092dbe12467f95d456d9bd7ff12b88dc3bcdd28b4a7aa159" {
		panic("Invalid checkpoint 6")
	}

	mainnetPrev7, _ := chainhash.NewHashFromStr("00000000000000000b3ec9df7aebc319bb12491ba651337f9b3541e78446eca8")
	mainnetMerk7, _ := chainhash.NewHashFromStr("9bd3929bd9b575d103c43e5dc131f2956f1a8925f4f2596e603bd5f64deabbdf")
	mainnetCheckpoints = append(mainnetCheckpoints, Checkpoint{
		Height: 348768,
		Header: wire.BlockHeader{
			Version:    2,
			PrevBlock:  *mainnetPrev7,
			MerkleRoot: *mainnetMerk7,
			Timestamp:  time.Unix(1427068411, 0),
			Bits:       404195570,
			Nonce:      903044756,
		},
	})
	if mainnetCheckpoints[7].Header.BlockHash().String() != "0000000000000000145865530e632fdcc732e2309645d759d738c7c9f4edb503" {
		panic("Invalid checkpoint 7")
	}

	mainnetPrev8, _ := chainhash.NewHashFromStr("0000000000000000074f9edbfc07648dc74392ba8248f0983ffea63431b3bc20")
	mainnetMerk8, _ := chainhash.NewHashFromStr("0ed1b9a40f94aec95e2843369bdcabaa42f860c82391c54874a7c193d7268eaa")
	mainnetCheckpoints = append(mainnetCheckpoints, Checkpoint{
		Height: 399168,
		Header: wire.BlockHeader{
			Version:    4,
			PrevBlock:  *mainnetPrev8,
			MerkleRoot: *mainnetMerk8,
			Timestamp:  time.Unix(1455885256, 0),
			Bits:       403093919,
			Nonce:      3889666804,
		},
	})
	if mainnetCheckpoints[8].Header.BlockHash().String() != "000000000000000002ee8708b6c5ca0b4e055fbeed16189af256ed8a8d4b181d" {
		panic("Invalid checkpoint 8")
	}

	mainnetPrev9, _ := chainhash.NewHashFromStr("0000000000000000009f989a246ac4221ebdced8ccebae9b8d5c83b69bb5e7c8")
	mainnetMerk9, _ := chainhash.NewHashFromStr("97c6781a64bb3967f556263437252dba5baae9d46c3d0d6b6ab8073d6db4fff5")
	mainnetCheckpoints = append(mainnetCheckpoints, Checkpoint{
		Height: 455616,
		Header: wire.BlockHeader{
			Version:    536870914,
			PrevBlock:  *mainnetPrev9,
			MerkleRoot: *mainnetMerk9,
			Timestamp:  time.Unix(1488567886, 0),
			Bits:       402809567,
			Nonce:      27288148,
		},
	})
	if mainnetCheckpoints[9].Header.BlockHash().String() != "000000000000000001f87b34342c67eadb07a5c48a6127ba1628dd43c0b97b45" {
		panic("Invalid checkpoint 9")
	}

	mainnetPrev10, _ := chainhash.NewHashFromStr("0000000000000000010e0373719b7538e713e47d8d7189826dce4264d85a79b8")
	mainnetMerk10, _ := chainhash.NewHashFromStr("28e3ae14925aeead5b44a885bca88e06ec37483ae905e65d57aa8cd60cb74b2f")
	mainnetCheckpoints = append(mainnetCheckpoints, Checkpoint{
		Height: 477792,
		Header: wire.BlockHeader{
			Version:    536870914,
			PrevBlock:  *mainnetPrev10,
			MerkleRoot: *mainnetMerk10,
			Timestamp:  time.Unix(1501153434, 0),
			Bits:       402736949,
			Nonce:      165193412,
		},
	})
	if mainnetCheckpoints[10].Header.BlockHash().String() != "00000000000000000016ba7786309176445b838b36a16bd1ef3c3e3020473206" {
		panic("Invalid checkpoint 10")
	}

	// Testnet3
	BitcoinCashTestnet3ForkBlock, _ = chainhash.NewHashFromStr("00000000f17c850672894b9a75b63a1e72830bbd5f4c8889b5c1a80e7faef138")

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
		panic("Invalid checkpoint 0")
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

func ForkHeight(params *chaincfg.Params) uint32 {
	switch params.Name {
	case chaincfg.MainNetParams.Name:
		return 478559
	case chaincfg.TestNet3Params.Name:
		return 1155875
	}
	return 0
}

func IsForkBlock(params *chaincfg.Params, header wire.BlockHeader) bool {
	ckHash := header.BlockHash()
	switch params.Name {
	case chaincfg.MainNetParams.Name:
		return ckHash.IsEqual(BitcoinCashMainnetForkBlock)
	case chaincfg.TestNet3Params.Name:
		return ckHash.IsEqual(BitcoinCashTestnet3ForkBlock)
	}
	return false
}
