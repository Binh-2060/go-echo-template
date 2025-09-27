package dbinserts

import (
	"context"

	"github.com/Binh-2060/go-echo-template/api/schema/requestbody"
	"github.com/Binh-2060/go-echo-template/config/db"
	dbpkg "github.com/Binh-2060/go-echo-template/pkg/db-pkg"
)

func InsertNewUserTx(ctx context.Context, tx dbpkg.DBTX, req requestbody.UserRequestBody) error {
	psql := db.GetPSQLCommand()
	query := psql.
		Insert("users").
		Columns("name", "email").
		Values(req.Name, req.Email)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	return err
}
