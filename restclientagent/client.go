package restclientagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"td5/vtypes"
)

type RestClientAgent struct {
	id   string
	url  string
	alts []vtypes.Alternative
}

func NewRestClientAgent(id string, url string, alts []vtypes.Alternative) *RestClientAgent {
	return &RestClientAgent{id, url, alts}
}

func (rca *RestClientAgent) treatResponse(r *http.Response) int {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	var resp vtypes.Response
	json.Unmarshal(buf.Bytes(), &resp)

	return resp.Result
}

func (rca *RestClientAgent) doRequest() (res int, err error) {
	req := vtypes.Request{
		Choice: rca.alts,
	}

	// sérialisation de la requête
	url := rca.url + "/vote"
	data, _ := json.Marshal(req)

	// envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	// traitement de la réponse
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		return
	}
	res = rca.treatResponse(resp)

	return
}

func (rca *RestClientAgent) Start() {
	log.Printf("démarrage de %s", rca.id)
	_, err := rca.doRequest()

	if err != nil {
		log.Fatal(rca.id, "error:", err.Error())
	} else {
		log.Printf(rca.id, rca.alts)
	}
}
