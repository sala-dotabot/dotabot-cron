package dota

type MatchesResult struct {
	Result Matches `json:"result"`
}

type MatchResult struct {
	Result MatchDetails `json:"result"`
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

type MatchDetails struct {
	Players      []PlayerDetails `json:"players"`
	RadiantWin   bool            `json:"radiant_win"`
	Duration     uint64          `json:"duration"`
	StartTime    uint64          `json:"start_time"`
	MatchId      int64           `json:"match_id"`
	MatchSeqNum  int64           `json:"match_seq_num"`
	RadiantScore uint16          `json:"radiant_score"`
	DireScore    uint16          `json:"dire_score"`
}

type PlayerDetails struct {
	AccountId  int64 `json:"account_id"`
	PlayerSlot uint8 `json:"player_slot"`
	HeroId     uint8 `json:"hero_id"`
	TeamNumber uint8 `json:"team_number"`

	Kills   uint16 `json:"kills"`
	Deaths  uint16 `json:"deaths"`
	Assists uint16 `json:"assists"`

	LastHists  uint16 `json:"last_hists"`
	Denies     uint16 `json:"denies"`
	GoldPerMin uint16 `json:"gold_per_min"`
	XpPerMin   uint16 `json:"hp_per_min"`

	HeroDamage  uint32 `json:"hero_damage"`
	TowerDamage uint32 `json:"tower_damage"`
	HeroHealing uint32 `json:"hero_healing"`
}
