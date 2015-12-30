package main

import (
	"fmt"
	"github.com/JanBerktold/sse"
	"net/http"
	"time"
)

var myChannel chan int

func HandleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	conn, err := sse.Upgrade(w, r)

	if err != nil {
		fmt.Printf("Error occured: %q", err.Error())
		return
	}

	for {
		select {
		case number := <-myChannel:
			conn.WriteString(fmt.Sprintf("%d", number))
		}
	}
}

func main() {

	myChannel = make(chan int)

	go producer()

	http.HandleFunc("/event", HandleSSE)

	http.ListenAndServe(":3000", nil)
}

func producer() {
	i := 1

	for {
		myChannel <- i
		if i == 100 {
			i = 1
		} else {
			i++
		}
		time.Sleep(time.Millisecond * 100)
	}
}
