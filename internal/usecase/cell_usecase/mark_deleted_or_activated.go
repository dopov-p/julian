package cell_usecase

import (
	"context"
	"fmt"
)

type MarkDeletedOrActivatedReq struct {
	Name       string
	IsActivate bool
}

func (u *UseCase) MarkDeletedOrActivated(
	ctx context.Context,
	req *MarkDeletedOrActivatedReq,
) error {
	if !req.IsActivate {
		err := u.cellRepo.MarkDeleted(ctx, req.Name)
		if err != nil {
			return fmt.Errorf("cellRepo.MarkDeleted: %w", err)
		}
	} else {
		err := u.cellRepo.MarkActivated(ctx, req.Name)
		if err != nil {
			return fmt.Errorf("cellRepo.MarkActivated: %w", err)
		}
	}

	return nil
}
