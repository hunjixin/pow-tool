package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"runtime"

	"github.com/holiman/uint256"
	"github.com/lilypad-tech/lilypad/pkg/resourceprovider"
)

func main() {
	ctx := context.Background()
	nodeId := "mock node id"
	taskCh := make(chan resourceprovider.Task)

	resultCh := make(chan *big.Int)
	submitWork := func(nonce *big.Int) {
		resultCh <- nonce
	}
	miner := resourceprovider.NewCpuMiner(nodeId, runtime.NumCPU()*2, taskCh, submitWork)
	go miner.Start(ctx)

	challenge := [32]byte{}
	rand.Read(challenge[:])

	//maxDifficulty := new(uint256.Int).Sub(uint256.NewInt(0), uint256.NewInt(1))
	minDifficulty := new(uint256.Int)
	fmt.Println(minDifficulty.String())
	taskCh <- resourceprovider.Task{
		Challenge:  challenge,
		Difficulty: minDifficulty,
	}
	_ = <-resultCh // every result can pass this difficult
}
