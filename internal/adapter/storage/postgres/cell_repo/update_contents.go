package cell_repo

import (
	"context"
	"encoding/json"

	sq "github.com/Masterminds/squirrel"
	"github.com/dopov-p/julian/internal/domain/dto"
	"github.com/dopov-p/julian/internal/domain/model"
)

func (r *Repo) UpdateContents(ctx context.Context, req dto.UpdateContentsRequest) error {
	contentsJSON, err := json.Marshal(req.Contents)
	if err != nil {
		return err
	}

	query := sq.Update(tableName).
		Set("contents", contentsJSON).
		Set("updated_at", r.timer.NowUTC()).
		Where(sq.Eq{"id": req.ID}).
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
