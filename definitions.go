package main

type SearchQuery struct {
	Query      string
	MaxResults int64
	Page       int64
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
	MaxResults    int64
	CurPage       int64
	ShowPrev      bool
	PrevPage      int64
	NextPage      int64
	Sort          string
}

type DocumentResponse struct {
	GlobalVars    GlobalVars
	SearchResults interface{}
	SummaryPDF    string
}

type GlobalVars struct {
	Name string
}
