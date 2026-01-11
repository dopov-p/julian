package cell_repo

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/dopov-p/julian/internal/domain/model"
)

func (r *Repo) MarkActivated(ctx context.Context, name string) error {
	now := r.timer.NowUTC()

	query := sq.Update(tableName).
		Set("deleted_at", nil).
		Set("updated_at", now).
		Where(sq.Eq{"name": name}).
		Where(sq.NotEq{"deleted_at": nil}).
		OrderBy("created_at DESC").
		Limit(1)

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
