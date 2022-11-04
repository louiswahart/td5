package vtypes

import "time"

type Request struct {
	Choice []Alternative `json:"alts"`
}

type NewBallotRequest struct {
	Rule     string    `json:"rule"`
	Deadline time.Time `json:"deadline"`
	Voters   []string  `json:"voter-ids"`
	NbrAlts  int       `json:"#alts"`
}

type VoteRequest struct {
	AgentID string        `json:"agent-id"`
	VoteID  string        `json:"vote-id"`
	Prefs   []Alternative `json:"prefs"`
	Options []int         `json:"options"`
}

type ResultRequest struct {
	BallotID string `json:"ballot-id"`
}

type Response struct {
	Result int `json:"res"`
}

type NewBallotResponse struct {
	BallotID string `json:"ballot-id"`
}

type ResultResponse struct {
	Winner  Alternative   `json:"winner"`
	Ranking []Alternative `json:"ranking"`
}

type Alternative int

type Profile [][]Alternative

type Count map[Alternative]int
