package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	stan "github.com/nats-io/go-nats-streaming"
)

// 10 OMIT
const (
	clusterID = "test-cluster"
	clientID  = "newForOld-markMilestone"
	svrURL    = "nats://192.168.99.100:4222"
)

type Payload struct {
	TimeStamp time.Time
	Event     string // eg. IDMapped
	OldID     string //     o123
	NewID     string //     nABC
}
type Milestone struct {
	TimeStamp time.Time
	Event     string // eg. MilestoneMarked
	N         int
}

// 20 OMIT
func main() {
	const milestoneInterval = 2
	fmt.Println("newForOld Mark Milestone Command Handler", time.Now())
	sc, _ := stan.Connect(clusterID, clientID, stan.NatsURL(svrURL),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	defer sc.Close()
	subject := "newForOld"
	var n int
	// 30 OMIT
	startOpt := stan.DeliverAllAvailable()
	sc.Subscribe(subject, func(msg *stan.Msg) { // HL
		if !bytes.Contains(msg.Data, []byte(`"IDMapped"`)) {
			return
		}
		pl := Payload{}
		json.Unmarshal(msg.Data, &pl)
		n++
		if n%milestoneInterval == 0 {
			ms := Milestone{pl.TimeStamp, "MilestoneMarked", n}
			milestone, _ := json.Marshal(ms)
			sc.Publish(subject, milestone) // HL
			fmt.Printf("Milestone: %d marked at %v\n", n, pl.TimeStamp)
		}
	}, startOpt)
	// 40 OMIT
	select {} // wait forever
}
