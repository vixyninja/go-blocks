package main

import (
	"encoding/json"
	"net/http"

	"github.com/vixyninja/go-blocks/chi"
	"github.com/vixyninja/go-blocks/logx"
)

func main() {
	server := chi.NewServer(
		chi.WithPrintRoutes(true),
		chi.WithLogger(logx.NewDefaultLogger()),
		chi.WithAddr(":5098"),
	)

	server.Router().Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	})

	server.Run()

}
