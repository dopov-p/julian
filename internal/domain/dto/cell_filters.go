package dto

import "github.com/dopov-p/julian/internal/domain/model"

type (
	FindAllCellFilter struct {
		Names         []string
		Kind          *model.CellKind
		ContentType   *model.CellContentType
		ContentsEmpty *bool
		IsDeleted     *bool
	}

	UpdateCellContentsRequest struct {
		Name     string
		Contents []model.CellContents
	}
)
