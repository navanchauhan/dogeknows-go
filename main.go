package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/meilisearch/meilisearch-go"
)

type SearchQuery struct {
	Query      string
	MaxResults int64
}

type SearchResponse struct {
	Success       bool
	SearchResults []interface{}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}
	meili_host, ok := os.LookupEnv("MEILI_HOST")
	if !ok {
		fmt.Println("Error loading MEILI_HOST from .env file")
		os.Exit(1)
	}

	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host: meili_host,
	})

	index := client.Index("fda510k")

	tmpl := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}
		query := SearchQuery{
			Query:      r.FormValue("query"),
			MaxResults: 100,
		}

		res, err := index.Search(query.Query, &meilisearch.SearchRequest{
			Limit: query.MaxResults,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println(res.Hits)

		tmpl.Execute(w, SearchResponse{Success: true, SearchResults: res.Hits})
	})

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
