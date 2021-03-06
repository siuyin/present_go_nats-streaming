package main

import (
	"fmt"
	"log"
	"time"

	stan "github.com/nats-io/go-nats-streaming"
)

// 10 OMIT
const (
	clusterID = "test-cluster"
	clientID  = "helloworld-sub"
	svrURL    = "nats://192.168.99.100:4222"
)

// 20 OMIT
// 30 OMIT
func main() {
	fmt.Println("Hello World Subscriber", time.Now())
	sc, _ := stan.Connect(clusterID, clientID, stan.NatsURL(svrURL),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	defer sc.Close()
	startOpt := stan.DeliverAllAvailable() // 1. Change me // HL
	// StartAtSequence(n),StartAtTime(time.Date(2018,M,D, 0,0,0,0,time.UTC))
	// StartAtTimeDelta(dur), StartWithLastReceived(), DeliverAllAvailable()
	subject := "helloWorld"
	sc.Subscribe(subject, func(msg *stan.Msg) {
		fmt.Printf("%s\n", msg)
	}, startOpt, stan.DurableName("")) // 2. try setting to "whiteboard-1" // HL
	select {} // wait forever
}

// 40 OMIT
