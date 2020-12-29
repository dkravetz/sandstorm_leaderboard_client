package main

import "fmt"

type Player struct{
	Rank int `json:"rank"`
	SteamId string `json:"steam_id"`
	PlayerAlias string `json:"player_alias"`
	Tier string `json:"tier"`
	DivisionText string `json:"division_text"`
	PlacementsLeft int `json:"placements_left"`
}

type Players []Player

func (p Player) String() string {
	return fmt.Sprintf("%s is %s %s. Current competitive rank is: %d", p.PlayerAlias, p.Tier, p.DivisionText, p.Rank)
}

// String (int) is necessary to satisfy the Source type for the fuzzy finder
func (p Players) String(i int) string {
	return p[i].PlayerAlias
}

// Len is necessary to satisfy the Source type for the fuzzy finder
func (p Players) Len() int {
	return len(p)
}