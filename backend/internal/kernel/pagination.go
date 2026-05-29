package kernel

// Pagination holds standard pagination params.
// Deprecated: use PaginationOptions instead.
type Pagination struct {
	Page     int
	PageSize int
}

// PaginationOptions holds standard pagination params (manifesto pattern).
type PaginationOptions struct {
	Page     int
	PageSize int
}

func NewPagination(page, pageSize int) Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return Pagination{Page: page, PageSize: pageSize}
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p Pagination) Limit() int {
	return p.PageSize
}

// NewPaginationOptions creates PaginationOptions with safe defaults.
func NewPaginationOptions(page, pageSize int) PaginationOptions {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return PaginationOptions{Page: page, PageSize: pageSize}
}

func (p PaginationOptions) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p PaginationOptions) Limit() int {
	return p.PageSize
}

// PaginatedResult wraps a list result with total count.
// Deprecated: use Paginated[T] instead.
type PaginatedResult[T any] struct {
	Items      []T `json:"items"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}

func NewPaginatedResult[T any](items []T, total int, p Pagination) PaginatedResult[T] {
	totalPages := total / p.PageSize
	if total%p.PageSize > 0 {
		totalPages++
	}
	return PaginatedResult[T]{
		Items:      items,
		Total:      total,
		Page:       p.Page,
		PageSize:   p.PageSize,
		TotalPages: totalPages,
	}
}

// Paginated wraps a list result with total count (manifesto pattern).
type Paginated[T any] struct {
	Items      []T `json:"items"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}

// NewPaginated creates a Paginated[T] result.
func NewPaginated[T any](items []T, page, pageSize, total int) Paginated[T] {
	if pageSize < 1 {
		pageSize = 20
	}
	totalPages := total / pageSize
	if total%pageSize > 0 {
		totalPages++
	}
	return Paginated[T]{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
