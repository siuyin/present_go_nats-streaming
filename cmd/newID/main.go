package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	stan "github.com/nats-io/go-nats-streaming"
)

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

func main() {
	fmt.Println("newForOld newID Query Handler", time.Now())
	sc, _ := stan.Connect(clusterID, clientID, stan.NatsURL(svrURL),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	defer sc.Close()

	idAccess := sync.Mutex{}
	idMap := make(map[string]string)
	// 10 OMIT
	startOpt := stan.DeliverAllAvailable()
	subject := "newForOld"
	sc.Subscribe(subject, func(msg *stan.Msg) {
		if !bytes.Contains(msg.Data, []byte(`"IDMapped"`)) { // 1. filter // HL
			return
		}
		pl := Payload{}
		json.Unmarshal(msg.Data, &pl)
		idAccess.Lock()
		idMap[pl.OldID] = pl.NewID // 2. update "DB" // HL
		idAccess.Unlock()
	}, startOpt)
	// 20 OMIT
	// 30 OMIT
	http.HandleFunc("/id", func(w http.ResponseWriter, req *http.Request) {
		id := req.FormValue("id")
		idAccess.Lock()
		fmt.Fprintf(w, "new ID for:%q is :%q\n", id, idMap[id])
		idAccess.Unlock()
	})
	// 40 OMIT
	// 50 OMIT
	http.HandleFunc("/ids", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Listing all IDs. Time:%v\n", time.Now().Format("15:04:05.000"))
		idAccess.Lock()
		for k, v := range idMap {
			fmt.Fprintf(w, "%s %s\n", k, v)
		}
		idAccess.Unlock()
	})
	// 60 OMIT
	log.Fatal(http.ListenAndServe(":8080", nil))
}
