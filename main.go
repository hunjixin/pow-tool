package main

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"
	"runtime"
	"time"

	"github.com/holiman/uint256"
	"github.com/lilypad-tech/lilypad/pkg/resourceprovider"
)

func main() {
	ctx := context.Background()
	nodeId := "123"
	taskCh := make(chan resourceprovider.Task)

	submitWork := func(nonce *big.Int) {
	}
	miner := resourceprovider.NewCpuMiner(nodeId, runtime.NumCPU()*2, taskCh, submitWork)
	go miner.Start(ctx)

	challenge := [32]byte{}
	challengeBytes, err := hex.DecodeString("86372f9059c8c8abc6f3330779e33c97ca1e712cab6c9da8abff294cc2f218f8")
	if err != nil {
		log.Fatal(err)
		return
	}
	copy(challenge[:], challengeBytes)

	bigDifficulty, _ := new(big.Int).SetString("1999658518639694904537102086817924629401430693671659749396090098758451", 10)
	taskCh <- resourceprovider.Task{
		Challenge:  challenge,
		Difficulty: uint256.MustFromBig(bigDifficulty),
	}
	time.Sleep(time.Hour)
}
