package main

import (
	"fmt"
	"math/big"

	"github.com/ALTree/bigfloat"
)

var (
	uint256Max, _ = new(big.Float).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935")
)

// 40m 74061421902600066637965639403500996888198243478266023770329538625536
func main() {
	targetDifficulty := esitimateDifficulty(1000, 30, 0.00001)

	number := esitimateHashRate(targetDifficulty, 30*12, 0.00001)
	fmt.Print(number)
}

func esitimateDifficulty(power int, windowsSize int, probability float64) *big.Float {
	n := int64(power * 1000 * 1000 * windowsSize * 12)

	targetProb := new(big.Float).SetFloat64(probability)

	nth := new(big.Float).Quo(big.NewFloat(1), new(big.Float).SetInt64(n))
	failChance := bigfloat.Pow(targetProb, nth)
	fmt.Printf("n fail chance %s\n", failChance.Text('f', 100))

	successChance := new(big.Float).Sub(big.NewFloat(1), failChance)
	fmt.Printf("success chance %s\n", successChance.Text('f', 100))

	result := new(big.Float).Mul(successChance, uint256Max)
	fmt.Printf("result %s\n", result.Text('f', 0))
	return result
}

func esitimateHashRate(difficulty *big.Float, durS int, probability float64) *big.Float {
	ne := new(big.Float).Sub(uint256Max, difficulty)
	not_probability := new(big.Float).Quo(ne, uint256Max)

	targetProb := new(big.Float).SetFloat64(probability)

	tmp := new(big.Float).Mul(big.NewFloat(float64(durS)), bigfloat.Log(not_probability))
	hashRate := new(big.Float).Quo(bigfloat.Log(targetProb), tmp)
	hashRate = hashRate.Quo(hashRate, big.NewFloat(1000*1000))
	fmt.Println(hashRate.Text('f', 0))
	return hashRate
}
