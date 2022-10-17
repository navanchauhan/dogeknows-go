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

type BaseResponse struct {
	GlobalVars GlobalVars
}

type SearchResponse struct {
	GlobalVars    GlobalVars
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

type DocumentResponse struct {
	GlobalVars    GlobalVars
	SearchResults interface{}
	SummaryPDF    string
}

type GlobalVars struct {
	Name string
}

var globalVariables = GlobalVars{
	Name: "510K Search",
}

func create_pdf_url(year string, knumber string) string {
	year_int, _ := strconv.Atoi(year[0:2])
	if year[0] == '0' || year[0] == '1' || year[0] == '2' || year[0] == '3' || year[0] == '4' || year[0] == '5' && year[0:2] != "01" && year[0:2] != "00" {
		return fmt.Sprintf("https://www.accessdata.fda.gov/cdrh_docs/pdf%d/%s.pdf", year_int, knumber)
	} else {
		return fmt.Sprintf("https://www.accessdata.fda.gov/cdrh_docs/pdf/%s.pdf", knumber)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request, index *meilisearch.Index, searchTemplate *template.Template) {
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

		searchTemplate.Execute(w, SearchResponse{
			GlobalVars:    globalVariables,
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

	index := client.Index("510k")

	http.HandleFunc("/classic/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("templates/search.gtpl")
		t.Execute(w, BaseResponse{
			GlobalVars: globalVariables,
		})
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("templates/home.html")
		t.Execute(w, BaseResponse{
			GlobalVars: globalVariables,
		})
	})

	funcMap := template.FuncMap{
		"unescapeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	// Classic UI
	searchResTemplate := template.Must(template.New("results.gtpl").Funcs(funcMap).ParseFiles("templates/results.gtpl"))

	// v2.0 UI
	searchResultsTemplate2 := template.Must(template.New("search_results.html").Funcs(funcMap).ParseFiles("templates/search_results.html"))
	documentDetailsTemplate2 := template.Must(template.ParseFiles("templates/document_details.html"))

	http.HandleFunc("/dbentry", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fmt.Println(r.Form)
		var res interface{}
		var documentID string = r.FormValue("id")
		if r.Form["id"] != nil {
			index.GetDocument(documentID, &meilisearch.DocumentQuery{
				Fields: []string{
					"title",
					"applicant",
					"decision",
					"decision_date",
					"full_text",
					"id",
					"predicates",
					"submission_date",
					"contact",
					"STREET1",
					"STREET2",
					"CITY",
					"STATE",
					"ZIP",
					"COUNTRY_CODE",
					"postal_code",
					"REVIEWADVISECOMM",
					"PRODUCTCODE",
					"STATEORSUMM",
					"CLASSADVISECOMM",
					"SSPINDICATOR",
					"TYPE",
					"THIRDPARTY",
					"EXPEDITEDREVIEW",
				}}, &res)
			fmt.Println(res)
			var year = documentID[1:3]

			documentDetailsTemplate2.Execute(w, DocumentResponse{
				GlobalVars:    globalVariables,
				SearchResults: res,
				SummaryPDF:    create_pdf_url(year, documentID),
			})
		} else {
			fmt.Println("No ID provided")
		}

	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		searchHandler(w, r, index, searchResultsTemplate2)
	})

	http.HandleFunc("/classic/search", func(w http.ResponseWriter, r *http.Request) {
		searchHandler(w, r, index, searchResTemplate)
	})

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
