package blockchain

import (
	. "github.com/Lawliet-Chan/yu/common"
	"github.com/Lawliet-Chan/yu/keypair"
	. "github.com/Lawliet-Chan/yu/result"
	. "github.com/Lawliet-Chan/yu/txn"
	"gorm.io/gorm"
)

type TxnScheme struct {
	TxnHash   string `gorm:"primaryKey"`
	Pubkey    string
	Signature string
	RawTxn    string

	BlockHash string
}

func (TxnScheme) TableName() string {
	return "txns"
}

func newTxnScheme(blockHash Hash, stxn *SignedTxn) (TxnScheme, error) {
	txnSm, err := toTxnScheme(stxn)
	if err != nil {
		return TxnScheme{}, err
	}
	txnSm.BlockHash = blockHash.String()
	return txnSm, nil
}

func toTxnScheme(stxn *SignedTxn) (TxnScheme, error) {
	rawTxnByt, err := stxn.GetRaw().Encode()
	if err != nil {
		return TxnScheme{}, err
	}
	return TxnScheme{
		TxnHash:   stxn.GetTxnHash().String(),
		Pubkey:    stxn.GetPubkey().StringWithType(),
		Signature: ToHex(stxn.GetSignature()),
		RawTxn:    ToHex(rawTxnByt),
		BlockHash: "",
	}, nil
}

func (t TxnScheme) toTxn() (*SignedTxn, error) {
	ut := &UnsignedTxn{}
	rawTxn, err := ut.Decode(FromHex(t.RawTxn))
	if err != nil {
		return nil, err
	}
	pubkey, err := keypair.PubkeyFromStr(t.Pubkey)
	if err != nil {
		return nil, err
	}
	return &SignedTxn{
		Raw:       rawTxn,
		TxnHash:   HexToHash(t.TxnHash),
		Pubkey:    pubkey,
		Signature: FromHex(t.Signature),
	}, nil
}

type EventScheme struct {
	gorm.Model
	Caller     string
	BlockStage string
	BlockHash  string
	Height     BlockNum
	TripodName string
	ExecName   string
	Value      string
}

func (EventScheme) TableName() string {
	return "events"
}

func toEventScheme(event *Event) (EventScheme, error) {
	return EventScheme{
		Caller:     event.Caller.String(),
		BlockStage: event.BlockStage,
		BlockHash:  event.BlockHash.String(),
		Height:     event.Height,
		TripodName: event.TripodName,
		ExecName:   event.ExecName,
		Value:      event.Value,
	}, nil
}

func (e EventScheme) toEvent() (*Event, error) {
	return &Event{
		Caller:     HexToAddress(e.Caller),
		BlockStage: e.BlockStage,
		BlockHash:  HexToHash(e.BlockHash),
		Height:     e.Height,
		TripodName: e.TripodName,
		ExecName:   e.ExecName,
		Value:      e.Value,
	}, nil

}

type ErrorScheme struct {
	gorm.Model
	Caller     string
	BlockStage string
	BlockHash  string
	Height     BlockNum
	TripodName string
	ExecName   string
	Error      string
}

func (ErrorScheme) TableName() string {
	return "errors"
}

func toErrorScheme(err *Error) ErrorScheme {
	return ErrorScheme{
		Caller:     err.Caller.String(),
		BlockStage: err.BlockStage,
		BlockHash:  err.BlockHash.String(),
		Height:     err.Height,
		TripodName: err.TripodName,
		ExecName:   err.ExecName,
		Error:      err.Err,
	}
}

func (e ErrorScheme) toError() *Error {
	return &Error{
		Caller:     HexToAddress(e.Caller),
		BlockStage: e.BlockStage,
		BlockHash:  HexToHash(e.BlockHash),
		Height:     e.Height,
		TripodName: e.TripodName,
		ExecName:   e.ExecName,
		Err:        e.Error,
	}
}
