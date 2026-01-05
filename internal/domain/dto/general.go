package dto

type OrderDirection string

const (
	OrderAsc  OrderDirection = "asc"
	OrderDesc OrderDirection = "desc"
)

type Pagination struct {
	Limit   *uint64
	Offset  *uint64
	OrderBy *string
	Order   *OrderDirection
}
