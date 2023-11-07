package wallet

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func PrintEthereumAccount() {
	privateKey, _ := crypto.HexToECDSA("396b5b77f7ad32e3b7013e64b8bfd919d7cfd27c0d54f1b0af636ab8d8efe3e6")
	publicKey := crypto.FromECDSAPub(&privateKey.PublicKey)
	fmt.Printf("以太坊publicKey:%x\n", publicKey)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Println("以太坊Address: ", address.Hex())
}
