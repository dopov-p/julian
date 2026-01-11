package cell_handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/dopov-p/julian/internal/adapter/handlers/grpc"
	cellPb "github.com/dopov-p/julian/internal/pb/cell/api"
)

func (s *Service) FullDevastateCellContents(
	ctx context.Context,
	req *cellPb.FullDevastateCellContentsRequest,
) (*cellPb.FullDevastateCellContentsResponse, error) {
	name, err := validateFullDevastateCellContentsRequest(req)
	if err != nil {
		return &cellPb.FullDevastateCellContentsResponse{}, err
	}

	err = s.cellRepo.FullDevastate(ctx, name)
	if err != nil {
		return &cellPb.FullDevastateCellContentsResponse{}, grpc.InternalError(
			fmt.Errorf("cellRepo.FullDevastate: %w", err),
		)
	}

	return &cellPb.FullDevastateCellContentsResponse{}, nil
}

func validateFullDevastateCellContentsRequest(
	req *cellPb.FullDevastateCellContentsRequest,
) (string, error) {
	var validationErrors []grpc.ValidationErrorItem

	if req == nil || req.GetData() == nil {
		return "", grpc.ValidationError([]grpc.ValidationErrorItem{
			grpc.NewValidationErrorItem("data", "data is required"),
		})
	}

	name := strings.TrimSpace(req.GetData().GetName())
	if name == "" {
		validationErrors = append(
			validationErrors,
			grpc.NewValidationErrorItem("name", "name is required"),
		)
	}

	if len(validationErrors) > 0 {
		return "", grpc.ValidationError(validationErrors)
	}

	return name, nil
}
