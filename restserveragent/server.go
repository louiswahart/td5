package restserveragent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"td5/vtypes"
	"time"
)

type RestServerAgent struct {
	sync.Mutex
	id       string
	reqCount int
	addr     string
	Ballot   vtypes.Profile
}

func NewRestServerAgent(addr string) *RestServerAgent {
	ballot := make(vtypes.Profile, 0)
	return &RestServerAgent{id: addr, addr: addr, Ballot: ballot}
}

// Test de la méthode
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

func (rsa *RestServerAgent) doVote(w http.ResponseWriter, r *http.Request) {
	// mise à jour du nombre de requêtes
	rsa.Lock()
	defer rsa.Unlock()
	rsa.reqCount++

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	rsa.Ballot = append(rsa.Ballot, req.Choice)

	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(0)
	w.Write(serial)

	return

}

func (rsa *RestServerAgent) doReqcount(w http.ResponseWriter, r *http.Request) {
	if !rsa.checkMethod("GET", w, r) {
		return
	}

	w.WriteHeader(http.StatusOK)
	rsa.Lock()
	defer rsa.Unlock()
	serial, _ := json.Marshal(rsa.reqCount)
	w.Write(serial)
}

func (rsa *RestServerAgent) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/reqcount", rsa.doReqcount)
	mux.HandleFunc("/vote", rsa.doVote)

	// création du serveur http
	s := &http.Server{
		Addr:           rsa.addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// lancement du serveur
	log.Println("Listening on", rsa.addr)
	go log.Fatal(s.ListenAndServe())
}
