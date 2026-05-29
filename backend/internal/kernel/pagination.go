package kernel

// Pagination holds standard pagination params.
type Pagination struct {
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

// PaginatedResult wraps a list result with total count.
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

// PaginationOptions is the manifesto-style pagination parameter type.
type PaginationOptions struct {
	Page     int
	PageSize int
}

func (p PaginationOptions) Offset() int {
	pg := p.normalized()
	return (pg.Page - 1) * pg.PageSize
}

func (p PaginationOptions) Limit() int {
	return p.normalized().PageSize
}

func (p PaginationOptions) normalized() PaginationOptions {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p
}

// Paginated is the manifesto-style paginated result type.
type Paginated[T any] struct {
	Items      []T `json:"items"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}

// NewPaginated builds a Paginated result from items, pagination params, and total count.
func NewPaginated[T any](items []T, page, pageSize, total int) Paginated[T] {
	if page < 1 {
		page = 1
	}
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
