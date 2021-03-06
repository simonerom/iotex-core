// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package blockindex

import (
	"context"
	"sync"

	"github.com/iotexproject/go-pkgs/bloom"

	"github.com/iotexproject/iotex-core/action"
	filter "github.com/iotexproject/iotex-core/api/logfilter"
	"github.com/iotexproject/iotex-core/blockchain/block"
	"github.com/iotexproject/iotex-core/blockchain/blockdao"
	"github.com/iotexproject/iotex-core/config"
	"github.com/iotexproject/iotex-core/db"
	"github.com/iotexproject/iotex-core/db/batch"
	"github.com/iotexproject/iotex-core/pkg/util/byteutil"
	"github.com/pkg/errors"
)

const (
	// BlockBloomFilterNamespace indicated the kvstore namespace to store BlockBloomFilterNamespace
	BlockBloomFilterNamespace = "BlockBloomFilters"
	// RangeBloomFilterNamespace indicates the kvstore namespace to store RangeBloomFilters
	RangeBloomFilterNamespace = "RangeBloomFilters"
	// CurrentHeightKey indicates the key of current bf indexer height in underlying DB
	CurrentHeightKey = "CurrentHeight"
)

type (
	// BloomFilterIndexer is the interface for bloomfilter indexer
	BloomFilterIndexer interface {
		blockdao.BlockIndexer
		// RangeBloomFilterBlocks returns the number of blocks that each rangeBloomfilter includes
		RangeBloomFilterBlocks() uint64
		// BlockFilterByHeight returns the block-level bloomfilter which includes not only topic but also address of logs info by given block height
		BlockFilterByHeight(uint64) (bloom.BloomFilter, error)
		// RangeFilterByHeight returns the range bloomfilter for the height
		RangeFilterByHeight(uint64) (bloom.BloomFilter, error)
		// FilterBlocksInRange returns the block numbers by given logFilter in range from start to end
		FilterBlocksInRange(*filter.LogFilter, uint64, uint64) ([]uint64, error)
	}

	// bloomfilterIndexer is a struct for bloomfilter indexer
	bloomfilterIndexer struct {
		mutex               sync.RWMutex // mutex for curRangeBloomfilter
		flusher             db.KVStoreFlusher
		rangeSize           uint64
		bfSize              uint64
		bfNumHash           uint64
		curRangeBloomfilter bloom.BloomFilter
	}
)

// NewBloomfilterIndexer creates a new bloomfilterindexer struct by given kvstore and rangebloomfilter size
func NewBloomfilterIndexer(kv db.KVStore, cfg config.Indexer) (BloomFilterIndexer, error) {
	if kv == nil {
		return nil, errors.New("empty kvStore")
	}
	flusher, err := db.NewKVStoreFlusher(kv, batch.NewCachedBatch())
	if err != nil {
		return nil, err
	}
	return &bloomfilterIndexer{
		flusher:   flusher,
		rangeSize: cfg.RangeBloomFilterBlocks,
		bfSize:    cfg.RangeBloomFilterSize,
		bfNumHash: cfg.RangeBloomFilterNumHash,
	}, nil
}

// Start starts the bloomfilter indexer
func (bfx *bloomfilterIndexer) Start(ctx context.Context) error {
	if err := bfx.flusher.KVStoreWithBuffer().Start(ctx); err != nil {
		return err
	}
	bfx.mutex.Lock()
	defer bfx.mutex.Unlock()
	tipHeightData, err := bfx.flusher.KVStoreWithBuffer().Get(RangeBloomFilterNamespace, []byte(CurrentHeightKey))
	switch errors.Cause(err) {
	case nil:
		tipHeight := byteutil.BytesToUint64(tipHeightData)
		if tipHeight%bfx.rangeSize == 0 {
			bfx.curRangeBloomfilter, _ = bloom.NewBloomFilter(bfx.bfSize, bfx.bfNumHash)
		} else {
			bfx.curRangeBloomfilter, err = bfx.rangeBloomFilter(tipHeight)
			if err != nil {
				return errors.Wrapf(err, "failed to read curRangeBloomfilter from DB")
			}
		}
	case db.ErrNotExist:
		if err = bfx.flusher.KVStoreWithBuffer().Put(RangeBloomFilterNamespace, []byte(CurrentHeightKey), byteutil.Uint64ToBytes(0)); err != nil {
			return err
		}
		if err := bfx.flusher.Flush(); err != nil {
			return errors.Wrapf(err, "failed to flush")
		}
		bfx.curRangeBloomfilter, _ = bloom.NewBloomFilter(bfx.bfSize, bfx.bfNumHash)
	default:
		return err
	}
	return nil
}

// Stop stops the bloomfilter indexer
func (bfx *bloomfilterIndexer) Stop(ctx context.Context) error {
	return bfx.flusher.KVStoreWithBuffer().Stop(ctx)
}

// Height returns the tipHeight from underlying DB
func (bfx *bloomfilterIndexer) Height() (uint64, error) {
	h, err := bfx.flusher.KVStoreWithBuffer().Get(RangeBloomFilterNamespace, []byte(CurrentHeightKey))
	if err != nil {
		return 0, err
	}
	return byteutil.BytesToUint64(h), nil
}

// PutBlock processes new block by adding logs into rangebloomfilter, and if necessary, updating underlying DB
func (bfx *bloomfilterIndexer) PutBlock(ctx context.Context, blk *block.Block) (err error) {
	bfx.mutex.Lock()
	defer bfx.mutex.Unlock()
	bfx.addLogsToRangeBloomFilter(ctx, blk.Height(), blk.Receipts)
	// commit into DB and update tipHeight
	if err := bfx.commit(blk.Height(), bfx.calculateBlockBloomFilter(ctx, blk.Receipts)); err != nil {
		return err
	}
	if blk.Height()%bfx.rangeSize == 0 {
		bfx.curRangeBloomfilter, err = bloom.NewBloomFilter(bfx.bfSize, bfx.bfNumHash)
		if err != nil {
			return errors.Wrapf(err, "Can not create new bloomfilter")
		}
	}
	return nil
}

// DeleteTipBlock deletes tip height from underlying DB if necessary
func (bfx *bloomfilterIndexer) DeleteTipBlock(blk *block.Block) (err error) {
	bfx.mutex.Lock()
	defer bfx.mutex.Unlock()
	height := blk.Height()
	if err := bfx.delete(height); err != nil {
		return err
	}
	bfx.curRangeBloomfilter = nil
	return nil
}

// RangeBloomFilterBlocks returns the number of blocks that each rangeBloomfilter includes
func (bfx *bloomfilterIndexer) RangeBloomFilterBlocks() uint64 {
	bfx.mutex.RLock()
	defer bfx.mutex.RUnlock()
	return bfx.rangeSize
}

// BlockFilterByHeight returns the block-level bloomfilter which includes not only topic but also address of logs info by given block height
func (bfx *bloomfilterIndexer) BlockFilterByHeight(height uint64) (bloom.BloomFilter, error) {
	bfBytes, err := bfx.flusher.KVStoreWithBuffer().Get(BlockBloomFilterNamespace, byteutil.Uint64ToBytes(height))
	if err != nil {
		return nil, err
	}
	return bloom.BloomFilterFromBytes(bfBytes)
}

// RangeFilterByHeight returns the range bloomfilter for the height
func (bfx *bloomfilterIndexer) RangeFilterByHeight(height uint64) (bloom.BloomFilter, error) {
	return bfx.rangeBloomFilter(height)
}

// FilterBlocksInRange returns the block numbers by given logFilter in range [start, end]
func (bfx *bloomfilterIndexer) FilterBlocksInRange(l *filter.LogFilter, start, end uint64) ([]uint64, error) {
	bfx.mutex.RLock()
	defer bfx.mutex.RUnlock()
	if start == 0 || end == 0 {
		return nil, errors.New("start/end height should be bigger than zero")
	}
	blockNumbers := make([]uint64, 0)
	queryHeight := bfx.rangeBloomfilterKey(start)  // range which includes start
	endQueryHeight := bfx.rangeBloomfilterKey(end) // range which includes end
	for queryHeight <= endQueryHeight {
		bigBloom, err := bfx.rangeBloomFilter(queryHeight)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get rangeBloomFilter from indexer by given height %d", queryHeight)
		}
		if l.ExistInBloomFilterv2(bigBloom) {
			blocks := l.SelectBlocksFromRangeBloomFilter(bigBloom, queryHeight-bfx.rangeSize+1, queryHeight)
			for _, num := range blocks {
				if num >= start && num <= end {
					blockNumbers = append(blockNumbers, num)
				}
			}
		}
		queryHeight += bfx.rangeSize
	}

	return blockNumbers, nil
}

func (bfx *bloomfilterIndexer) rangeBloomfilterKey(blockNumber uint64) uint64 {
	numRange := (blockNumber + bfx.rangeSize - 1) / bfx.rangeSize
	return bfx.rangeSize * numRange
}

// rangeBloomFilter reads rangebloomfilter by given block number from underlying DB
func (bfx *bloomfilterIndexer) rangeBloomFilter(blockNumber uint64) (bloom.BloomFilter, error) {
	rangeBloomfilterKey := bfx.rangeBloomfilterKey(blockNumber)
	bfBytes, err := bfx.flusher.KVStoreWithBuffer().Get(RangeBloomFilterNamespace, byteutil.Uint64ToBytes(rangeBloomfilterKey))
	if err != nil {
		return nil, err
	}
	return bloom.BloomFilterFromBytes(bfBytes)
}

func (bfx *bloomfilterIndexer) delete(blockNumber uint64) error {
	rangeBloomfilterKey := bfx.rangeBloomfilterKey(blockNumber)
	if err := bfx.flusher.KVStoreWithBuffer().Delete(RangeBloomFilterNamespace, byteutil.Uint64ToBytes(rangeBloomfilterKey)); err != nil {
		return err
	}
	if err := bfx.flusher.KVStoreWithBuffer().Delete(BlockBloomFilterNamespace, byteutil.Uint64ToBytes(blockNumber)); err != nil {
		return err
	}
	if err := bfx.flusher.KVStoreWithBuffer().Put(RangeBloomFilterNamespace, []byte(CurrentHeightKey), byteutil.Uint64ToBytes(rangeBloomfilterKey-bfx.rangeSize)); err != nil {
		return err
	}

	return bfx.flusher.Flush()
}

func (bfx *bloomfilterIndexer) commit(blockNumber uint64, blkBloomfilter bloom.BloomFilter) error {
	rangeBloomfilterKey := bfx.rangeBloomfilterKey(blockNumber)
	if err := bfx.flusher.KVStoreWithBuffer().Put(RangeBloomFilterNamespace, byteutil.Uint64ToBytes(rangeBloomfilterKey), bfx.curRangeBloomfilter.Bytes()); err != nil {
		return err
	}
	if err := bfx.flusher.KVStoreWithBuffer().Put(BlockBloomFilterNamespace, byteutil.Uint64ToBytes(blockNumber), blkBloomfilter.Bytes()); err != nil {
		return err
	}
	if err := bfx.flusher.KVStoreWithBuffer().Put(RangeBloomFilterNamespace, []byte(CurrentHeightKey), byteutil.Uint64ToBytes(blockNumber)); err != nil {
		return err
	}

	return bfx.flusher.Flush()
}

func (bfx *bloomfilterIndexer) calculateBlockBloomFilter(ctx context.Context, receipts []*action.Receipt) bloom.BloomFilter {
	bloom, _ := bloom.NewBloomFilter(2048, 3)
	for _, receipt := range receipts {
		for _, l := range receipt.Logs() {
			bloom.Add([]byte(l.Address))
			for i, topic := range l.Topics {
				bloom.Add(append(byteutil.Uint64ToBytes(uint64(i)), topic[:]...)) //position-sensitive
			}
		}
	}
	return bloom
}

func (bfx *bloomfilterIndexer) addLogsToRangeBloomFilter(ctx context.Context, blockNumber uint64, receipts []*action.Receipt) {
	Heightkey := append([]byte(filter.BlockHeightPrefix), byteutil.Uint64ToBytes(blockNumber)...)

	for _, receipt := range receipts {
		for _, l := range receipt.Logs() {
			bfx.curRangeBloomfilter.Add([]byte(l.Address))
			bfx.curRangeBloomfilter.Add(append(Heightkey, []byte(l.Address)...)) // concatenate with block number
			for i, topic := range l.Topics {
				bfx.curRangeBloomfilter.Add(append(byteutil.Uint64ToBytes(uint64(i)), topic[:]...)) //position-sensitive
				bfx.curRangeBloomfilter.Add(append(Heightkey, topic[:]...))                         // concatenate with block number
			}
		}
	}
	return
}
