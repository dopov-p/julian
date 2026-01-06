package cell_repo

import (
	"context"
	"encoding/json"

	sq "github.com/Masterminds/squirrel"
	"github.com/dopov-p/julian/internal/domain/model"
)

func (r *Repo) Create(ctx context.Context, cell *model.Cell) error {
	contentsJSON, err := json.Marshal(cell.Contents)
	if err != nil {
		return err
	}

	now := r.timer.NowUTC()

	query := sq.Insert(tableName).
		Columns(
			"id",
			"name",
			"kind",
			"content_type",
			"contents",
			"created_at",
			"updated_at",
		).
		Values(
			cell.ID,
			cell.Name,
			cell.Kind,
			cell.ContentType,
			contentsJSON,
			now,
			now,
		).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.getConn(ctx).Exec(ctx, sqlQuery, args...)
	if err != nil {
		return err
	}

	return nil
}
