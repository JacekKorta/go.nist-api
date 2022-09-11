package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var tmpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
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

	fs := http.FileServer(http.Dir("assets"))

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/", indexHandler)

	http.ListenAndServe(":"+port, mux)


}


