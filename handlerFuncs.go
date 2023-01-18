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
		var pageNo int64
		var maxHits int64
		if r.Form["page"] != nil {
			convertedInt, _ := strconv.ParseInt(r.FormValue("page"), 10, 64)
			if convertedInt > 0 {
				pageNo = convertedInt
			} else {
				pageNo = 1
			}
		} else {
			pageNo = 1
		}
		if r.Form["maxHits"] != nil {
			convertedInt, _ := strconv.ParseInt(r.FormValue("maxHits"), 10, 64)
			if convertedInt > 0 {
				maxHits = convertedInt
			} else {
				maxHits = 20
			}
		} else {
			maxHits = 20
		}

		if maxHits > 100 {
			maxHits = 100
		}
		query := SearchQuery{
			Query:      r.FormValue("query"),
			MaxResults: maxHits,
			Page:       pageNo,
		}

		res, err := index.Search(query.Query, &meilisearch.SearchRequest{
			HitsPerPage: query.MaxResults,
			Page:        query.Page,
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

		searchTemplate.Execute(w, SearchResponse{
			GlobalVars:    globalVariables,
			Success:       true,
			SearchResults: res.Hits,
			NumResults:    int(res.TotalHits) - (int(res.TotalPages)-int(res.Page))*int(res.HitsPerPage),
			TotalResults:  res.TotalHits,
			MoreResults:   res.Page < res.TotalPages,
			OriginalQuery: query,
			NumPages:      int(res.TotalPages),
			MaxResults:    maxHits,
			CurPage:       res.Page,
			ShowPrev:      res.Page > 1,
			PrevPage:      res.Page - 1,
			NextPage:      res.Page + 1,
		})

	} else {
		// Return empty search results
		searchTemplate.Execute(w, SearchResponse{
			GlobalVars:    globalVariables,
			Success:       true,
			SearchResults: []interface{}{},
			NumResults:    0,
			TotalResults:  0,
			MoreResults:   false,
			OriginalQuery: SearchQuery{},
			NumPages:      0,
			MaxResults:    0,
			CurPage:       0,
			ShowPrev:      false,
			PrevPage:      0,
			NextPage:      0,
		})
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
