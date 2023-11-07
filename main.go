package main

import (
	"fmt"
	"zkblockchain/block"
	"zkblockchain/wallet"
	"log"
	"math/big"
	"strconv"
	"github.com/fatih/color"
)

func init() {

	color.Blue("██╗  ██╗  █████╗ ██╗")
	color.Blue("██║ ██╔╝ ██╔══██╗██║")
	color.Blue("█████╔╝  ███████║██║")
	color.Blue("██╔═██╗  ██╔══██║██║")
	color.Blue("██║  ██╗ ██║  ██║██║")
	color.Blue("╚═╝  ╚═╝ ╚═╝  ╚═╝╚═╝")

	color.Red("██████╗ ██╗      ██████╗  ██████╗██╗  ██╗ ██████╗██╗  ██╗ █████╗ ██╗███╗   ██╗")
	color.Red("██╔══██╗██║     ██╔═══██╗██╔════╝██║ ██╔╝██╔════╝██║  ██║██╔══██╗██║████╗  ██║")
	color.Red("██████╔╝██║     ██║   ██║██║     █████╔╝ ██║     ███████║███████║██║██╔██╗ ██║")
	color.Red("██╔══██╗██║     ██║   ██║██║     ██╔═██╗ ██║     ██╔══██║██╔══██║██║██║╚██╗██║")
	color.Red("██████╔╝███████╗╚██████╔╝╚██████╗██║  ██╗╚██████╗██║  ██║██║  ██║██║██║ ╚████║")
	color.Red("╚═════╝ ╚══════╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝")

	log.SetPrefix("Blockchain: ")
}

func main() {
	// 1、根据矿工(姓名拼音zhoukai)私钥加载钱包，输出矿工区块链地址。(Loadwallet)
	wallet_zhoukai := wallet.LoadWallet(block.ZHOUKAI_ACCOUNT_ADDRESS) //矿工
	color.Yellow("矿工的account:%s\n", wallet_zhoukai.BlockchainAddress())
	
	// 2、生成2个账户account1、account2 的钱包，输出account1、account2区块链地址 (NewWallet)
	account1 := wallet.NewWallet()     //account1
	account2 := wallet.NewWallet()     //account2

	color.Yellow("account1:%s\n", account1.BlockchainAddress())
	color.Yellow("account2:%s\n", account2.BlockchainAddress())
	

	// 3、新建一条链
	blockchain := block.NewBlockchain(wallet_zhoukai.BlockchainAddress())
	fmt.Printf("account1[%s] %d\n", account1.BlockchainAddress(), blockchain.CalculateTotalAmount(account1.BlockchainAddress()))
	fmt.Printf("account2[%s] %d\n", account2.BlockchainAddress(), blockchain.CalculateTotalAmount(account2.BlockchainAddress()))
	fmt.Printf("矿工[%s]   %d\n", wallet_zhoukai.BlockchainAddress(), blockchain.CalculateTotalAmount(wallet_zhoukai.BlockchainAddress()))

	// 4、转账交易 矿工->account1 数量2e+19 (20000000000000000000)
	reward, _ := strconv.Atoi(fmt.Sprintf("%1.0f", 2e+19))
	fmt.Printf("t: %v\n", reward)
	t := wallet_zhoukai.Transfer(account1.BlockchainAddress(), 2e+19)
	isAdded := blockchain.AddTransaction(
		wallet_zhoukai.BlockchainAddress(),
		account1.BlockchainAddress(),
		big.NewInt(int64(reward)),
		wallet_zhoukai.PublicKey(),
		t.GenerateSignature())

	color.HiGreen("这笔交易验证通过吗? %v\n", isAdded)

	// 5、转账交易account1->account2 数量2000
	t2 := account1.Transfer(account2.BlockchainAddress(), 2000)

	//区块链 打包交易
	isAdded = blockchain.AddTransaction(
		account1.BlockchainAddress(),
		account2.BlockchainAddress(),
		big.NewInt(2000),
		account1.PublicKey(),
		t2.GenerateSignature())

	color.Blue("这笔交易验证通过吗? %v\n", isAdded)

	// 6、打包区块 Mining
	blockchain.Mining()

	// 7、查询区块号 GetBlockByNumber
	blockchain.GetBlockByNumber(1)

	// 8、查询区块hash GetBlockByHash
	hash := blockchain.LastBlock().Hash()
	blockchain.GetBlockByHash([]byte(hash[:]))
	// 9、查询步骤4的交易GetTransactionByHash
	blockchain.GetTransactionByHash(nil)
	// 10、输出区块信息(区块头和该区块的交易) Print
	blockchain.Print()

}
