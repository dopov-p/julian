package cell_repo

import (
	"context"
	"encoding/json"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/dopov-p/julian/internal/domain/model"
	"github.com/dopov-p/julian/internal/pkg"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) GetByName(ctx context.Context, name string, isActive bool) (*model.Cell, error) {
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
		Where(sq.Eq{"name": name}).
		PlaceholderFormat(sq.Dollar)

	if pkg.HasTx(ctx) {
		query = query.Suffix("FOR UPDATE")
	}

	if isActive {
		query = query.Where(sq.Eq{"deleted_at": nil})
	} else {
		query = query.
			Where(sq.NotEq{"deleted_at": nil}).
			OrderBy("created_at DESC").
			Limit(1)
	}

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		cell         model.Cell
		contentsJSON []byte
	)

	err = r.getConn(ctx).QueryRow(ctx, sqlQuery, args...).Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrCellNotFound
		}

		return nil, err
	}

	if len(contentsJSON) > 0 {
		err = json.Unmarshal(contentsJSON, &cell.Contents)
		if err != nil {
			return nil, err
		}
	}

	return &cell, nil
}
