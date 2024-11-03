package dto

import (
	"fmt"
	"math"
	"strconv"
)

type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Links Links       `json:"links"`
	Meta  Meta        `json:"meta"`
}

type Links struct {
	First string      `json:"first"`
	Last  string      `json:"last"`
	Prev  interface{} `json:"prev"` // interface{} to allow null
	Next  interface{} `json:"next"` // interface{} to allow null
}

type Meta struct {
	CurrentPage int        `json:"current_page"`
	From        int        `json:"from"`
	LastPage    int        `json:"last_page"`
	Links       []PageLink `json:"links"`
	Path        string     `json:"path"`
	PerPage     int        `json:"per_page"`
	To          int        `json:"to"`
	Total       int        `json:"total"`
}

type PageLink struct {
	URL    interface{} `json:"url"` // interface{} to allow null
	Label  string      `json:"label"`
	Active bool        `json:"active"`
}

type DataDBFilter struct {
	OrderBy  string
	OrderDir string
	Page     int
	PerPage  int
}

func NewPaginatedResponse(data interface{}, page, perPage, total int, baseURL string) PaginatedResponse {
	lastPage := int(math.Ceil(float64(total) / float64(perPage)))

	// Calculate from and to
	from := ((page - 1) * perPage) + 1
	to := from + perPage - 1
	if to > total {
		to = total
	}
	if total == 0 {
		from = 0
	}

	// Create base response
	response := PaginatedResponse{
		Data: data,
		Links: Links{
			First: fmt.Sprintf("%s?page=1", baseURL),
			Last:  fmt.Sprintf("%s?page=%d", baseURL, lastPage),
			Prev:  nil,
			Next:  nil,
		},
		Meta: Meta{
			CurrentPage: page,
			From:        from,
			LastPage:    lastPage,
			Path:        baseURL,
			PerPage:     perPage,
			To:          to,
			Total:       total,
		},
	}

	// Set prev/next links
	if page > 1 {
		response.Links.Prev = fmt.Sprintf("%s?page=%d", baseURL, page-1)
	}
	if page < lastPage {
		response.Links.Next = fmt.Sprintf("%s?page=%d", baseURL, page+1)
	}

	// Generate pagination links
	response.Meta.Links = make([]PageLink, 0)

	// Previous link
	response.Meta.Links = append(response.Meta.Links, PageLink{
		URL:    response.Links.Prev,
		Label:  "&laquo; Previous",
		Active: false,
	})

	// Number links
	for i := 1; i <= lastPage; i++ {
		response.Meta.Links = append(response.Meta.Links, PageLink{
			URL:    fmt.Sprintf("%s?page=%d", baseURL, i),
			Label:  strconv.Itoa(i),
			Active: i == page,
		})
	}

	// Next link
	response.Meta.Links = append(response.Meta.Links, PageLink{
		URL:    response.Links.Next,
		Label:  "Next &raquo;",
		Active: false,
	})

	return response
}
