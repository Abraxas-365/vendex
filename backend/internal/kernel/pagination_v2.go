package kernel

// PaginationOptions holds pagination parameters for list queries.
// It is the manifesto-compatible replacement for the older Pagination type.
type PaginationOptions struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// page returns a sanitised page number (min 1).
func (p PaginationOptions) page() int {
	if p.Page < 1 {
		return 1
	}
	return p.Page
}

// pageSize returns a sanitised page size (1–100, default 20).
func (p PaginationOptions) pageSize() int {
	if p.PageSize < 1 {
		return 20
	}
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}

// Offset returns the SQL OFFSET value for this page.
func (p PaginationOptions) Offset() int {
	return (p.page() - 1) * p.pageSize()
}

// Limit returns the SQL LIMIT value for this page.
func (p PaginationOptions) Limit() int {
	return p.pageSize()
}

// Paginated is a generic paginated list result.
// It is the manifesto-compatible replacement for the older PaginatedResult type.
type Paginated[T any] struct {
	Items      []T `json:"items"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}

// NewPaginated constructs a Paginated result from a slice of items and pagination metadata.
func NewPaginated[T any](items []T, page, pageSize, total int) Paginated[T] {
	opts := PaginationOptions{Page: page, PageSize: pageSize}
	ps := opts.pageSize()
	pg := opts.page()

	totalPages := total / ps
	if total%ps > 0 {
		totalPages++
	}

	if items == nil {
		items = []T{}
	}

	return Paginated[T]{
		Items:      items,
		Total:      total,
		Page:       pg,
		PageSize:   ps,
		TotalPages: totalPages,
	}
}
