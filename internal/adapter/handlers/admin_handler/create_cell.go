package admin_handler

import (
	"context"

	"github.com/dopov-p/julian/internal/pb/admin"
)

// CreateCell creates a new cell
func (s *Service) CreateCell(ctx context.Context, req *admin.CreateCellRequest) (*admin.CreateCellResponse, error) {
	// TODO: Implement CreateCell logic
	_ = req // Temporary to avoid unused variable error
	return &admin.CreateCellResponse{}, nil
}
