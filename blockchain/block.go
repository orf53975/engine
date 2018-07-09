package blockchain

import (
	"encoding/json"
	"time"

	"bytes"

	"errors"

	"github.com/it-chain/yggdrasill/common"
	"fmt"
)

var ErrDecodingEmptyBlock = errors.New("Empty Block decoding failed")
var ErrTransactionType = errors.New("Wrong transaction type")

type Block = common.Block

type BlockHeight = uint64

type DefaultBlock struct {
	Seal      []byte
	PrevSeal  []byte
	Height    uint64
	TxList    []Transaction
	TxSeal    [][]byte
	Timestamp time.Time
	Creator   []byte
}

func (block *DefaultBlock) SetSeal(seal []byte) {
	block.Seal = seal
}

func (block *DefaultBlock) SetPrevSeal(prevSeal []byte) {
	block.PrevSeal = prevSeal
}

func (block *DefaultBlock) SetHeight(height uint64) {
	block.Height = height
}

func (block *DefaultBlock) PutTx(transaction Transaction) error {

	if block.TxList == nil {
		block.TxList = make([]Transaction, 0)
	}

	block.TxList = append(block.TxList, transaction)

	return nil

	return ErrTransactionType
}

func (block *DefaultBlock) SetTxSeal(txSeal [][]byte) {
	block.TxSeal = txSeal
}

func (block *DefaultBlock) SetCreator(creator []byte) {
	block.Creator = creator
}

func (block *DefaultBlock) SetTimestamp(currentTime time.Time) {
	block.Timestamp = currentTime
}

func (block *DefaultBlock) GetSeal() []byte {
	return block.Seal
}

func (block *DefaultBlock) GetPrevSeal() []byte {
	return block.PrevSeal
}

func (block *DefaultBlock) GetHeight() uint64 {
	return block.Height
}

func (block *DefaultBlock) GetTxList() []Transaction {
	txList := make([]Transaction, 0)
	for _, tx := range block.TxList {
		txList = append(txList, tx)
	}
	return txList
}

func (block *DefaultBlock) GetTxSeal() [][]byte {
	return block.TxSeal
}

func (block *DefaultBlock) GetCreator() []byte {
	return block.Creator
}

func (block *DefaultBlock) GetTimestamp() time.Time {
	return block.Timestamp
}

func (block *DefaultBlock) Serialize() ([]byte, error) {
	data, err := json.Marshal(block)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (block *DefaultBlock) Deserialize(serializedBlock []byte) error {
	if len(serializedBlock) == 0 {
		return ErrDecodingEmptyBlock
	}

	err := json.Unmarshal(serializedBlock, block)
	if err != nil {
		return err
	}

	return nil
}

func (block *DefaultBlock) IsReadyToPublish() bool {
	return block.Seal != nil
}

func (block *DefaultBlock) IsPrev(serializedPrevBlock []byte) bool {
	prevBlock := &DefaultBlock{}
	prevBlock.Deserialize(serializedPrevBlock)

	return bytes.Compare(prevBlock.GetSeal(), block.GetPrevSeal()) == 0
}

func NewEmptyBlock(prevSeal []byte, height uint64, creator []byte) *DefaultBlock {
	block := &DefaultBlock{}

	block.SetPrevSeal(prevSeal)
	block.SetHeight(height)
	block.SetCreator(creator)

	return block
}

// interface of api gateway query api
type BlockQueryApi interface {
	GetBlockByHeight(blockHeight uint64) (Block, error)
	GetBlockBySeal(seal []byte) (Block, error)
	GetBlockByTxID(txid string) (Block, error)
	GetLastBlock() (Block, error)
	GetTransactionByTxID(txid string) (Transaction, error)
}


type Action interface {
	DoAction(block Block)
}

func CreateSaveOrSyncAction(checkResult int64, pool BlockPool) Action {
	if checkResult > 0 {
		return NewSyncAction()
	} else if checkResult == 0 {
		return NewSaveAction(pool)
	} else {
		return NewDefaultAction()
	}
}

type SyncAction struct {}

func NewSyncAction() *SyncAction {
	return &SyncAction{}
}

func (syncAction *SyncAction) DoAction(block Block) {
	// TODO: Start synchronize
}


type SaveAction struct {
	blockPool BlockPool
}

func NewSaveAction(blockPool BlockPool) *SaveAction {
	return &SaveAction{
		blockPool: blockPool,
	}
}

func (saveAction *SaveAction) DoAction(block Block) {
	saveAction.blockPool.Add(block)
}

type DefaultAction struct {}

func NewDefaultAction() *DefaultAction{
	return &DefaultAction{}
}

func (defaultAction *DefaultAction) DoAction(block Block) {
	fmt.Printf("got shorter height block [%v]", block.GetHeight())
}