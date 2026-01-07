package cell_repo

import (
	"context"
	"encoding/json"

	sq "github.com/Masterminds/squirrel"
	"github.com/dopov-p/julian/internal/domain/dto"
	"github.com/dopov-p/julian/internal/domain/model"
)

func (r *Repo) FindAllByFilter(
	ctx context.Context,
	filter dto.FindAllCellFilter,
	pagination dto.Pagination,
) ([]*model.Cell, error) {
	query := sq.Select(
		"id",
		"name",
		"kind",
		"content_type",
		"contents",
		"created_at",
		"updated_at",
		"deleted_at",
	).
		From(tableName).
		PlaceholderFormat(sq.Dollar)

	if len(filter.Names) > 0 {
		query = query.Where(sq.Eq{"name": filter.Names})
	}
	if filter.Kind != nil {
		query = query.Where(sq.Eq{"type": *filter.Kind})
	}
	if filter.ContentType != nil {
		query = query.Where(sq.Eq{"content_type": *filter.ContentType})
	}
	if filter.IsDeleted != nil {
		if *filter.IsDeleted {
			query = query.Where(sq.Eq{"is_deleted": nil})
		} else {
			query = query.Where(sq.NotEq{"is_deleted": nil})
		}
	}
	if filter.ContentsEmpty != nil {
		if *filter.ContentsEmpty {
			query = query.Where(sq.Or{
				sq.Eq{"contents": nil},
				sq.Expr("contents = '[]'::jsonb"),
				sq.Expr("contents = 'null'::jsonb"),
			})
		} else {
			query = query.Where(sq.And{
				sq.NotEq{"contents": nil},
				sq.Expr("contents != '[]'::jsonb"),
				sq.Expr("contents != 'null'::jsonb"),
			})
		}
	}

	if pagination.OrderBy != nil && *pagination.OrderBy != "" {
		order := "asc"
		if pagination.Order != nil {
			order = string(*pagination.Order)
		}
		query = query.OrderBy(*pagination.OrderBy + " " + order)
	}

	limit := defaultBatchSize
	if pagination.Limit != nil && *pagination.Limit > 0 && *pagination.Limit <= 50 {
		limit = *pagination.Limit
	}
	query = query.Limit(limit)

	if pagination.Offset != nil && *pagination.Offset > 0 {
		query = query.Offset(*pagination.Offset)
	}

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.getConn(ctx).Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cells []*model.Cell
	for rows.Next() {
		var cell model.Cell
		var contentsJSON []byte

		err = rows.Scan(
			&cell.ID,
			&cell.Name,
			&cell.Kind,
			&cell.ContentType,
			&contentsJSON,
			&cell.CreatedAt,
			&cell.UpdatedAt,
			&cell.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		if len(contentsJSON) > 0 {
			err = json.Unmarshal(contentsJSON, &cell.Contents)
			if err != nil {
				return nil, err
			}
		}

		cells = append(cells, &cell)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cells, nil
}
