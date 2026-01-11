package cell_handler

import (
	"context"
	"strings"

	"github.com/dopov-p/julian/internal/adapter/handlers/grpc"
	"github.com/dopov-p/julian/internal/domain/model"
	cellPb "github.com/dopov-p/julian/internal/pb/cell/api"
	"github.com/dopov-p/julian/internal/usecase/cell_usecase"
	"github.com/samber/lo"
)

func (s *Service) DevastateCellContents(
	ctx context.Context,
	req *cellPb.DevastateCellContentsRequest,
) (*cellPb.DevastateCellContentsResponse, error) {
	ucReq, err := validateDevastateCellContentsRequest(req)
	if err != nil {
		return &cellPb.DevastateCellContentsResponse{}, err
	}

	err = s.cellUseCase.DevastateContents(ctx, ucReq)
	if err != nil {
		return nil, err
	}

	return &cellPb.DevastateCellContentsResponse{}, nil
}

func validateDevastateCellContentsRequest(
	req *cellPb.DevastateCellContentsRequest,
) (*cell_usecase.DevastateContentsReq, error) {
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
			grpc.NewValidationErrorItem("data.name", "name is required and cannot be empty"),
		)
	}

	contents := data.GetCellContents()
	if len(contents) == 0 {
		validationErrors = append(
			validationErrors,
			grpc.NewValidationErrorItem("data.cell_contents", "cell_contents is required"),
		)
	}

	if len(validationErrors) > 0 {
		return nil, grpc.ValidationError(validationErrors)
	}

	return &cell_usecase.DevastateContentsReq{
		Name: name,
		Contents: lo.Map(
			contents,
			func(protoContent *cellPb.DevastateCellContentsRequest_Data_CellContents, _ int) model.CellContents {
				return model.CellContents{
					ExternalOrderID: protoContent.ExternalOrderId,
					SKU:             protoContent.GetSku(),
					Quantity:        protoContent.GetQuantity(),
				}
			},
		),
	}, nil
}
