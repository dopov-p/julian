package cell_handler

import (
	"context"

	"github.com/dopov-p/julian/internal/pb/cell"
)

// FullDevastateCellContents fully devastates cell contents
func (s *Service) FullDevastateCellContents(ctx context.Context, req *cell.FullDevastateCellContentsRequest) (*cell.FullDevastateCellContentsResponse, error) {
	// TODO: Implement FullDevastateCellContents logic
	_ = req // Temporary to avoid unused variable error
	return &cell.FullDevastateCellContentsResponse{}, nil
}
