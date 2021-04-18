package master

import (
	"encoding/json"
	. "yu/common"
	. "yu/txn"
	"yu/yerror"
)

type HandShakeInfo struct {
	GenesisBlockHash Hash

	// When POW, these two will be always 0 and null.
	FinalizeHeight    BlockNum
	FinalizeBlockHash Hash

	EndHeight    BlockNum
	EndBlockHash Hash
}

func (m *Master) NewHsInfo() (*HandShakeInfo, error) {
	gBlock, err := m.chain.GetGenesis()
	if err != nil {
		return nil, err
	}
	fBlock, err := m.chain.GetFinalizedBlock()
	if err != nil {
		return nil, err
	}
	eBlock, err := m.chain.GetEndBlock()
	if err != nil {
		return nil, err
	}

	return &HandShakeInfo{
		GenesisBlockHash:  gBlock.GetHeader().GetHash(),
		FinalizeHeight:    fBlock.GetHeader().GetHeight(),
		FinalizeBlockHash: fBlock.GetHeader().GetHash(),
		EndHeight:         eBlock.GetHeader().GetHeight(),
		EndBlockHash:      eBlock.GetHeader().GetHash(),
	}, nil

}

func (hs *HandShakeInfo) Compare(other *HandShakeInfo) (*BlocksRange, error) {
	if hs.GenesisBlockHash != other.GenesisBlockHash {
		return nil, yerror.GenesisBlockIllegal
	}

	if hs.EndHeight < other.EndHeight || hs.FinalizeHeight < other.FinalizeHeight {
		if other.FinalizeHeight == 0 {
			return &BlocksRange{
				StartHeight: hs.EndHeight,
				EndHeight:   other.EndHeight,
			}, nil
		}

	}

	return nil, nil
}

func (hs *HandShakeInfo) Encode() ([]byte, error) {
	return json.Marshal(hs)
}

func DecodeHsInfo(data []byte) (*HandShakeInfo, error) {
	var hs HandShakeInfo
	err := json.Unmarshal(data, &hs)
	return &hs, err
}

type HandShakeResp struct {
	Br  *BlocksRange
	Err error
}

func (hs *HandShakeResp) Encode() ([]byte, error) {
	return json.Marshal(hs)
}

func DecodeHsResp(data []byte) (*HandShakeResp, error) {
	var hs HandShakeResp
	err := json.Unmarshal(data, &hs)
	return &hs, err
}

type BlocksRange struct {
	StartHeight BlockNum
	EndHeight   BlockNum
}

type PackedTxns struct {
	BlockHash string
	TxnsBytes []byte
}

func NewPackedTxns(blockHash Hash, txns SignedTxns) (*PackedTxns, error) {
	byt, err := txns.Encode()
	if err != nil {
		return nil, err
	}
	return &PackedTxns{
		BlockHash: blockHash.String(),
		TxnsBytes: byt,
	}, nil
}

func (pt *PackedTxns) Resolve() (Hash, SignedTxns, error) {
	txns := SignedTxns{}
	stxns, err := txns.Decode(pt.TxnsBytes)
	if err != nil {
		return NullHash, nil, err
	}
	return HexToHash(pt.BlockHash), stxns, nil
}
