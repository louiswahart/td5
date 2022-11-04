package main

import (
	"fmt"
	"log"
	"math/rand"
	"td5/restclientagent"
	"td5/restserveragent"
	"td5/vtypes"
	"time"
)

func main() {
	const n = 5
	const url1 = ":8080"
	const url2 = "http://localhost:8080"

	clAgts := make([]restclientagent.RestClientAgent, 0, n)
	servAgt := restserveragent.NewRestServerAgent(url1)
	log.Println("Démarrage du serveur...")
	go servAgt.Start()

	log.Println("Démarrage des clients...")
	for i := 0; i < n; i++ {
		id := fmt.Sprintf("id%02d", i)
		myChoice := make([]vtypes.Alternative, n)
		a := make([]int, n)
		for i := range a {
			a[i] = i + 1
		}
		rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
		for j := 0; j < n; j++ {
			myChoice[j] = vtypes.Alternative(a[j])
		}
		agt := restclientagent.NewRestClientAgent(id, url2, myChoice)
		clAgts = append(clAgts, *agt)
	}

	for _, agt := range clAgts {
		// Attention, obligation de passer par cette lambda pour faire capturer la valeur de l'itération par la goroutine
		func(agt restclientagent.RestClientAgent) {
			go agt.Start()
		}(agt)
	}
	time.Sleep(time.Second * 5)
	fmt.Println(servAgt.Ballot)
	fmt.Scanln()
}
