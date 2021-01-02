package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		res, _ := json.Marshal(map[string]interface{}{
			"name":    "go-socket",
			"version": "1.0.0",
		})

		w.Write(res)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	channel := newChannel()
	go channel.run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(channel, w, r)
	})

	http.ListenAndServe(":3000", nil)
}
