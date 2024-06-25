package main

import (
	"embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:embed client/views/*
var views embed.FS

var t = template.Must(template.ParseFS(views, "client/views/*"))

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/pikachu")
		if err != nil {
			panic(err)
		}
		data := Pokemon{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Fatalf("Failed to decode JSON: %v", err)
		}

		if err := t.ExecuteTemplate(w, "index.html", data); err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("POST /poke", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Fatal(err)
		}

		resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + strings.ToLower(r.FormValue("pokemon")))
		if err != nil {
			panic(err)
		}
		data := Pokemon{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Fatalf("Failed to decode JSON: %v", err)

		}
		if err := t.ExecuteTemplate(w, "response.html", data); err != nil {
			log.Fatal(err)
		}
	})

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
