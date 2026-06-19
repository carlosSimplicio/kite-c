package services

import (
	"context"

	sql "database/sql"
	db "kite-c/database/sqlc"
)

type CommandRepository struct {
	dbClient *sql.DB
	queries  *db.Queries
}

func NewCommandRepository(dbClient *sql.DB) *CommandRepository {
	return &CommandRepository{
		dbClient: dbClient,
		queries:  db.New(dbClient),
	}
}

func (r *CommandRepository) CreateCommand(ctx context.Context, command db.CreateCommandParams) (db.Command, error) {
	return r.queries.CreateCommand(ctx, command)
}

const searchCommandQuery = `-- name: SearchCommand :many
SELECT rowid, rank, name, description FROM command_fts_idx WHERE command_fts_idx MATCH ? ORDER BY rank
`

type SearchCommandRow struct {
	RowId       int64
	Rank        float64
	Name        string
	Description string
}

func (r *CommandRepository) SearchCommand(ctx context.Context, queryString string) ([]SearchCommandRow, error) {
	rows, err := r.dbClient.QueryContext(ctx, searchCommandQuery, queryString+"*")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []SearchCommandRow
	for rows.Next() {
		var row SearchCommandRow
		err := rows.Scan(
			&row.RowId,
			&row.Rank,
			&row.Name,
			&row.Description,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
