package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
)

type RewardsPerWallet struct {
	Address     string `csv:"Address"`
	TodayPoints string `csv:"TodayPoints"`
	TotalPoints string `csv:"TotalPoints"`
	Day         int    `csv:"Day"`
	Phase       int    `csv:"Phase"`
	CreatedAt   time.Time
}

type POWHashrate struct {
	ID        string  `gorm:"primaryKey;index"`
	Address   string  // public address for resource provider
	Date      int64   // timestamp for hashrate calculation
	Hashrate  float64 // hashrate value
	CreatedAt time.Time
}

var addressesStr string

func main() {
	hashrateFile, err := os.Open("hashrates.csv")
	if err != nil {
		panic(err)
	}
	defer hashrateFile.Close()

	var hashrates []POWHashrate
	if err := gocsv.UnmarshalFile(hashrateFile, &hashrates); err != nil { // Load clients from file
		panic(err)
	}

	rewardFile, err := os.Open("rewards.csv")
	if err != nil {
		panic(err)
	}
	defer hashrateFile.Close()

	var rewards []RewardsPerWallet
	if err := gocsv.UnmarshalFile(rewardFile, &rewards); err != nil { // Load clients from file
		panic(err)
	}

	flag.StringVar(&addressesStr, "address", "", "")

	addressesStr = strings.ToUpper(addressesStr)
	addresses := strings.Split(addressesStr, ",")

	latestDate := findMaxDate(rewards).Truncate(time.Hour * 24)

	for _, address := range addresses {
		fmt.Println("Address: ", address)
		for i := 0; i < 5; i++ {
			startDate := latestDate.Add(time.Hour * -24)
			endDate := latestDate

			reward := findRewardsInTheDay(rewards, startDate, endDate)[0]
			hashratesOfAddress := findHashsrateInTheDay(hashrates, startDate, endDate)

			fmt.Println("	", reward.CreatedAt.String(), "rewards: ", reward.TodayPoints, reward.TotalPoints, "hashratesNumber ", len(hashratesOfAddress))
			latestDate = startDate
		}
	}
}

func findMaxDate(rewards []RewardsPerWallet) time.Time {
	latestDate := rewards[0].CreatedAt
	for _, reward := range rewards {
		if reward.CreatedAt.Compare(latestDate) >= 0 {
			latestDate = reward.CreatedAt
		}
	}
	return latestDate
}

func findHashsrateInTheDay(hashrates []POWHashrate, start, end time.Time) []POWHashrate {
	var results []POWHashrate
	hasSeen := make(map[string]struct{})
	for _, hashrate := range hashrates {
		if (hashrate.CreatedAt.Compare(start) == 1 || hashrate.CreatedAt.Compare(start) == 0) && hashrate.CreatedAt.Compare(end) == -1 {
			key := hashrate.Address + strconv.FormatInt(hashrate.Date, 10)
			if _, ok := hasSeen[key]; ok {
				continue
			}
			hasSeen[key] = struct{}{}
			results = append(results, hashrate)
		}
	}
	return results
}

func findRewardsInTheDay(rewards []RewardsPerWallet, start, end time.Time) []RewardsPerWallet {
	var results []RewardsPerWallet
	hasSeen := make(map[string]struct{})
	for _, reward := range rewards {
		if (reward.CreatedAt.Compare(start) == 1 || reward.CreatedAt.Compare(start) == 0) && reward.CreatedAt.Compare(end) == -1 {
			key := reward.Address + strconv.Itoa(reward.Day) + strconv.Itoa(reward.Phase)
			if _, ok := hasSeen[key]; ok {
				continue
			}
			hasSeen[key] = struct{}{}
			results = append(results, reward)
		}
	}
	return results
}
