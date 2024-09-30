package main

import (
	"fmt"
	"math/big"

	"github.com/ALTree/bigfloat"
)

var (
	BlockInterval = 12
	uint256Max, _ = new(big.Float).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935")
)

// 40m 92576780592126171815437600338300430792573009392238517278497593884672

// 8m 555460709263765739036470010701196062214039696708679004195670928130048
func main() {
	targetDifficulty := esitimateDifficulty(8, 300, 0.00001)

	number := esitimateHashRate(targetDifficulty, 30*BlockInterval, 0.00001)
	fmt.Printf("%s %s\n", targetDifficulty.Text('f', 10), number.String())

	fmt.Println("********all possible hashrates***********")
	//12 到 12*30
	for i := 0; i < 30; i++ {
		changeRange := new(big.Float).Quo(targetDifficulty, big.NewFloat(10))

		minDifficulty := new(big.Float).Sub(targetDifficulty, changeRange)
		minHashrates := esitimateHashRate(minDifficulty, BlockInterval*(i+1), 0.00001)

		maxDifficulty := new(big.Float).Add(targetDifficulty, changeRange)
		maxHashrates := esitimateHashRate(maxDifficulty, BlockInterval*(i+1), 0.00001)

		fmt.Printf("%ds	min: %s		max: %s\n", BlockInterval*(i+1), minHashrates.String(), maxHashrates.String())
	}

}

func esitimateDifficulty(power int, durS int, probability float64) *big.Float {
	n := int64(power * 1000 * 1000 * durS)

	targetProb := new(big.Float).SetFloat64(probability)

	nth := new(big.Float).Quo(big.NewFloat(1), new(big.Float).SetInt64(n))
	failChance := bigfloat.Pow(targetProb, nth)
	fmt.Printf("n fail chance %s\n", failChance.Text('f', 100))

	successChance := new(big.Float).Sub(big.NewFloat(1), failChance)
	fmt.Printf("success chance %s\n", successChance.Text('f', 100))

	result := new(big.Float).Mul(successChance, uint256Max)
	return result
}

func esitimateHashRate(difficulty *big.Float, durS int, probability float64) *big.Float {
	ne := new(big.Float).Sub(uint256Max, difficulty)
	not_probability := new(big.Float).Quo(ne, uint256Max)

	targetProb := new(big.Float).SetFloat64(probability)

	tmp := new(big.Float).Mul(big.NewFloat(float64(durS)), bigfloat.Log(not_probability))
	hashRate := new(big.Float).Quo(bigfloat.Log(targetProb), tmp)
	hashRate = hashRate.Quo(hashRate, big.NewFloat(1000*1000))
	return hashRate
}
