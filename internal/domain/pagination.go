package domain

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// Pagination holds normalized page/size parameters for list queries.
type Pagination struct {
	Page     int
	PageSize int
}

// Offset returns the SQL OFFSET value for the current page.
func (p Pagination) Offset() int {
	if p.Page <= 1 {
		return 0
	}
	return (p.Page - 1) * p.PageSize
}

// Normalize clamps page/pageSize to valid values.
func (p Pagination) Normalize() Pagination {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = DefaultPageSize
	}
	if p.PageSize > MaxPageSize {
		p.PageSize = MaxPageSize
	}
	return p
}

// TotalPages returns the number of pages for a given total record count.
func TotalPages(total int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	pages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		pages++
	}
	return pages
}
