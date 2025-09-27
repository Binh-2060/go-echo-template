package services

import (
	"context"

	"github.com/Binh-2060/go-echo-template/api/schema/requestbody"
	dbpkg "github.com/Binh-2060/go-echo-template/pkg/db-pkg"
	dbinserts "github.com/Binh-2060/go-echo-template/pkg/db-pkg/db-inserts"
	dbquery "github.com/Binh-2060/go-echo-template/pkg/db-pkg/db-query"
	dbschema "github.com/Binh-2060/go-echo-template/pkg/db-pkg/db-schema"
)

func CreateUserService(ctx context.Context, req requestbody.UserRequestBody) error {
	tx := dbpkg.GetTransactionManager()
	err := tx.WithTransaction(ctx, func(context context.Context) error {
		db := dbpkg.GetDBFromContext(context)
		return dbinserts.InsertNewUserTx(context, db, req)

	})

	return err
}

func GetUserService(ctx context.Context) ([]dbschema.GetUserDataDBSchema, error) {
	result, err := dbquery.GetUserDataDBQuery(ctx)
	return result, err
}
