package model

import (
	"errors"
	"time"
)

var (
	ErrCellNotFound  = errors.New("cell not found")
	ErrEmptyContents = errors.New("empty contents")
)

type (
	CellKind        uint32
	CellContentType uint32
)

const (
	CellTypeUnspecified              CellKind = 0
	CellTypeWarehouseStorage         CellKind = 1
	CellTypeWarehouseHandleContainer CellKind = 2
	CellTypeWarehouseBox             CellKind = 3

	CellContentTypeUnspecified       CellContentType = 0
	CellContentTypeExpendableProduct CellContentType = 1
	CellContentTypeMainProduct       CellContentType = 2
	CellContentTypeOversizedProduct  CellContentType = 3
	CellContentTypeExpendableLiquid  CellContentType = 4
	CellContentTypeMainLiquid        CellContentType = 5
)

type CellContents struct {
	ExternalOrderID *string `json:"external_order_id,omitempty"`
	SKU             string  `json:"sku"`
	Quantity        uint64  `json:"quantity"`
}

type Cell struct {
	ID              string
	Name            string
	Kind            CellKind
	ContentType     CellContentType
	Contents        []CellContents
	CanHasFewOrders bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

func NewCell(
	id, name string,
	kind CellKind,
	contentType CellContentType,
	canHasFewOrders bool,
) *Cell {
	return &Cell{
		ID:              id,
		Name:            name,
		Kind:            kind,
		ContentType:     contentType,
		CanHasFewOrders: canHasFewOrders,
	}
}
