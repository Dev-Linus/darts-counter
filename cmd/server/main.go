package main

import (
	"log"
	"net/http"

	handler "darts-counter/cmd/server/http"
	"darts-counter/darts"
	"darts-counter/response"
	"darts-counter/storage"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	store := storage.NewStorage("darts.db")
	responseBuilder := response.NewBuilder()
	service := darts.NewService(store, responseBuilder)
	api, err := handler.NewApi(store, service)
	if err != nil {
		log.Fatal("Api could be initialized")
	}

	mux := http.NewServeMux()

	// CRUD Player
	mux.HandleFunc("/createPlayer", api.CreatePlayer)
	mux.HandleFunc("/updatePlayer", api.UpdatePlayer)
	mux.HandleFunc("/listPlayers", api.ListPlayers)
	mux.HandleFunc("/deletePlayer", api.DeletePlayer)

	// CRD Matches
	mux.HandleFunc("/createMatch", api.CreateMatch)
	mux.HandleFunc("/listMatches", api.ListMatches)
	mux.HandleFunc("/deleteMatch", api.DeleteMatch)
	mux.HandleFunc("/getMatch", api.GetMatch)

	// gameplay
	mux.HandleFunc("/playerThrow", api.PlayerThrow)

	// misc
	mux.HandleFunc("/statistics", api.Statistics)
	// mux.HandleFunc("/settings", api.Settings)

	// media streaming
	mux.HandleFunc("/streamFile", api.StreamFile)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(mux)))
}
