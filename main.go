package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/meilisearch/meilisearch-go"
)

type SearchQuery struct {
	Query      string
	MaxResults int64
	Offset     int64
}

type SearchResponse struct {
	Success       bool
	SearchResults []interface{}
	NumResults    int
	TotalResults  int64
	MoreResults   bool
	OriginalQuery SearchQuery
	Offset        int64
	LastOffset    int64
	NumPages      int
}

func pageCount(total int, perPage int) int {
	return int(math.Ceil(float64(total) / float64(perPage)))
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

	http.HandleFunc("/classic/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("templates/search.gtpl")
		t.Execute(w, nil)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("templates/home.html")
		t.Execute(w, nil)
	})

	funcMap := template.FuncMap{
		"unescapeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	//searchResTemplate := template.Must(template.ParseFiles("results.gtpl"))
	searchResTemplate := template.Must(template.New("results.gtpl").Funcs(funcMap).ParseFiles("templates/results.gtpl"))

	// v2.0 UI
	searchResultsTemplate2 := template.Must(template.New("search_results.html").Funcs(funcMap).ParseFiles("templates/search_results.html"))

	if err != nil {
		fmt.Println("Error parsing template")
		os.Exit(1)
	}

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fmt.Println(r.Form)
		if r.Form["query"] != nil || r.FormValue("query") != "" {
			fmt.Println("query:", r.Form["query"])
			var myOffset int64
			if r.Form["offset"] != nil {
				offset, _ := strconv.ParseInt(r.FormValue("offset"), 10, 64)
				myOffset = offset
				if offset < 0 {
					myOffset = 0
				}
			} else {
				offset := int64(0)
				myOffset = offset
			}
			query := SearchQuery{
				Query:      r.FormValue("query"),
				MaxResults: 100,
				Offset:     myOffset,
			}

			res, err := index.Search(query.Query, &meilisearch.SearchRequest{
				Limit:  query.MaxResults,
				Offset: query.Offset,
				AttributesToRetrieve: []string{
					"title",
					"applicant",
					"submission_date",
					"predicates",
					"id",
				},
				AttributesToCrop:      []string{"full_text"},
				AttributesToHighlight: []string{"full_text"},
				HighlightPreTag:       "<mark>",
				HighlightPostTag:      "</mark>",
			})

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			numPages := pageCount(int(res.EstimatedTotalHits), int(query.MaxResults))

			searchResultsTemplate2.Execute(w, SearchResponse{
				Success:       true,
				SearchResults: res.Hits,
				NumResults:    len(res.Hits) + int(query.Offset),
				TotalResults:  res.EstimatedTotalHits,
				MoreResults:   res.EstimatedTotalHits > query.MaxResults,
				OriginalQuery: query,
				Offset:        query.Offset + query.MaxResults,
				LastOffset:    query.Offset - query.MaxResults,
				NumPages:      numPages,
			})
		} else {
			fmt.Println("query is empty")
		}
	})

	http.HandleFunc("/classic/search", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fmt.Println(r.Form)
		if r.Form["query"] != nil || r.FormValue("query") != "" {
			fmt.Println("query:", r.Form["query"])
			var myOffset int64
			if r.Form["offset"] != nil {
				offset, _ := strconv.ParseInt(r.FormValue("offset"), 10, 64)
				myOffset = offset
				if offset < 0 {
					myOffset = 0
				}
			} else {
				offset := int64(0)
				myOffset = offset
			}
			query := SearchQuery{
				Query:      r.FormValue("query"),
				MaxResults: 100,
				Offset:     myOffset,
			}

			res, err := index.Search(query.Query, &meilisearch.SearchRequest{
				Limit:  query.MaxResults,
				Offset: query.Offset,
				AttributesToRetrieve: []string{
					"title",
					"applicant",
					"submission_date",
					"predicates",
					"id",
				},
				AttributesToCrop:      []string{"full_text"},
				AttributesToHighlight: []string{"full_text"},
				HighlightPreTag:       "<mark>",
				HighlightPostTag:      "</mark>",
			})

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			numPages := pageCount(int(res.EstimatedTotalHits), int(query.MaxResults))

			searchResTemplate.Execute(w, SearchResponse{
				Success:       true,
				SearchResults: res.Hits,
				NumResults:    len(res.Hits) + int(query.Offset),
				TotalResults:  res.EstimatedTotalHits,
				MoreResults:   res.EstimatedTotalHits > query.MaxResults,
				OriginalQuery: query,
				Offset:        query.Offset + query.MaxResults,
				LastOffset:    query.Offset - query.MaxResults,
				NumPages:      numPages,
			})
		} else {
			fmt.Println("query is empty")
		}
	})

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
