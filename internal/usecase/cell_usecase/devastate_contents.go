package cell_usecase

import (
	"context"
	"fmt"

	"github.com/dopov-p/julian/internal/domain/dto"
	"github.com/dopov-p/julian/internal/domain/model"
	"github.com/samber/lo"
)

type DevastateContentsReq struct {
	Name     string
	Contents []model.CellContents
}

func (u *UseCase) DevastateContents(ctx context.Context, req *DevastateContentsReq) error {
	if err := u.txManager.WithTx(ctx, func(txCtx context.Context) error {
		currentCell, err := u.cellRepo.GetContentsByName(txCtx, req.Name)
		if err != nil {
			return fmt.Errorf("cellRepo.GetContentsByName: %w", err)
		}

		resultContents := u.processDevastation(currentCell.Contents, req.Contents)

		err = u.cellRepo.UpdateContents(txCtx, dto.UpdateContentsRequest{
			ID:       currentCell.ID,
			Contents: resultContents,
		})
		if err != nil {
			return fmt.Errorf("cellRepo.UpdateContents: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("txManager.WithTx: %w", err)
	}

	return nil
}

func (u *UseCase) processDevastation(
	currentContents []model.CellContents,
	requestContents []model.CellContents,
) []model.CellContents {
	// Helper to create a composite key for matching items
	createKey := func(item model.CellContents) string {
		orderID := lo.FromPtr(item.ExternalOrderID)
		// Pre-allocate with estimated size for better performance
		key := make([]byte, 0, len(item.SKU)+len(orderID)+1)
		key = append(key, item.SKU...)
		key = append(key, ':')
		key = append(key, orderID...)
		return string(key)
	}

	// Create a map of indices for O(1) lookup
	// Using indices instead of pointers is safer and avoids potential issues
	// with slice reallocation or concurrent modifications
	contentsMap := make(map[string]int, len(currentContents))
	for i := range currentContents {
		key := createKey(currentContents[i])
		contentsMap[key] = i
	}

	// Process each request item and update quantities
	// Work with indices to safely modify the slice
	for _, reqItem := range requestContents {
		key := createKey(reqItem)

		idx, exists := contentsMap[key]
		if !exists {
			continue
		}

		// Safely modify the element using index
		if reqItem.Quantity >= currentContents[idx].Quantity {
			currentContents[idx].Quantity = 0
		} else {
			currentContents[idx].Quantity -= reqItem.Quantity
		}
	}

	// Filter out items with quantity <= 0 using lo.Filter
	result := lo.Filter(currentContents, func(item model.CellContents, _ int) bool {
		return item.Quantity > 0
	})

	return result
}
