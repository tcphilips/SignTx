package main

import (
	"C"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
)

type txdata struct {
	AccountNonce uint64          `json:"nonce"    gencodec:"required"`
	Price        *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit     *big.Int        `json:"gas"      gencodec:"required"`
	Recipient    *common.Address `json:"to"       rlp:"nil"` // nil means contract creation
	Amount       *big.Int        `json:"value"    gencodec:"required"`
	Payload      []byte          `json:"input"    gencodec:"required"`

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	// Hash *common.Hash `json:"hash" rlp:"-"`

	// Extra
	ChainId *big.Int `json:"chainId" rlp:"-"`
}

func (data txdata) Hash() []byte {
	return rlpHash([]interface{}{
		data.AccountNonce,
		data.Price,
		data.GasLimit,
		data.Recipient,
		data.Amount,
		data.Payload,
		data.ChainId, uint(0), uint(0),
	}).Bytes()
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func sign(hash []byte, prv *btcec.PrivateKey) ([]byte, error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}
	if prv.Curve != btcec.S256() {
		return nil, fmt.Errorf("private key curve is not secp256k1")
	}
	sig, err := btcec.SignCompact(btcec.S256(), prv, hash, false)
	if err != nil {
		return nil, err
	}
	// Convert to Ethereum signature format with 'recovery id' v at the end.
	v := sig[0] - 27
	copy(sig, sig[1:])
	sig[64] = v
	return sig, nil
}

//export SignTxWithPrivKey
func SignTxWithPrivKey(txJson, privKey *C.char) *C.char {
	rawTx, err := signTxWithPrivKey(C.GoString(txJson), C.GoString(privKey), big.NewInt(0))
	if err == nil {
		return C.CString(rawTx)
	} else {
		return C.CString("ERROR: " + err.Error())
	}
}

func signTxWithPrivKey(txJson, privKey string, maxValue *big.Int) (string, error) {
	var data txdata
	err := json.Unmarshal([]byte(txJson), &data)
	if err != nil {
		return "", err
	}

	if maxValue.Cmp(big.NewInt(0)) == 1 && maxValue.Cmp(data.Amount) == -1  {
		return "", fmt.Errorf("exceeds max value")
	}

	prv, _ := btcec.PrivKeyFromBytes(btcec.S256(), common.FromHex(privKey))

	sig, err := sign(data.Hash(), prv)
	if err != nil {
		return "", err
	}

	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])
	v := new(big.Int).SetBytes([]byte{sig[64] + 27})

	if data.ChainId.Sign() != 0 {
		v = big.NewInt(int64(sig[64] + 35))
		v.Add(v, new(big.Int).Mul(data.ChainId, big.NewInt(2)))
	}

	data.R, data.S, data.V = r, s, v

	raw, err := rlp.EncodeToBytes(data)
	if err != nil {
		return "", err
	}
	return common.ToHex(raw), nil
}

func main() {}
