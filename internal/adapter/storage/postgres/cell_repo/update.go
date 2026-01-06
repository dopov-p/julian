package cell_repo

import (
	"context"
	"encoding/json"

	sq "github.com/Masterminds/squirrel"
	"github.com/dopov-p/julian/internal/domain/model"
)

func (r *Repo) Update(ctx context.Context, cell *model.Cell) error {
	contentsJSON, err := json.Marshal(cell.Contents)
	if err != nil {
		return err
	}

	query := sq.Update(tableName).
		Set("type", cell.Kind).
		Set("content_type", cell.ContentType).
		Set("contents", contentsJSON).
		Set("updated_at", r.timer.NowUTC()).
		Where(sq.Eq{"name": cell.Name}).
		Where(sq.Eq{"deleted_at": nil}).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	result, err := r.getConn(ctx).Exec(ctx, sqlQuery, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrCellNotFound
	}

	return nil
}
