package admin_handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/dopov-p/julian/internal/adapter/handlers/grpc"
	adminPb "github.com/dopov-p/julian/internal/pb/admin/api"
	"github.com/dopov-p/julian/internal/usecase/cell_usecase"
)

func (s *Service) MarkCellDeletedOrActivated(
	ctx context.Context,
	req *adminPb.MarkCellDeletedOrActivatedRequest,
) (*adminPb.MarkCellDeletedOrActivatedResponse, error) {
	ucReq, err := validateMarkCellDeletedOrActivatedRequest(req)
	if err != nil {
		return &adminPb.MarkCellDeletedOrActivatedResponse{}, err
	}

	err = s.cellUseCase.MarkDeletedOrActivated(ctx, ucReq)
	if err != nil {
		return &adminPb.MarkCellDeletedOrActivatedResponse{}, grpc.InternalError(
			fmt.Errorf("cellUseCase.MarkDeletedOrActivated: %w", err),
		)
	}

	return &adminPb.MarkCellDeletedOrActivatedResponse{}, nil
}

func validateMarkCellDeletedOrActivatedRequest(
	req *adminPb.MarkCellDeletedOrActivatedRequest,
) (*cell_usecase.MarkDeletedOrActivatedReq, error) {
	var validationErrors []grpc.ValidationErrorItem

	if req == nil || req.GetData() == nil {
		return nil, grpc.ValidationError([]grpc.ValidationErrorItem{
			grpc.NewValidationErrorItem("data", "data is required"),
		})
	}

	data := req.GetData()
	name := strings.TrimSpace(data.GetName())
	if name == "" {
		validationErrors = append(
			validationErrors,
			grpc.NewValidationErrorItem("name", "name is required"),
		)
	}

	if len(validationErrors) > 0 {
		return nil, grpc.ValidationError(validationErrors)
	}

	return &cell_usecase.MarkDeletedOrActivatedReq{
		Name:       name,
		IsActivate: req.GetData().GetIsActive(),
	}, nil
}
