package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	stan "github.com/nats-io/go-nats-streaming"
)

const (
	clusterID = "test-cluster"
	clientID  = "newForOld-renameMedia"
	svrURL    = "nats://192.168.99.100:4222"
)

type Payload struct {
	TimeStamp time.Time
	Event     string // eg. RenameRequested / Renamed / RenameFailed
	OldID     string //     o123
	NewID     string //     nABC
}

func main() {
	fmt.Println("newForOld renameMedia reactor")
	sc, _ := stan.Connect(clusterID, clientID, stan.NatsURL(svrURL),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	defer sc.Close()

	// 10 OMIT
	startOpt := stan.DeliverAllAvailable()
	subject := "newForOld"
	sc.Subscribe(subject, func(msg *stan.Msg) {
		if !bytes.Contains(msg.Data, []byte(`"RenameRequested"`)) { // HL
			return
		}
		pl := Payload{}
		json.Unmarshal(msg.Data, &pl)
		// Actually do renaming here ...
		fmt.Printf("!! renamed: %q -> %q\n", pl.OldID, pl.NewID)
		pl.Event = "Renamed" // HL
		plBytes, _ := json.Marshal(pl)
		sc.Publish(subject, plBytes) // HL
		fmt.Printf("%s event created\n", plBytes)
	}, startOpt)
	// 20 OMIT
	// 30 OMIT
	http.HandleFunc("/rename", func(w http.ResponseWriter, req *http.Request) {
		oID := req.FormValue("oID")
		nID := req.FormValue("nID")
		subject := "newForOld"
		pl := Payload{time.Now(), "RenameRequested", oID, nID}
		plBytes, _ := json.Marshal(pl)
		sc.Publish(subject, plBytes) // HL
		fmt.Fprintf(w, "%s event created\n", plBytes)
	})
	// 40 OMIT

	log.Fatal(http.ListenAndServe(":8080", nil))
}
