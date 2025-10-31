package main

import (
	"encoding/json"
	"net/http"

	"github.com/vixyninja/go-blocks/pkg/chix"
	"github.com/vixyninja/go-blocks/pkg/logx"
)

func main() {
	server := chix.NewServer(
		chix.WithPrintRoutes(true),
		chix.WithLogger(logx.NewDefaultLogger()),
		chix.WithAddr(":5098"),
	)

	server.Router().Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	})

	server.Run()

}
