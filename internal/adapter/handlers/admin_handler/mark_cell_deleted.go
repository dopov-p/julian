package admin_handler

import (
	"context"

	"github.com/dopov-p/julian/internal/pb/admin"
)

// MarkCellDeleted marks a cell as deleted
func (s *Service) MarkCellDeleted(ctx context.Context, req *admin.MarkCellDeletedRequest) (*admin.MarkCellDeletedResponse, error) {
	// TODO: Implement MarkCellDeleted logic
	_ = req // Temporary to avoid unused variable error
	return &admin.MarkCellDeletedResponse{}, nil
}
