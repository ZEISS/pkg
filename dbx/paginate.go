package dbx

import (
	"fmt"
	"math"

	"github.com/zeiss/pkg/utilx"

	"gorm.io/gorm"
)

// DefaultLimits is a list of default limits.
var DefaultLimits = []int{5, 10, 25, 50}

// NewPaginated returns a new Paginated struct.
func NewPaginated[T any]() *Paginated[T] {
	return &Paginated[T]{
		Limit:  10,
		Offset: 0,
		Sort:   SortDesc,
	}
}

const (
	// SortNone is the no sort order.
	SortNone = ""
	// SortAsc is the ascending sort order.
	SortAsc = "asc"
	// SortDesc is the descending sort order.
	SortDesc = "desc"
)

// Paginated is a struct that contains the properties of a pagination.
type Paginated[T any] struct {
	// Limit is the number of items to return.
	Limit int `json:"limit" xml:"limit" form:"limit" query:"limit"`
	// Offset is the number of items to skip.
	Offset int `json:"offset" xml:"offset" form:"offset" query:"offset"`
	// Search is the search term to filter the results.
	Search string `json:"search,omitempty" xml:"search" form:"search" query:"search"`
	// Sort is the sorting order.
	Sort string `json:"sort,omitempty" xml:"sort" form:"sort" query:"sort"`
	// Value is the value to paginate.
	Value T `json:"value,omitempty" xml:"value" form:"value" query:"value"`
}

// GetLimit returns the limit.
func (p *Paginated[T]) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}

	return p.Limit
}

// GetOffset returns the page.
func (p *Paginated[T]) GetOffset() int {
	if p.Offset < 0 {
		p.Offset = 0
	}

	return p.Offset
}

// GetSort returns the sort.
func (p *Paginated[T]) GetSort() string {
	if p.Sort == SortNone {
		p.Sort = SortDesc
	}

	return p.Sort
}

// GetSearch returns the search.
func (p *Paginated[T]) GetSearch() string {
	return p.Search
}

// Paginate returns a function that paginates the results.
func Paginate[T any](_ interface{}, pagination *Paginated[T], _ *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

// RowsPtr is a function that returns the rows as pointers.
func RowsPtr[T any](rows []T) []*T {
	rowsPtr := make([]*T, 0, len(rows))
	for _, row := range rows {
		rowsPtr = append(rowsPtr, &row)
	}

	return rowsPtr
}

// NewResults returns a new Results struct.
func NewResults[T any]() *Results[T] {
	return &Results[T]{
		Limit:  10,
		Offset: 0,
		Sort:   SortDesc,
	}
}

// Results is a struct that contains the results of a query.
type Results[T any] struct {
	// Limit is the number of items to return.
	Limit int `json:"limit" xml:"limit" form:"limit" query:"limit"`
	// Offset is the number of items to skip.
	Offset int `json:"offset" xml:"offset" form:"offset" query:"offset"`
	// Search is the search term to filter the results.
	Search string `json:"search,omitempty" xml:"search" form:"search" query:"search"`
	// SearchFields is the search term to filter the results.
	SearchFields []string `json:"-"`
	// Sort is the sorting order.
	Sort string `json:"sort,omitempty" xml:"sort" form:"sort" query:"sort"`
	// TotalRows is the total number of rows.
	TotalRows int `json:"total_rows"`
	// TotalPages is the total number of pages.
	TotalPages int `json:"total_pages"`
	// Rows is the items to return.
	Rows []T `json:"rows" xml:"rows"`
}

// GetLimit returns the limit.
func (p *Results[T]) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}

	return p.Limit
}

// GetOffset returns the page.
func (p *Results[T]) GetOffset() int {
	if p.Offset < 0 {
		p.Offset = 0
	}

	return p.Offset
}

// GetSort returns the sort.
func (p *Results[T]) GetSort() string {
	if p.Sort == SortNone {
		p.Sort = SortDesc
	}

	return p.Sort
}

// GetSearch returns the search.
func (p *Results[T]) GetSearch() string {
	return p.Search
}

// GetRows returns the rows as pointers.
func (p *Results[T]) GetRows() []*T {
	rows := make([]*T, 0, len(p.Rows))
	for _, row := range p.Rows {
		rows = append(rows, &row)
	}

	return rows
}

// GetTotalRows returns the total rows.
func (p *Results[T]) GetTotalRows() int {
	return p.TotalRows
}

// GetTotalPages returns the total pages.
func (p *Results[T]) GetTotalPages() int {
	return p.TotalPages
}

// GetLen returns the length of the rows.
func (p *Results[T]) GetLen() int {
	return len(p.Rows)
}

// PaginatedResults returns a function that paginates the results.
func PaginatedResults[T any](value interface{}, pagination *Results[T], db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Scopes(SearchScope(pagination)).Count(&totalRows).Where("deleted_at IS NULL")

	pagination.TotalRows = int(totalRows)
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		db = db.Scopes(SearchScope(pagination))
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

// SearchScope is a function that returns a scope that searches the given fields.
func SearchScope[T any](pagination *Results[T]) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if utilx.Empty(pagination.GetSearch()) {
			return db
		}

		for _, field := range pagination.SearchFields {
			db = db.Where(fmt.Sprintf("%s LIKE ? ", field), fmt.Sprintf("%%%s%%", pagination.GetSearch()))
		}

		return db
	}
}
