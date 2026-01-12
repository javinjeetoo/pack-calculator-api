package ui

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed web/*
var content embed.FS

// Handler returns an http.Handler that serves the embedded UI files.
func Handler() (http.Handler, error) {
	sub, err := fs.Sub(content, "web")
	if err != nil {
		return nil, err
	}
	return http.FileServer(http.FS(sub)), nil
}
