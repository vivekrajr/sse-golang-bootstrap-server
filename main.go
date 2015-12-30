//Simple example for server sent events
package main

import (
	"fmt"
	"github.com/JanBerktold/sse"
	"net/http"
	"time"
)

//The channel which gets the events
var eventChannel chan int

//When a client listens, for SSE, this is the handler function
func HandleSSE(w http.ResponseWriter, r *http.Request) {

	//Add header to allow requests from any origin
	w.Header().Add("Access-Control-Allow-Origin", "*")

	conn, err := sse.Upgrade(w, r)

	if err != nil {
		fmt.Printf("Error occured: %q", err.Error())
		return
	}

	for {
		select {

		//Listen on eventChannel for an integer
		case number := <-eventChannel:
			//Write the integer received to the stream
			conn.WriteString(fmt.Sprintf("%d", number))
		}
	}

}

func main() {

	//create channel for events
	eventChannel = make(chan int)

	//Start a goroutine for producing "events"
	go producer()

	http.HandleFunc("/event", HandleSSE)

	http.ListenAndServe(":3000", nil)
}

//Infinitely send 1 to 100 to channel with 100ms delay
func producer() {
	i := 1

	for {
		eventChannel <- i
		if i == 100 {
			i = 1
		} else {
			i++
		}
		time.Sleep(time.Millisecond * 100)
	}
}
