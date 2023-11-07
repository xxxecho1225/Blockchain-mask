package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"zkblockchain/utils"
	"math/big"
	"strconv"
	"github.com/btcsuite/btcd/btcutil/base58"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

func NewWallet() *Wallet {

	w := new(Wallet)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.privateKey = privateKey
	w.publicKey = &w.privateKey.PublicKey

	h := sha256.New()
	h.Write(w.publicKey.X.Bytes())
	h.Write(w.publicKey.Y.Bytes())
	digest := h.Sum(nil)
	address := base58.Encode(digest)
	w.blockchainAddress = address

	return w
}

// 为什么要写以下返回私钥和公钥的方法
func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {

	return w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) BlockchainAddress() string {
	return w.blockchainAddress
}


func LoadWallet(privkey string) *Wallet {
	privateKey := privkey
	privateKeyInt := new(big.Int)
	privateKeyInt.SetString(privateKey, 16)
	fmt.Println("privateKeyInt:", privateKeyInt)
	// 曲线
	curve := elliptic.P256()
	// 获取公钥
	x, y := curve.ScalarBaseMult(privateKeyInt.Bytes())
	publicKey := ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}
	wallet := new(Wallet)
	wallet.publicKey = &publicKey
	wallet.privateKey = &ecdsa.PrivateKey{
		PublicKey: publicKey,
		D:     privateKeyInt,
	}
	h := sha256.New()
	h.Write(wallet.publicKey .X.Bytes())
	h.Write(wallet.publicKey.Y.Bytes())
	digest := h.Sum(nil)
	fmt.Printf("digest: #{digest}\n")
	address := base58.Encode(digest)
	wallet.blockchainAddress = address
	return wallet
}

type Transaction struct {
	senderPrivateKey           *ecdsa.PrivateKey
	senderPublicKey            *ecdsa.PublicKey
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      *big.Int
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     *big.Int `json:"value"`
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}

func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey,
	sender string, recipient string, value *big.Int) *Transaction {
	return &Transaction{privateKey, publicKey, sender, recipient, value}
}

func (t *Transaction) GenerateSignature() *utils.Signature {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	return &utils.Signature{R: r, S: s}
}
func (wallet *Wallet) Transfer(recipientBlockchainAddress string, value float32) *Transaction {
	reward, _ := strconv.Atoi(fmt.Sprintf("%1.0f", value))
	return &Transaction{wallet.PrivateKey(), wallet.PublicKey(), wallet.BlockchainAddress(), recipientBlockchainAddress, big.NewInt(int64(reward))}
}