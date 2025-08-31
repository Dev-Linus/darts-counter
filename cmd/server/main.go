package main

import (
	handler "darts-counter/cmd/server/http"
	"darts-counter/darts"
	"darts-counter/storage"
	"log"
	"net/http"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	store := storage.NewStorage("darts.db")
	darts := darts.NewService(store)
	svc, err := handler.NewApi(store, darts)
	if err != nil {
		log.Fatal("Api could be initialized")
	}

	mux := http.NewServeMux()

	// CRUD Player
	mux.HandleFunc("/createPlayer", svc.CreatePlayer)
	mux.HandleFunc("/updatePlayer", svc.UpdatePlayer)
	mux.HandleFunc("/listPlayers", svc.ListPlayers)
	mux.HandleFunc("/deletePlayer", svc.DeletePlayer)

	// CRD Matches
	mux.HandleFunc("/createMatch", svc.CreateMatch)
	mux.HandleFunc("/listMatches", svc.ListMatches)
	mux.HandleFunc("/deleteMatch", svc.DeleteMatch)

	// gameplay
	mux.HandleFunc("/playerThrow", svc.PlayerThrow)

	// misc
	mux.HandleFunc("/statistics", svc.Statistics)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(mux)))
}
