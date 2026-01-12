package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/javinjeetoo/pack-calculator-api/internal/api"
	"github.com/javinjeetoo/pack-calculator-api/internal/ui"
)

func main() {
	defaultPackSizes := []int{250, 500, 1000, 2000, 5000}
	if env := os.Getenv("PACK_SIZES"); env != "" {
		if parsed, err := parsePackSizes(env); err == nil && len(parsed) > 0 {
			defaultPackSizes = parsed
		} else if err != nil {
			log.Printf("WARN: invalid PACK_SIZES, using defaults: %v", err)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	h := api.NewHandler(defaultPackSizes)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/calculate", h.Calculate)

	// UI served at /
	uiHandler, err := ui.Handler()
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", uiHandler)

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func parsePackSizes(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, nil
}
