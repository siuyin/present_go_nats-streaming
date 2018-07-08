package main

import (
	"fmt"
	"time"

	stan "github.com/nats-io/go-nats-streaming"
)

// 10 OMIT
const (
	clusterID = "test-cluster"
	clientID  = "helloworld-pub"
	svrURL    = "nats://192.168.99.100:4222"
)

// 20 OMIT
// 30 OMIT
func main() {
	fmt.Println("Hello World Publisher")
	sc, _ := stan.Connect(clusterID, clientID, stan.NatsURL(svrURL))
	defer sc.Close()
	subject := "helloWorld"
	for n := 0; n < 3; n++ {
		msg, _ := time.Now().MarshalText()
		sc.Publish(subject, msg)
		fmt.Printf("%s: %s\n", subject, msg)
	}
}

// 40 OMIT
