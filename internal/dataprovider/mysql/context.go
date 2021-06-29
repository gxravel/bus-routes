package mysql

import (
	"context"
	"fmt"

	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/internal/logger"
	"github.com/gxravel/bus-routes/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type ResultType uint8

const (
	TypeBus ResultType = iota
	TypeCity
	TypeStop
	TypeRoute
	TypeUser
)

func execContext(ctx context.Context, qb interface{}, entity string, txer dataprovider.Txer) error {
	query, args, codewords, err := toSql(ctx, qb, entity)
	if err != nil {
		return err
	}

	f := func(tx *dataprovider.Tx) error {
		result, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return errors.Wrapf(err, codewords+" with query %s", query)
		}
		num, err := result.RowsAffected()
		if err != nil {
			return errors.Wrap(err, "failed to call RowsAffected")
		}
		if num == 0 {
			return errors.New("no rows affected: wrong id")
		}
		return nil
	}

	return dataprovider.BeginAutoCommitedTx(ctx, txer, f)
}

func selectContext(ctx context.Context, qb sq.SelectBuilder, entity string, db sqlx.ExtContext, resultType ResultType) (interface{}, error) {
	query, args, codewords, err := toSql(ctx, qb, entity)
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf(codewords+" by filter with query %s", query)

	switch resultType {
	case TypeCity:
		var result = make([]*model.City, 0)
		if err := sqlx.SelectContext(ctx, db, &result, query, args...); err != nil {
			return nil, errors.Wrapf(err, msg)
		}
		return result, nil
	case TypeBus:
		var result = make([]*model.Bus, 0)
		if err := sqlx.SelectContext(ctx, db, &result, query, args...); err != nil {
			return nil, errors.Wrapf(err, msg)
		}
		return result, nil
	case TypeStop:
		var result = make([]*model.Stop, 0)
		if err := sqlx.SelectContext(ctx, db, &result, query, args...); err != nil {
			return nil, errors.Wrapf(err, msg)
		}
		return result, nil
	case TypeRoute:
		var result = make([]*model.Route, 0)
		if err := sqlx.SelectContext(ctx, db, &result, query, args...); err != nil {
			return nil, errors.Wrapf(err, msg)
		}
		return result, nil
	case TypeUser:
		var result = make([]*model.User, 0)
		if err := sqlx.SelectContext(ctx, db, &result, query, args...); err != nil {
			return nil, errors.Wrapf(err, msg)
		}
		return result, nil
	default:
		return nil, errors.New("wrong result type")
	}
}

func toSql(ctx context.Context, qb interface{}, entity string) (string, []interface{}, string, error) {
	var (
		query     string
		args      []interface{}
		codewords string
		err       error
	)
	switch qb := qb.(type) {
	case sq.SelectBuilder:
		codewords = "selecting "
		query, args, err = qb.ToSql()
	case sq.CaseBuilder:
		codewords = "casing "
		query, args, err = qb.ToSql()
	case sq.InsertBuilder:
		codewords = "inserting "
		query, args, err = qb.ToSql()
	case sq.UpdateBuilder:
		codewords = "updating "
		query, args, err = qb.ToSql()
	case sq.DeleteBuilder:
		codewords = "deleting "
		query, args, err = qb.ToSql()
	default:
		err = errors.New("wrong query builder")
	}
	codewords += entity
	if err != nil {
		return "", nil, "", errors.Wrap(err, "creating sql query for "+codewords)
	}

	logger.FromContext(ctx).WithFields(
		"query", query,
		"args", args).
		Debug(codewords)

	return query, args, codewords, nil
}
