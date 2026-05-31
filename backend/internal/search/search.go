package search

import (
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// SearchQuery holds all parameters for a product search operation.
type SearchQuery struct {
	Query      string   `json:"query"`
	CategoryID string   `json:"category_id,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	Status     string   `json:"status,omitempty"`
	MinPrice   *int64   `json:"min_price,omitempty"`
	MaxPrice   *int64   `json:"max_price,omitempty"`
	SortBy     string   `json:"sort_by,omitempty"` // "relevance", "price_asc", "price_desc", "name", "created_at"
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
}

// SearchResult is the paginated result of a product search.
type SearchResult struct {
	Products   []ProductHit `json:"products"`
	Total      int          `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
	Query      string       `json:"query"`
}

// ProductHit is a product returned from a search.
type ProductHit struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Price       kernel.Money `json:"price"`
	SKU         string       `json:"sku"`
	Images      []string     `json:"images"`
	CategoryID  string       `json:"category_id"`
	Tags        []string     `json:"tags"`
	Status      string       `json:"status"`
	Stock       int          `json:"stock"`
	CreatedAt   time.Time    `json:"created_at"`
	Rank        float64      `json:"rank,omitempty"`
}

// SearchSuggestion is an autocomplete suggestion for a search term prefix.
type SearchSuggestion struct {
	Term  string `json:"term"`
	Count int    `json:"count"`
}
