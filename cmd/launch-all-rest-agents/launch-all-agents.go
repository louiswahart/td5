package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"td5/restclientagent"
	"td5/restserveragent"
	"td5/vtypes"
	"time"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Veuillez respecter la syntaxe demandée")
		return
	}
	nbrAgts, err1 := strconv.Atoi(os.Args[1])
	nbrAlts, err2 := strconv.Atoi(os.Args[2])
	if err1 != nil || err2 != nil {
		fmt.Println(err1)
		fmt.Println(err2)
		return
	}
	rule := os.Args[3]
	const url1 = ":8080"
	const url2 = "http://localhost:8080"

	// Démarrage du serveur
	votingAgts := make([]restclientagent.VoteRestClientAgent, 0, nbrAgts)
	servAgt := restserveragent.NewRestServerAgent(url1)
	log.Println("Démarrage du serveur...")
	go servAgt.Start()

	time.Sleep(time.Second * 2)

	voterIDs := make([]string, nbrAgts)
	for i := 0; i < nbrAgts; i++ {
		voterIDs[i] = fmt.Sprintf("id%02d", i)
	}

	// Démarrage du NewBallotAgent
	var c chan string = make(chan string)
	newBallotAgt := restclientagent.NewNewBallotRestClientAgent("id00", url2, c, rule, time.Now().Add(5*time.Minute), voterIDs, nbrAlts)
	go newBallotAgt.Start()

	ballotID := <-c

	// Démarrage des VotingAgents

	log.Println("Démarrage des clients...")
	for i := 0; i < nbrAgts; i++ {
		prefs := make([]vtypes.Alternative, nbrAlts)
		a := make([]int, nbrAlts)
		for j := range a {
			a[j] = j + 1
		}
		rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
		for j := 0; j < nbrAlts; j++ {
			prefs[j] = vtypes.Alternative(a[j])
		}
		thresholds := []int{}
		if rule == "approval" {
			thresholds = []int{rand.Int() % nbrAlts}
		}
		votingAgt := restclientagent.NewVoteRestClientAgent(voterIDs[i], url2, ballotID, prefs, thresholds)
		votingAgts = append(votingAgts, *votingAgt)
	}

	for _, agt := range votingAgts {
		// Attention, obligation de passer par cette lambda pour faire capturer la valeur de l'itération par la goroutine
		func(agt restclientagent.VoteRestClientAgent) {
			go agt.Start()
		}(agt)
	}
	time.Sleep(time.Second * 5)

	// Démarrage du ResultAgent
	resultAgt := restclientagent.NewResultRestClientAgent("id00", url2, ballotID)
	go resultAgt.Start()

	fmt.Scanln()
}
