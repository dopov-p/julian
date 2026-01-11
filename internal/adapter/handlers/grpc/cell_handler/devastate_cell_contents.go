package cell_handler

import (
	"context"

	"github.com/dopov-p/julian/internal/domain/model"
	cell "github.com/dopov-p/julian/internal/pb/cell/api"
	"github.com/dopov-p/julian/internal/usecase/cell_usecase"
)

// DevastateCellContents partially devastates cell contents.
func (s *Service) DevastateCellContents(ctx context.Context, req *cell.DevastateCellContentsRequest) (*cell.DevastateCellContentsResponse, error) {
	if req.Data == nil {
		return &cell.DevastateCellContentsResponse{}, nil
	}

	// Convert proto CellContents to domain model
	var contents []model.CellContents
	if req.Data.CellContents != nil {
		contents = []model.CellContents{
			{
				ExternalOrderID: req.Data.CellContents.ExternalOrderId, // Already *string from proto optional field
				SKU:             req.Data.CellContents.Sku,
				Quantity:        int64(req.Data.CellContents.Quantity), //nolint:gosec
			},
		}
	}

	err := s.cellUseCase.DevastateContents(ctx, cell_usecase.DevastateContentsReq{
		Name:     req.Data.Name,
		Contents: contents,
	})
	if err != nil {
		return nil, err
	}

	return &cell.DevastateCellContentsResponse{}, nil
}
