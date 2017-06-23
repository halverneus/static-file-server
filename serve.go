package main

import (
	"net/http"
	"os"
)

func main() {
	host := env("HOST", "")
	port := env("PORT", "8080")
	folder := env("FOLDER", "/web") + "/"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, folder+r.URL.Path)
	})
	http.ListenAndServe(host+":"+port, nil)
}

func env(key, fallback string) string {
	if value := os.Getenv(key); 0 < len(value) {
		return value
	}
	return fallback
}
