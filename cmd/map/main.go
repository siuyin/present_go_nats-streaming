package main

import (
	"encoding/json"
	"fmt"
	"time"

	stan "github.com/nats-io/go-nats-streaming"
)

// 10 OMIT
const (
	clusterID = "test-cluster"
	clientID  = "newForOld-map"
	svrURL    = "nats://192.168.99.100:4222"
)

type Payload struct {
	TimeStamp time.Time
	Event     string // eg. IDMapped
	OldID     string //     o123
	NewID     string //     nABC
}

// 20 OMIT
// 30 OMIT
func main() {
	fmt.Println("Old ID to New ID Mapper")
	sc, _ := stan.Connect(clusterID, clientID, stan.NatsURL(svrURL))
	defer sc.Close()
	subject := "newForOld"
	pl := Payload{time.Now(), "IDMapped", "o123", "nABC"}
	plBytes, _ := json.Marshal(pl)
	sc.Publish(subject, plBytes)
	fmt.Printf("%s event created\n", plBytes)
}

// 40 OMIT
