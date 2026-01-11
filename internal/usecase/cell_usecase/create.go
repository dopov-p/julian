package cell_usecase

import (
	"context"
	"fmt"

	"github.com/dopov-p/julian/internal/domain/model"
)

type CreateReq struct {
	Name            string
	CellKind        model.CellKind
	CellContentType model.CellContentType
	CanHasFewOrders bool
}

func (u *UseCase) Create(ctx context.Context, req *CreateReq) error {
	id, err := u.generator.GenerateUUID()
	if err != nil {
		return fmt.Errorf("generator.GenerateUUID: %w", err)
	}

	cell := model.NewCell(id, req.Name, req.CellKind, req.CellContentType, req.CanHasFewOrders)
	err = u.cellRepo.Create(ctx, cell)
	if err != nil {
		return fmt.Errorf("cellRepo.Create: %w", err)
	}

	return nil
}
