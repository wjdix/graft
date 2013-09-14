package graft

import (
	"encoding/json"
	"github.com/bmizerany/pat"
	"io/ioutil"
	"net/http"
)

func HttpHandler(server *Server) http.Handler {
	return handler(server)
}

func extractMessage(r *http.Request, message interface{}) error {
	body, _ := ioutil.ReadAll(r.Body)
	return json.Unmarshal(body, message)
}

func handler(server *Server) http.Handler {
	m := pat.New()

	m.Post("/request_vote", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var message RequestVoteMessage
		err := extractMessage(r, &message)
		if err != nil {
			w.WriteHeader(422)
			return
		}

		response, _ := server.ReceiveRequestVote(message)

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(response)
	}))

	m.Post("/append_entries", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var message AppendEntriesMessage
		err := extractMessage(r, &message)
		if err != nil {
			w.WriteHeader(422)
			return
		}

		response, _ := server.ReceiveAppendEntries(message)

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(response)
	}))

	return m
}
