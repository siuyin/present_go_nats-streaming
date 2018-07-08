package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	stan "github.com/nats-io/go-nats-streaming"
)

// 10 OMIT
const (
	clusterID = "test-cluster"
	clientID  = "newForOld-newID"
	svrURL    = "nats://192.168.99.100:4222"
)

type Payload struct {
	TimeStamp time.Time
	Event     string // eg. IDMapped
	OldID     string //     o123
	NewID     string //     nABC
}

// 20 OMIT
func main() {
	fmt.Println("newForOld newID Query Handler", time.Now())
	// 30 OMIT
	sc, _ := stan.Connect(clusterID, clientID, stan.NatsURL(svrURL),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	defer sc.Close()
	// 31 OMIT
	idAccess := sync.Mutex{}
	idMap := make(map[string]string)
	newData := make(chan bool)
	// 32 OMIT
	// 33 OMIT
	startOpt := stan.DeliverAllAvailable()
	subject := "newForOld"
	sc.Subscribe(subject, func(msg *stan.Msg) {
		pl := Payload{}
		if !bytes.Contains(msg.Data, []byte(`"IDMapped"`)) {
			return
		}
		json.Unmarshal(msg.Data, &pl)
		idAccess.Lock()
		idMap[pl.OldID] = pl.NewID
		idAccess.Unlock()
		newData <- true
	}, startOpt)
	// 34 OMIT
	listIDs(idMap, newData)
	select {} // wait forever
	// 40 OMIT
}

// 50 OMIT
func listIDs(id map[string]string, ch <-chan bool) {
	go func() {
		for {
			<-ch // wait for new data
			fmt.Println("New Listing", time.Now())
			for k, v := range id {
				fmt.Printf("%s %s\n", k, v)
			}
		}
	}()
}

// 60 OMIT
