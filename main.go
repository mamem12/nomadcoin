package main

import (
	"github.com/nomadcoin/cli"
	"github.com/nomadcoin/db"
)

func main() {
	defer db.Close()
	cli.Start()

	// difficulty := 3
	// target := strings.Repeat("0", difficulty)
	// nonce := 1
	// fmt.Println(target)

	// for {
	// 	hash := fmt.Sprintf("%x\n", sha256.Sum256([]byte("hello"+fmt.Sprint(nonce))))
	// 	fmt.Printf("Hash:%s\nTarget:%s\nNonce:%d\n\n\n", hash, target, nonce)
	// 	if strings.HasPrefix(hash, target) {
	// 		break
	// 	} else {
	// 		nonce++
	// 	}
	// }

	// 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
}
