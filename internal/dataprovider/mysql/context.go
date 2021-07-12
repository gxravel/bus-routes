package mysql

import (
	"context"

	log "github.com/gxravel/bus-routes/internal/logger"
	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

var (
	errNoRowsAffected = errors.New("no rows affected")
)

// execContext builds the query that doesn't return rows and executes it.
func execContext(ctx context.Context, qb interface{}, entity string, db sqlx.ExtContext) error {
	query, args, codewords, err := toSql(ctx, qb, entity)
	if err != nil {
		return err
	}

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrapf(err, codewords+" with query %s", query)
	}

	num, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to call RowsAffected")
	}

	if num == 0 {
		return errNoRowsAffected
	}

	return nil
}

// toSql builds the query into a SQL string and bound args, and logs the result.
func toSql(ctx context.Context, qb interface{}, entity string) (string, []interface{}, string, error) {
	var (
		query     string
		args      []interface{}
		codewords string
		err       error
	)

	switch qb := qb.(type) {
	case sq.SelectBuilder:
		codewords = "select "
		query, args, err = qb.ToSql()

	case sq.CaseBuilder:
		codewords = "case "
		query, args, err = qb.ToSql()

	case sq.InsertBuilder:
		codewords = "insert "
		query, args, err = qb.ToSql()

	case sq.UpdateBuilder:
		codewords = "update "
		query, args, err = qb.ToSql()

	case sq.DeleteBuilder:
		codewords = "delete "
		query, args, err = qb.ToSql()

	default:
		err = errors.New("wrong query builder")
	}

	codewords += entity
	if err != nil {
		return "", nil, "", errors.Wrap(err, "create sql query for "+codewords)
	}

	log.
		FromContext(ctx).
		WithFields(
			"query", query,
			"args", args).
		Debug(codewords)

	return query, args, codewords, nil
}
