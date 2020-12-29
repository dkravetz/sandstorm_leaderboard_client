package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func GetPlayers(allPlayers *Players, wg *sync.WaitGroup) error {
	totalPlayerCount := 50000
	client := http.Client{Timeout: 10 * time.Second}
	// https://leaderboard.sandstorm.game/api/v0/Players/GetPlayerCount

	pagination := 500

	//ch := make(chan Players)
	for startRank := 1; startRank < totalPlayerCount; startRank += pagination {
		go func(i int) {
			wg.Add(1)
			req, err := http.NewRequest(
				"GET",
				"https://leaderboard.sandstorm.game/api/v0/Players/GetRankedPlayers",
				nil,
			)
			if err != nil {
				log.Fatal("Couldn't request URL. ", err)
			}
			q := req.URL.Query()
			q.Add("startRank", fmt.Sprintf("%d", i))
			q.Add("endRank", fmt.Sprintf("%d", i+pagination))
			req.URL.RawQuery = q.Encode()

			res, err := client.Do(req)
			if err != nil {
				log.Fatal("Couldn't request URL. ", err)
			}

			if res.StatusCode != 200 {
				log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
			}

			var pagePlayers Players
			err = json.NewDecoder(res.Body).Decode(&pagePlayers)
			res.Body.Close()
			if err != nil {
			}
			*allPlayers = append(*allPlayers, pagePlayers...)
			wg.Done()
		}(startRank)
	}
	return nil
}
