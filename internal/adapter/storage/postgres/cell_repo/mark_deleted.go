package cell_repo

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/dopov-p/julian/internal/domain/model"
)

func (r *Repo) MarkDeleted(ctx context.Context, name string) error {
	now := r.timer.NowUTC()

	query := sq.Update(tableName).
		Set("deleted_at", now).
		Set("updated_at", now).
		Where(sq.Eq{"name": name}).
		Where(sq.Eq{"deleted_at": nil}).
		Where(sq.Eq{"contents": nil})

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
