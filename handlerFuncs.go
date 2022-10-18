package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/meilisearch/meilisearch-go"
)

func searchHandler(w http.ResponseWriter, r *http.Request, index *meilisearch.Index, searchTemplate *template.Template) {
	r.ParseForm()
	fmt.Println(r.Form)
	if r.Form["query"] != nil || r.FormValue("query") != "" {
		//fmt.Println("query:", r.Form["query"])
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

func documentHandler510k(w http.ResponseWriter, r *http.Request, index *meilisearch.Index, template *template.Template) {
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
		//fmt.Println(res)
		var year = documentID[1:3]

		template.Execute(w, DocumentResponse{
			GlobalVars:    globalVariables,
			SearchResults: res,
			SummaryPDF:    create_pdf_url(year, documentID),
		})
	} else {
		fmt.Println("No ID provided")
	}

}
