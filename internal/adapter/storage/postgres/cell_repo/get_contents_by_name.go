package cell_repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/dopov-p/julian/internal/domain/dto"
	"github.com/dopov-p/julian/internal/domain/model"
	"github.com/dopov-p/julian/internal/pkg"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) GetContentsByName(ctx context.Context, name string) (*dto.GetContentsByNameResponse, error) {
	query := sq.Select("id", "contents").
		From(tableName).
		Where(sq.Eq{"name": name}).
		Where(sq.Eq{"deleted_at": nil}).
		PlaceholderFormat(sq.Dollar)

	// Add row lock if in transaction for atomicity
	if pkg.HasTx(ctx) {
		query = query.Suffix("FOR UPDATE")
	}

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		id           string
		contentsJSON []byte
	)

	err = r.getConn(ctx).QueryRow(ctx, sqlQuery, args...).Scan(&id, &contentsJSON)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w [Name = %s]", model.ErrCellNotFound, name)
		}

		return nil, err
	}

	if len(contentsJSON) == 0 {
		return nil, fmt.Errorf("%w [Name = %s]", model.ErrEmptyContents, name)
	}

	var contents []model.CellContents
	err = json.Unmarshal(contentsJSON, &contents)
	if err != nil {
		return nil, err
	}

	return &dto.GetContentsByNameResponse{
		ID:       id,
		Contents: contents,
	}, nil
}
