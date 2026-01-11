package admin_handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/dopov-p/julian/internal/adapter/handlers/grpc"
	"github.com/dopov-p/julian/internal/domain/model"
	adminPb "github.com/dopov-p/julian/internal/pb/admin/api"
	"github.com/dopov-p/julian/internal/usecase/cell_usecase"
)

func (s *Service) CreateCell(ctx context.Context, req *adminPb.CreateCellRequest) (*adminPb.CreateCellResponse, error) {
	cell, err := s.validateCreateCellRequest(req)
	if err != nil {
		return nil, err
	}

	err = s.cellUseCase.Create(ctx, cell)
	if err != nil {
		return nil, grpc.InternalError(fmt.Sprintf("failed to create cell: %v", err))
	}

	return &adminPb.CreateCellResponse{}, nil
}

func (s *Service) validateCreateCellRequest(req *adminPb.CreateCellRequest) (*cell_usecase.CreateReq, error) {
	var validationErrors []grpc.ValidationErrorItem

	if req == nil || req.Data == nil {
		return nil, grpc.ValidationError([]grpc.ValidationErrorItem{
			grpc.NewValidationErrorItem("data", "data is required"),
		})
	}

	data := req.Data
	if data.Cell == nil {
		return nil, grpc.ValidationError([]grpc.ValidationErrorItem{
			grpc.NewValidationErrorItem("data.cell", "data.cell is required"),
		})
	}

	protoCell := data.Cell
	name := strings.TrimSpace(protoCell.GetName())
	if name == "" {
		validationErrors = append(validationErrors, grpc.NewValidationErrorItem("data.cell.name", "name is required and cannot be empty"))
	}

	kind := protoCell.GetKind()
	if kind == adminPb.CellKind_CELL_KIND_UNSPECIFIED {
		validationErrors = append(validationErrors, grpc.NewValidationErrorItem("data.cell.kind", "kind is required and cannot be UNSPECIFIED"))
	}

	contentType := protoCell.GetContentType()
	if contentType == adminPb.CellContentType_CELL_CONTENT_TYPE_UNSPECIFIED {
		validationErrors = append(validationErrors, grpc.NewValidationErrorItem("data.cell.content_type", "content_type is required and cannot be UNSPECIFIED"))
	}

	if len(validationErrors) > 0 {
		return nil, grpc.ValidationError(validationErrors)
	}

	return &cell_usecase.CreateReq{
		Name:            name,
		CellKind:        model.CellKind(kind),               //nolint:gosec
		CellContentType: model.CellContentType(contentType), //nolint:gosec
		CanHasFewOrders: protoCell.GetCanHasFewOrders(),
	}, nil
}
