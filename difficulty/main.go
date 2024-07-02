package main

import (
	"fmt"
	"math/big"

	"github.com/ALTree/bigfloat"
)

// 1000m 4443685057045916916839494277496125038705878205143482420972918669312
// 500m 8887370114091833833678988554992250077411756410286964841945837338624
// 10m 740614270448018082667345211377796923799185715493156038752666115899392
// 8m  555460709263765739036470010701196062214039696708679004195670928130048
// 6m  740614270448018082667345211377796923799185715493156038752666115899392
// 3m 1481228540896036165334690422755593847598371430986312077505332231798784
// 2m 2221842798488549893930113429797694032668256326301844165995655665287168
func main() {
	esitimateDifficulty(500, 20, 0.00001)

	//difficulty, _ := new(big.Int).SetString("740614270448018082667345211377796923799185715493156038752666115899392", 10)
	//number := esitimateHashRate(difficulty, 20*15, 0.00001)
	//fmt.Print(number)
}

func esitimateDifficulty(power int, blockNumber int, probability float64) {
	n := int64(power * 1000 * 1000 * blockNumber * 15)
	uint256Max, _ := new(big.Float).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935")
	fmt.Printf("uint256_max %s\n", uint256Max.Text('f', 0))

	targetPer := new(big.Float).SetFloat64(probability)
	fmt.Printf("target_per %s\n", targetPer.Text('f', 5))

	nth := new(big.Float).Quo(big.NewFloat(1), new(big.Float).SetInt64(n))
	failChance := bigfloat.Pow(targetPer, nth)
	fmt.Printf("n fail chance %s\n", failChance.Text('f', 100))

	target := bigfloat.Pow(failChance, big.NewFloat(float64(n)))
	fmt.Printf("check target %s\n", target.Text('f', 100))

	successChance := new(big.Float).Sub(big.NewFloat(1), failChance)
	fmt.Printf("success chance %s\n", successChance.Text('f', 100))

	result := new(big.Float).Mul(successChance, uint256Max)
	fmt.Printf("result %s\n", result.Text('f', 0))
}

func esitimateHashRate(difficulty *big.Int, durS int, probability float64) int {
	uint256Max, _ := new(big.Float).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935")
	diffilcultyF, _ := new(big.Float).SetString(difficulty.String())
	ne := new(big.Float).Sub(uint256Max, diffilcultyF)
	not_probability := new(big.Float).Quo(ne, uint256Max)

	targetPer := new(big.Float).SetFloat64(probability)

	number := 1
	for {
		//todo find a bigfloat log n algo
		//this one is not accurate
		n := int64(number * 1000 * 1000)
		x := bigfloat.Pow(not_probability, new(big.Float).SetInt64(n))
		if x.Cmp(targetPer) == -1 {
			break
		}
		number++
	}

	return number / durS
}
