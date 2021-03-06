package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

// GetBlock returns a block at the specified height
func (ds DataStore) GetBlock(ctx sdk.Context, blockHeight *big.Int) (Block, bool) {
	key := GetBlockKey(blockHeight)
	data := ds.Get(ctx, key)
	if data == nil {
		return Block{}, false
	}

	block := Block{}
	if err := rlp.DecodeBytes(data, &block); err != nil {
		panic(fmt.Sprintf("block store corrupted: %s", err))
	}

	return block, true
}

// StoreBlock will store the plasma block and return the plasma block number
// in which it was stored at.
func (ds DataStore) StoreBlock(ctx sdk.Context, tmBlockHeight uint64, block plasma.Block) *big.Int {
	blockHeight := ds.NextPlasmaBlockHeight(ctx)

	blockKey := GetBlockKey(blockHeight)
	blockData, err := rlp.EncodeToBytes(&Block{block, tmBlockHeight})
	if err != nil {
		panic(fmt.Sprintf("error rlp encoding block: %s", err))
	}

	// store the block and updated the height counter
	ds.Set(ctx, blockKey, blockData)
	ds.Set(ctx, GetBlockHeightKey(), blockHeight.Bytes())

	return blockHeight
}

// PlasmaBlockHeight returns the current plasma block height. nil if no blocks exist
func (ds DataStore) PlasmaBlockHeight(ctx sdk.Context) *big.Int {
	var plasmaBlockNum *big.Int
	data := ds.Get(ctx, GetBlockHeightKey())
	if data == nil {
		return nil
	} else {
		plasmaBlockNum = new(big.Int).SetBytes(data)
	}

	return plasmaBlockNum
}

// NextPlasmaBlockHeight returns the next plasma block height
func (ds DataStore) NextPlasmaBlockHeight(ctx sdk.Context) *big.Int {
	height := ds.PlasmaBlockHeight(ctx)
	if height == nil {
		return big.NewInt(1)
	}

	return height.Add(height, utils.Big1)
}
