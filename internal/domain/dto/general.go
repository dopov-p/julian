package dto

type OrderDirection string

const (
	OrderAsc  OrderDirection = "ASC"
	OrderDesc OrderDirection = "DESC"
)

type Pagination struct {
	Limit  *uint64
	Offset *uint64
}

type Sorting struct {
	OrderBy *string
	Order   *OrderDirection
}
