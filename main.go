package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go-nist-api/cpe"
)

var tmpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

func searchHandler(cpeApi *cpe.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()
		keyword := params.Get("q")

		response, err := cpeApi.FetchAll(keyword)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("%+v", response)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("error during loading the .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	apiKey := os.Getenv("NEWS_API_KEY")
	// if apiKey == "" {
	// 	log.Fatal("Env: apiKey must be set")
	// }

	myClient := &http.Client{Timeout: 10 * time.Second}
	cpeApi := cpe.NewClient(myClient, apiKey)

	fs := http.FileServer(http.Dir("assets"))

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/search", searchHandler(cpeApi))
	mux.HandleFunc("/", indexHandler)

	http.ListenAndServe(":"+port, mux)

}
