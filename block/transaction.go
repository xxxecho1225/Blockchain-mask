package block

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"zkblockchain/utils"
	"log"
	"math/big"
	"strings"

	"github.com/fatih/color"
)

type Transaction struct {
	senderAddress  string
	receiveAddress string
	value          *big.Int
}

func NewTransaction(sender string, receive string, value *big.Int) *Transaction {
	t := Transaction{sender, receive, value}
	return &t
}

func (bc *Blockchain) AddTransaction(
	senderBlockchainAddress  string,
	recipientBlockchainAddress string,
	value *big.Int,
	senderPublicKey *ecdsa.PublicKey,
	s *utils.Signature) bool {
	t := NewTransaction(senderBlockchainAddress, recipientBlockchainAddress, value)

	//如果是挖矿得到的奖励交易，不验证
	if senderBlockchainAddress == ZHOUKAI_ACCOUNT_ADDRESS {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	// 判断有没有足够的余额
	if bc.CalculateTotalAmount(senderBlockchainAddress).Cmp(value) == -1{
		log.Printf("ERROR: %s ，你的钱包里没有足够的钱", senderBlockchainAddress)
		return false
	}

	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	} else {
		log.Println("ERROR: Verify Transaction")
	}
	return false

}

func (bc *Blockchain) VerifyTransactionSignature(
	senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

func (t *Transaction) Print() {
	color.Red("%s\n", strings.Repeat("~", 30))
	color.Cyan("发送地址             %s\n", t.senderAddress)
	color.Cyan("接受地址             %s\n", t.receiveAddress)
	color.Cyan("金额                 %d\n", t.value)

}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string `json:"sender_blockchain_address"`
		Recipient string `json:"recipient_blockchain_address"`
		Value     *big.Int  `json:"value"`
	}{
		Sender:    t.senderAddress,
		Recipient: t.receiveAddress,
		Value:     t.value,
	})
}
func (bc *Blockchain) GetTransactionByHash(hash []byte) *Transaction{
	var ha [32]byte
	copy(ha[:],hash)
	for i,block := range bc.block{
		if block.hash == ha{
			log.Printf("%-15v:%30d\n", "该交易所属区块为", i)
			return nil
		}
	}
	log.Printf("%-15v\n", "交易不存在")
	return nil
}
