package block

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"strconv"
	"time"
	"github.com/fatih/color"
)

const MINING_DIFFICULT = 3
const ZHOUKAI_ACCOUNT_ADDRESS = "7a6026c4bc5f70169cfc16f60cb9283a011a4e452a67151164c992fe50ac3fc4"
const MINING_REWARD = 1e+20

type Block struct {
	nonce        *big.Int
	timestamp    uint64
	number       *big.Int
	previousHash [32]byte
	tosize       uint16
	difficulty   *big.Int
	hash         [32]byte
	transactions []*Transaction
}

func NewBlock(nonce *big.Int, number *big.Int,previousHash [32]byte,txs []*Transaction) *Block {
	b := new(Block)
	b.timestamp = uint64(time.Now().UnixNano())
	b.nonce = nonce
	b.number = number
	b.difficulty = big.NewInt(2)
	b.previousHash = previousHash
	b.transactions = txs
	b.hash = b.Hash()
	return b
}


func (b *Block) Print() {
	log.Printf("%-15v:%30d\n", "nonce", b.nonce)
	log.Printf("%-15v:%30d\n", "number", b.number)
	log.Printf("%-15v:%30d\n", "txSize", b.tosize)
	log.Printf("%-15v:%30x\n", "previous_hash", b.previousHash)
	log.Printf("%-15v:%30x\n", "hash", b.hash)
	for _, t := range b.transactions {
		t.Print()
	}
}

type Blockchain struct {
	transactionPool   []*Transaction //交易池
	block             []*Block
	coinbase string   //区块奖励地址
 }

// 新建一条链的第一个区块
// NewBlockchain(blockchainAddress string) *Blockchain
// 函数定义了一个创建区块链的方法，它接收一个字符串类型的参数 blockchainAddress，
// 它返回一个区块链类型的指针。在函数内部，它创建一个区块链对象并为其设置地址，
// 然后创建一个创世块并将其添加到区块链中，最后返回区块链对象。
func NewBlockchain(coinbase string) *Blockchain {
	bc := new(Blockchain)
	b := &Block{}
	bc.CreateBlock(big.NewInt(0),big.NewInt(0),b.Hash())
	bc.coinbase = coinbase
	return bc
}

// (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block
//  函数是在区块链上创建新的区块，它接收两个参数：一个int类型的nonce和一个字节数组类型的 previousHash，
//  返回一个区块类型的指针。在函数内部，它使用传入的参数来创建一个新的区块，
//  然后将该区块添加到区块链的链上，并清空交易池。

func (bc *Blockchain) CreateBlock(nonce *big.Int, number *big.Int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, number,previousHash,bc.transactionPool)
	bc.block = append(bc.block, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) Print() {
	for i, block := range bc.block {
		color.Green("%s BLOCK %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	color.Yellow("%s\n\n\n", strings.Repeat("*", 50))
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    uint64          `json:"timestamp"`
		Nonce        *big.Int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.block[len(bc.block)-1]
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(t.senderAddress,
				t.receiveAddress,
				t.value))
	}
	return transactions
}

func (bc *Blockchain) ValidProof(nonce *big.Int,
	previousHash [32]byte,
	transactions []*Transaction,
	difficulty *big.Int,
) bool {
	zeros := strings.Repeat("0", int(difficulty.Int64()))
	//zeros := "1234"
	//tmpBlock := Block{nonce: nonce, previousHash: previousHash, transactions: transactions, timestamp: 0}
	tmpBlock := Block{nonce: nonce, previousHash: previousHash, transactions: transactions, timestamp: uint64(time.Now().UnixNano())}
	//log.Printf("tmpBlock%+v", tmpBlock)
	tmpHashStr := fmt.Sprintf("%x", tmpBlock.Hash())
	//log.Println("guessHashStr", tmpHashStr)
	return tmpHashStr[:int(difficulty.Int64())] == zeros
}
func (bc *Blockchain) ProofOfWork() *big.Int {
	transactions := bc.CopyTransactionPool() //选择交易？控制交易数量？
	previousHash := bc.LastBlock().Hash()
	nonce := big.NewInt(0)
	begin := time.Now()
	for !bc.ValidProof(nonce, previousHash, transactions,big.NewInt(MINING_DIFFICULT)) {
		nonce = nonce.Add(nonce, big.NewInt(1))
	}

	end := time.Now()
	//log.Printf("POW spend Time:%f", end.Sub(begin).Seconds())
	log.Printf("POW spend Time:%f Second", end.Sub(begin).Seconds())
	log.Printf("POW spend Time:%s", end.Sub(begin))

	return nonce
}

func (bc *Blockchain) Mining() bool {
	reward, _ := strconv.Atoi(fmt.Sprintf("%1.0f", MINING_REWARD))
	bc.AddTransaction(ZHOUKAI_ACCOUNT_ADDRESS, bc.coinbase, big.NewInt(int64(reward)), nil, nil) //因为是挖矿奖励，不用公钥和签名
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	var num = bc.LastBlock().number.Int64()
	bc.CreateBlock(nonce, big.NewInt(num).Add(big.NewInt(num), big.NewInt(1)), previousHash)
	color.Red("action=mining, status=success")
	return true
}

func (bc *Blockchain) CalculateTotalAmount(accountAddress string) *big.Int {
	var totalAmount = big.NewInt(0)
	for _, _chain := range bc.block {
		for _, _tx := range _chain.transactions {
			if accountAddress == _tx.receiveAddress {
				totalAmount = totalAmount.Add(totalAmount,_tx.value)
			}
			if accountAddress == _tx.senderAddress {
				totalAmount = totalAmount.Add(totalAmount,_tx.value)
			}
		}
	}
	return totalAmount
}
// 根据区块ID输出该结构体内容
func (blockchain *Blockchain) GetBlockByNumber(blockid uint64) (*Block, error) {
	for i, block := range blockchain.block {
		if big.NewInt(int64(i)).Cmp(big.NewInt(int64(blockid))) == 0 {
			log.Printf("%-15v:%30d\n", "nonce", block.nonce)
			log.Printf("%-15v:%30d\n", "timestamp", block.timestamp)
			log.Printf("%-15v:%30d\n", "number", block.number)
			log.Printf("%-15v:%30d\n", "difficulty", block.difficulty)
			log.Printf("%-15v:%30x\n", "previousHash", block.previousHash)
			log.Printf("%-15v:%30x\n", "hash", block.hash)
			log.Printf("%-15v:%30d\n", "tosize", block.tosize)
			return block, nil
		}
	}
	log.Printf("%-15v:%30s\n", "error", "没找到对应区块ID结构体内容")
	return nil, errors.New("没找到对应区块信息")
}
// 根据区块哈希输出该结构体内容
func (blockchain *Blockchain) GetBlockByHash(hash []byte) (*Block, error) {
	var ha [32]byte
	copy(ha[:], hash)
	for _, block := range blockchain.block {
		if ha == block.hash {
			log.Printf("%-15v:%30d\n", "nonce", block.nonce)
			log.Printf("%-15v:%30d\n", "timestamp", block.timestamp)
			log.Printf("%-15v:%30d\n", "number", block.number)
			log.Printf("%-15v:%30d\n", "difficulty", block.difficulty)
			log.Printf("%-15v:%30x\n", "previousHash", block.previousHash)
			log.Printf("%-15v:%30x\n", "hash", block.hash)
			log.Printf("%-15v:%30d\n", "tosize", block.tosize)
			return nil, nil
		}
	}
	log.Printf("%-15v:%30s\n", "error", "没找到对应区块哈希的结构体内容")
	return nil, nil
}
