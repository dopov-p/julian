package cell_repo

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/dopov-p/julian/internal/domain/model"
)

func (r *Repo) SetSellEmpty(ctx context.Context, name string) error {
	query := sq.Update(tableName).
		Set("contents", nil).
		Set("updated_at", r.timer.NowUTC()).
		Where(sq.Eq{"name": name}).
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
