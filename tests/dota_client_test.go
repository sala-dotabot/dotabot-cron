package tests

import (
	"dotabot-cron/dota"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGetMatches(t *testing.T) {
	json := `
	{
		"result":{
			"status":1,
			"num_results":1,
			"total_results":500,
			"results_remaining":499,
			"matches":[
				{
					"match_id":3514877261,
					"match_seq_num":3059114721,
					"start_time":1508598995,
					"lobby_type":8,
					"radiant_team_id":0,
					"dire_team_id":0,
					"players":[
						{
							"account_id":245120373,
							"player_slot":0,
							"hero_id":34
						},
						{
							"account_id":273379500,
							"player_slot":128,
							"hero_id":22
						}
					]
				}
			]			
		}
	}
	`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, json)
	}))
	defer ts.Close()

	client, err := dota.NewDotaClient(ts.URL, "asdasd")
	account_id := "359953025"
	var start_at_match_id uint64 = 0

	result, err := client.GetMatchesHistory(account_id, start_at_match_id)

	if err != nil {
		t.Fatal(err)
	}

	if result.Result.NumResults != 1 {
		t.Fatal("NumResults must be 1, got: " + strconv.Itoa(result.Result.NumResults))
	}

	match := result.Result.Matches[0]

	if match.MatchId != 3514877261 {
		t.Fatalf("MatchId must be 3514877261, got: %d", match.MatchId)
	}
	if match.MatchSeqNum != 3059114721 {
		t.Fatalf("MatchSeqNum must be 3059114721, got: %d", match.MatchSeqNum)
	}
	if match.StartTime != 1508598995 {
		t.Fatalf("StartTime must be 1508598995, got: %d", match.StartTime)
	}
	if match.LobbyType != 8 {
		t.Fatalf("LobbyType must be 8, got: %d", match.LobbyType)
	}
	if match.RadiantTeamId != 0 {
		t.Fatalf("RadiantTeamId must be %d, got: %d", 0, match.RadiantTeamId)
	}
	if match.DireTeamId != 0 {
		t.Fatalf("DireTeamId must be %d, got: %d", 0, match.DireTeamId)
	}
	if len(match.Players) != 2 {
		t.Fatalf("Players.len must be %d, got: %d", 2, len(match.Players))
	}
	for index, player := range match.Players {
		var accountId int64
		var heroId uint8
		var playerSlot uint8
		switch index {
		case 0:
			accountId = 245120373
			heroId = 34
			playerSlot = 0
			break
		case 1:
			accountId = 273379500
			heroId = 22
			playerSlot = 128
			break
		}

		if player.AccountId != accountId {
			t.Fatalf("AccountId must be %d, got: %d", accountId, player.AccountId)
		}
		if player.HeroId != heroId {
			t.Fatalf("HeroId must be %d, got: %d", heroId, player.HeroId)
		}
		if player.PlayerSlot != playerSlot {
			t.Fatalf("PlayerSlot must be %d, got: %d", playerSlot, player.PlayerSlot)
		}
	}
}
