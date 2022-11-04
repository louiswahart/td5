package restserveragent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"td5/comsoc"
	"td5/vtypes"
	"time"
)

var supportedVotingMethods []string = []string{"borda", "copeland", "majority", "stv"}

type RestServerAgent struct {
	sync.Mutex
	id       string
	addr     string
	rule     string
	deadline time.Time
	voterIDs []string
	nbrAlts  int
	ballotID string
	Ballot   vtypes.Profile
}

func NewRestServerAgent(addr string) *RestServerAgent {
	ballot := make(vtypes.Profile, 0)
	return &RestServerAgent{id: addr, addr: addr, Ballot: ballot}
}

func (rsa *RestServerAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}

func (*RestServerAgent) decodeRequest(r *http.Request) (req vtypes.Request, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (*RestServerAgent) decodeNewBallotRequest(r *http.Request) (req vtypes.NewBallotRequest, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return req, err
}

func (*RestServerAgent) decodeVoteRequest(r *http.Request) (req vtypes.VoteRequest, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return req, err
}

func (*RestServerAgent) decodeResultRequest(r *http.Request) (req vtypes.ResultRequest, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return req, err
}

func (rsa *RestServerAgent) doNewBallot(w http.ResponseWriter, r *http.Request) {
	rsa.Lock()
	defer rsa.Unlock()

	// Vérification de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// Décodage de la requête
	req, err := rsa.decodeNewBallotRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}

	// Traitement de la requête
	if !contains(supportedVotingMethods, req.Rule) {
		err = errors.New("Voting method not supported")
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, err.Error())
		return
	}
	rsa.rule = req.Rule
	rsa.deadline = req.Deadline
	rsa.voterIDs = req.Voters
	rsa.nbrAlts = req.NbrAlts
	rsa.ballotID = "ballot"

	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(vtypes.NewBallotResponse{BallotID: rsa.ballotID})
	w.Write(serial)

	return
}

func (rsa *RestServerAgent) doVote(w http.ResponseWriter, r *http.Request) {
	rsa.Lock()
	defer rsa.Unlock()

	// Vérification de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// Décodage de la requête
	req, err := rsa.decodeVoteRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}

	// Traitement de la requête
	err = comsoc.CheckPrefs(req.Prefs, rsa.nbrAlts)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}

	if req.VoteID != rsa.ballotID {
		err = errors.New("ID Ballot invalide")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}

	rsa.voterIDs, err = removeIfContains(rsa.voterIDs, req.AgentID)
	if err != nil {
		err = errors.New("Votant non inscrit sur la liste")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, err.Error())
		return
	}

	if time.Now().After(rsa.deadline) {
		err = errors.New("Date limite dépassée")
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, err.Error())
		return
	}

	rsa.Ballot = append(rsa.Ballot, req.Prefs)

	w.WriteHeader(http.StatusOK)

	return
}

func (rsa *RestServerAgent) doResult(w http.ResponseWriter, r *http.Request) {
	rsa.Lock()
	defer rsa.Unlock()

	// Vérification de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// Décodage de la requête
	req, err := rsa.decodeResultRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}

	// Traitement de la requête
	if req.BallotID != rsa.ballotID {
		err = errors.New("ID Ballot invalide")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, err.Error())
		return
	}

	if len(rsa.voterIDs) > 0 && time.Now().Before(rsa.deadline) {
		err = errors.New("La procédure de vote n'est pas terminée")
		w.WriteHeader(http.StatusTooEarly)
		fmt.Fprintf(w, err.Error())
		return
	}

	// Calcul de la SWF
	var count vtypes.Count
	switch rsa.rule {
	case "borda":
		count, err = comsoc.BordaSWF(rsa.Ballot)
		break
	case "copeland":
		count, err = comsoc.CopelandSWF(rsa.Ballot)
		break
	case "majority":
		count, err = comsoc.MajoritySWF(rsa.Ballot)
		break
	case "stv":
		count, err = comsoc.StvSWF(rsa.Ballot)
		break
	default:
		err = errors.New("Procédure de vote non implémentée")
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, err.Error())
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, err.Error())
		return
	}

	ranking := comsoc.SortCount(count)
	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(vtypes.ResultResponse{Winner: ranking[0], Ranking: ranking})
	w.Write(serial)
}

func (rsa *RestServerAgent) Start() {
	// Création du multiplexeur
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", rsa.doNewBallot)
	mux.HandleFunc("/vote", rsa.doVote)
	mux.HandleFunc("/result", rsa.doResult)

	// Création du serveur
	s := &http.Server{
		Addr:           rsa.addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// Lancement du serveur
	log.Println("Listening on", rsa.addr)
	go log.Fatal(s.ListenAndServe())
}

func contains(list []string, s string) bool {
	for _, r := range list {
		if r == s {
			return true
		}
	}
	return false
}

func removeIfContains(list []string, s string) (newList []string, err error) {
	for i, r := range list {
		if r == s {
			newList = append(list[:i], list[i+1:]...)
			return newList, nil
		}
	}
	err = errors.New("La liste ne contient pas l'élément souhaité")
	return list, err
}
