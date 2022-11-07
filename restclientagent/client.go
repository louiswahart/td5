package restclientagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"td5/vtypes"
	"time"
)

type RestClientAgent struct {
	id  string
	url string
}

type NewBallotRestClientAgent struct {
	RestClientAgent
	ballotIDchan chan string
	rule         string
	deadline     time.Time
	voterIDs     []string
	nbrAlts      int
}

type VoteRestClientAgent struct {
	RestClientAgent
	voteID  string
	prefs   []vtypes.Alternative
	options []int
}

type ResultRestClientAgent struct {
	RestClientAgent
	ballotID string
}

func NewNewBallotRestClientAgent(id string, url string, ballotIDchan chan string, rule string, deadline time.Time, voterIDs []string, nbrAlts int) *NewBallotRestClientAgent {
	return &NewBallotRestClientAgent{RestClientAgent{id, url}, ballotIDchan, rule, deadline, voterIDs, nbrAlts}
}

func NewVoteRestClientAgent(id string, url string, voteID string, prefs []vtypes.Alternative, options []int) *VoteRestClientAgent {
	return &VoteRestClientAgent{RestClientAgent{id, url}, voteID, prefs, options}
}

func NewResultRestClientAgent(id string, url string, ballotID string) *ResultRestClientAgent {
	return &ResultRestClientAgent{RestClientAgent{id, url}, ballotID}
}

func (rca *NewBallotRestClientAgent) doRequest() (ballotID string, err error) {
	req := vtypes.NewBallotRequest{Rule: rca.rule, Deadline: rca.deadline, Voters: rca.voterIDs, NbrAlts: rca.nbrAlts}

	// Serialisation de la requête
	url := rca.url + "/new_ballot"
	data, _ := json.Marshal(req)

	// Envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	// Traitement de la requête
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		return
	}
	ballotID = rca.treatResponse(resp)

	return ballotID, err
}

func (rca *VoteRestClientAgent) doRequest() (err error) {
	req := vtypes.VoteRequest{AgentID: rca.id, VoteID: rca.voteID, Prefs: rca.prefs, Options: rca.options}

	// Serialisation de la requête
	url := rca.url + "/vote"
	data, _ := json.Marshal(req)

	// Envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	// Traitement de la requête
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		return
	}

	return err
}

func (rca *ResultRestClientAgent) doRequest() (winner vtypes.Alternative, ranking []vtypes.Alternative, err error) {
	req := vtypes.ResultRequest{BallotID: rca.ballotID}

	// Serialisation de la requête
	url := rca.url + "/result"
	data, _ := json.Marshal(req)

	// Envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	// Traitement de la requête
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		return
	}

	winner, ranking = rca.treatResponse(resp)

	return winner, ranking, err
}

func (rca *NewBallotRestClientAgent) treatResponse(r *http.Response) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	var resp vtypes.NewBallotResponse
	json.Unmarshal(buf.Bytes(), &resp)

	return resp.BallotID
}

func (rca *ResultRestClientAgent) treatResponse(r *http.Response) (winner vtypes.Alternative, ranking []vtypes.Alternative) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	var resp vtypes.ResultResponse
	json.Unmarshal(buf.Bytes(), &resp)

	return resp.Winner, resp.Ranking
}

func (rca *NewBallotRestClientAgent) Start() {
	log.Printf("Démarrage de %s", rca.id)
	ballotID, err := rca.doRequest()

	if err != nil {
		log.Fatal(rca.id, "error:", err.Error())
	} else {
		log.Printf(rca.id, " a créé le ballot : ", ballotID)
		rca.ballotIDchan <- ballotID
	}
}

func (rca *VoteRestClientAgent) Start() {
	log.Printf("Démarrage de %s", rca.id)
	err := rca.doRequest()

	if err != nil {
		log.Fatal(rca.id, "error :", err.Error())
	} else {
		log.Printf("%s - a voté : %v", rca.id, rca.prefs)
	}
}

func (rca *ResultRestClientAgent) Start() {
	log.Printf("Démarrage de %s", rca.id)
	winner, ranking, err := rca.doRequest()

	if err != nil {
		log.Fatal(rca.id, "error:", err.Error())
	} else {
		log.Printf("%s - Gagnant du ballot (id : %s) : %d - ranking : %v", rca.id, rca.ballotID, winner, ranking)
	}
}
