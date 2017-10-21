package dota

type MatchesResult struct {
	Result Matches `json:"result"`
}

type Matches struct {
	Status           int     `json:"status"`
	NumResults       int     `json:"num_results"`
	TotalResults     int     `json:"total_results"`
	ResultsRemaining int     `json:"results_remaining"`
	Matches          []Match `json:"matches"`
}

type LobbyType int

const (
	PUBLIC LobbyType = iota
	PRACTICE
	TOURNAMENT
	TUTORIAL
	COOP_WITH_BOTS
	TEAM_MATCH
	SOLO_QUEUE
	RANKED_MATCHMAKING
	SOLO_MID
)

type Match struct {
	MatchId         int64     `json:"match_id"`
	MatchSeqNum     int64     `json:"match_seq_num"`
	StartTime       uint64    `json:"start_time"`
	LobbyType       LobbyType `json:"lobby_type"`
	RadiantTeamId   int64     `json:"radiant_team_id"`
	DireTeamId      int64     `json:"dire_team_id"`
	TournamentId    int64     `json:"tournament_id"`
	TournamentRound int64     `json:"tournament_round"`
	Players         []Player  `json:"players"`
}

type Player struct {
	AccountId  int64 `json:"account_id"`
	PlayerSlot uint8 `json:"player_slot"`
	HeroId     uint8 `json:"hero_id"`
}
