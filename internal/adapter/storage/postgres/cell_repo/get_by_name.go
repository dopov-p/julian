package cell_repo

import (
	"context"
	"encoding/json"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/dopov-p/julian/internal/domain/model"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) GetByName(ctx context.Context, name string) (*model.Cell, error) {
	query := sq.Select(
		"id",
		"name",
		"kind",
		"content_type", "contents", "created_at", "updated_at", "deleted_at").
		From("cells").
		Where(sq.Eq{"name": name}).
		Where(sq.Eq{"deleted_at": nil}).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		cell         model.Cell
		contentsJSON []byte
	)

	err = r.cluster.Conn.QueryRow(ctx, sqlQuery, args...).Scan(
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
			return nil, ErrNotFound
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
