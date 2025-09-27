package dbquery

import (
	"context"

	"github.com/Binh-2060/go-echo-template/config/db"
	dbpkg "github.com/Binh-2060/go-echo-template/pkg/db-pkg"
	dbschema "github.com/Binh-2060/go-echo-template/pkg/db-pkg/db-schema"
)

func GetUserDataDBQuery(ctx context.Context) ([]dbschema.GetUserDataDBSchema, error) {
	var res = []dbschema.GetUserDataDBSchema{}
	psql := db.GetPSQLCommand()
	query := psql.
		Select("id", "name", "email").
		From("users")
	sql, args, err := query.ToSql()
	if err != nil {
		return res, err
	}

	rows, err := dbpkg.DB.Query(ctx, sql, args...)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var item dbschema.GetUserDataDBSchema
		if err := rows.Scan(&item.ID, &item.Name, &item.Email); err != nil {
			return res, err
		}
		res = append(res, item)
	}

	return res, nil
}
