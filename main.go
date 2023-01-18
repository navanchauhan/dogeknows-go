package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/meilisearch/meilisearch-go"
	mytemplates "github.com/navanchauhan/dogeknows-go/templates"
)

var globalVariables = GlobalVars{
	Name: "510K Search",
}

//go:embed static/*
var static embed.FS

func create_pdf_url(year string, knumber string) string {
	year_int, _ := strconv.Atoi(year[0:2])
	if year[0] == '0' || year[0] == '1' || year[0] == '2' || year[0] == '3' || year[0] == '4' || year[0] == '5' && year[0:2] != "01" && year[0:2] != "00" {
		return fmt.Sprintf("https://www.accessdata.fda.gov/cdrh_docs/pdf%d/%s.pdf", year_int, knumber)
	} else {
		return fmt.Sprintf("https://www.accessdata.fda.gov/cdrh_docs/pdf/%s.pdf", knumber)
	}
}

func main() {
	sentry_err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://72a0913f00e24d5da4d44ccde6ebcc9c@o126149.ingest.sentry.io/4504524880674816",
		TracesSampleRate: 1.0,
	})
	if sentry_err != nil {
		log.Fatalf("sentry.Init: %s", sentry_err)
	}
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

	contentStatic, _ := fs.Sub(static, "static")

	funcMap := template.FuncMap{
		"unescapeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	// Classic UI Templates
	classicIndexTemplate := template.Must(template.ParseFS(mytemplates.Templates, "search.gtpl"))
	searchResTemplate := template.Must(template.New("results.gtpl").Funcs(funcMap).ParseFS(mytemplates.Templates, "results.gtpl"))

	// v2.0 UI Templates
	indexTemplate := template.Must(template.ParseFS(mytemplates.Templates, "home.html", "components/section.html", "components/header.html"))
	searchResultsTemplate2 := template.Must(template.New("search_results.html").Funcs(funcMap).ParseFS(mytemplates.Templates, "search_results.html", "components/section.html", "components/header.html"))
	documentDetailsTemplate2 := template.Must(template.New("document_details.html").Funcs(funcMap).ParseFS(mytemplates.Templates, "document_details.html", "components/section.html", "components/header.html"))

	index := client.Index("510k")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(contentStatic))))

	http.HandleFunc("/classic/", func(w http.ResponseWriter, r *http.Request) {
		classicIndexTemplate.Execute(w, BaseResponse{
			GlobalVars: globalVariables,
		})
	})

	http.HandleFunc("/classic/search", func(w http.ResponseWriter, r *http.Request) {
		searchHandler(w, r, index, searchResTemplate)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexTemplate.Execute(w, BaseResponse{
			GlobalVars: globalVariables,
		})
	})

	http.HandleFunc("/dbentry", func(w http.ResponseWriter, r *http.Request) {
		documentHandler510k(w, r, index, documentDetailsTemplate2)
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		searchHandler(w, r, index, searchResultsTemplate2)
	})

	fmt.Println("Listening on port 8752")
	http.ListenAndServe(":8752", nil)
}
