package pow

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	. "github.com/Lawliet-Chan/yu/blockchain"
	"github.com/Lawliet-Chan/yu/common"
	"github.com/sirupsen/logrus"
	"math"
	"math/big"
)

func Run(block IBlock, target *big.Int, targetBits int64) (nonce int64, hash common.Hash, err error) {
	var hashInt big.Int
	nonce = 0

	logrus.Info("[[[Mining a new Block!!!]]]")
	for nonce < math.MaxInt64 {
		var data []byte
		data, err = prepareData(block, nonce, targetBits)
		if err != nil {
			return
		}
		hash = sha256.Sum256(data)
		if math.Remainder(float64(nonce), 100000) == 0 {
			logrus.Infof("Hash is \r%x", hash.Bytes())
		}
		hashInt.SetBytes(hash.Bytes())

		if hashInt.Cmp(target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return
}

func Validate(block IBlock, target *big.Int, targetBits int64) bool {
	var hashInt big.Int

	var nonce uint64 = block.GetHeader().(*Header).Nonce
	data, err := prepareData(block, int64(nonce), targetBits)
	if err != nil {
		return false
	}
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(target) == -1
}

func prepareData(block IBlock, nonce, targetBits int64) ([]byte, error) {
	num := block.GetTimestamp()
	hex1, err := intToHex(int64(num))
	if err != nil {
		return nil, err
	}
	hex2, err := intToHex(targetBits)
	if err != nil {
		return nil, err
	}
	hex3, err := intToHex(nonce)
	if err != nil {
		return nil, err
	}
	data := bytes.Join(
		[][]byte{
			block.GetPrevHash().Bytes(),
			block.GetTxnRoot().Bytes(),
			hex1,
			hex2,
			hex3,
		},
		[]byte{},
	)

	return data, nil
}

func intToHex(num int64) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
