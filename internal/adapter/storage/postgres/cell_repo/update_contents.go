package cell_repo

import (
	"context"
	"encoding/json"

	sq "github.com/Masterminds/squirrel"
	"github.com/dopov-p/julian/internal/domain/dto"
)

func (r *Repo) UpdateContents(ctx context.Context, req dto.UpdateCellContentsRequest) error {
	contentsJSON, err := json.Marshal(req.Contents)
	if err != nil {
		return err
	}

	query := sq.Update("cells").
		Set("contents", contentsJSON).
		Set("updated_at", r.timer.NowUTC()).
		Where(sq.Eq{"name": req.Name}).
		Where(sq.Eq{"deleted_at": nil}).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	result, err := r.cluster.Conn.Exec(ctx, sqlQuery, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
